package application

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"polomnik/internal/domain"
	"polomnik/internal/ports"
)

const phoneChallengeTTL = 5 * time.Minute

const unavailableMessage = "Пока что недоступно, используйте другой вариант."

type AuthService struct {
	users     ports.UserRepository
	bookings  ports.BookingRepository
	phones    ports.PhoneVerifier
	jwtSecret []byte
	tokenTTL  time.Duration

	mu         sync.Mutex
	challenges map[string]phoneChallengeRecord
}

type phoneChallengeRecord struct {
	Phone     string
	ExpiresAt time.Time
}

func NewAuthService(
	users ports.UserRepository,
	bookings ports.BookingRepository,
	phones ports.PhoneVerifier,
	jwtSecret string,
	tokenTTL time.Duration,
) *AuthService {
	if tokenTTL <= 0 {
		tokenTTL = 7 * 24 * time.Hour
	}
	if phones == nil {
		phones = unavailablePhoneVerifier{}
	}
	return &AuthService{
		users:      users,
		bookings:   bookings,
		phones:     phones,
		jwtSecret:  []byte(jwtSecret),
		tokenTTL:   tokenTTL,
		challenges: make(map[string]phoneChallengeRecord),
	}
}

type unavailablePhoneVerifier struct{}

func (unavailablePhoneVerifier) Available() bool { return false }
func (unavailablePhoneVerifier) Start(context.Context, string) (ports.PhoneChallenge, error) {
	return ports.PhoneChallenge{}, ports.ErrPhoneVerifierNotConfigured
}
func (unavailablePhoneVerifier) Status(context.Context, string) (ports.PhoneCheckStatus, error) {
	return "", ports.ErrPhoneVerifierNotConfigured
}

type RegisterInput struct {
	Email         string
	Phone         string
	Name          string
	Password      string
	PhoneCheckID  string
}

type LoginInput struct {
	Login    string
	Password string
}

type AuthResult struct {
	Token string
	User  domain.User
}

type AuthClaims struct {
	UserID uuid.UUID `json:"uid"`
	jwt.RegisteredClaims
}

type AuthMethods struct {
	Password  bool
	PhoneCall AuthMethodStatus
}

type AuthMethodStatus struct {
	Available bool
	Message   string
}

type PhoneStartResult struct {
	CheckID         string
	CallPhone       string
	CallPhonePretty string
	ExpiresIn       int
}

type PhoneStatusResult struct {
	Status string
}

func (s *AuthService) AuthMethods() AuthMethods {
	phone := AuthMethodStatus{Available: false, Message: unavailableMessage}
	if s.phones.Available() {
		phone = AuthMethodStatus{Available: true}
	}
	return AuthMethods{
		Password:  true,
		PhoneCall: phone,
	}
}

func (s *AuthService) StartPhoneVerification(ctx context.Context, phone string) (PhoneStartResult, error) {
	if !s.phones.Available() {
		return PhoneStartResult{}, ErrPhoneVerificationUnavailable
	}
	normalized := domain.NormalizePhone(phone)
	if normalized == "" {
		return PhoneStartResult{}, domain.ErrInvalidPhone
	}

	challenge, err := s.phones.Start(ctx, normalized)
	if err != nil {
		if errors.Is(err, ports.ErrPhoneVerifierNotConfigured) {
			return PhoneStartResult{}, ErrPhoneVerificationUnavailable
		}
		return PhoneStartResult{}, err
	}

	expiresIn := challenge.ExpiresIn
	if expiresIn <= 0 {
		expiresIn = int(phoneChallengeTTL.Seconds())
	}

	s.mu.Lock()
	s.pruneChallengesLocked(time.Now().UTC())
	s.challenges[challenge.CheckID] = phoneChallengeRecord{
		Phone:     normalized,
		ExpiresAt: time.Now().UTC().Add(time.Duration(expiresIn) * time.Second),
	}
	s.mu.Unlock()

	return PhoneStartResult{
		CheckID:         challenge.CheckID,
		CallPhone:       challenge.CallPhone,
		CallPhonePretty: challenge.CallPhonePretty,
		ExpiresIn:       expiresIn,
	}, nil
}

func (s *AuthService) PhoneVerificationStatus(ctx context.Context, checkID string) (PhoneStatusResult, error) {
	if !s.phones.Available() {
		return PhoneStatusResult{}, ErrPhoneVerificationUnavailable
	}
	checkID = strings.TrimSpace(checkID)
	if checkID == "" {
		return PhoneStatusResult{}, domain.ErrInvalidID
	}

	if _, ok := s.lookupChallenge(checkID); !ok {
		return PhoneStatusResult{Status: string(ports.PhoneCheckExpired)}, nil
	}

	status, err := s.phones.Status(ctx, checkID)
	if err != nil {
		if errors.Is(err, ports.ErrPhoneVerifierNotConfigured) {
			return PhoneStatusResult{}, ErrPhoneVerificationUnavailable
		}
		return PhoneStatusResult{}, err
	}
	return PhoneStatusResult{Status: string(status)}, nil
}

func (s *AuthService) CompletePhoneLogin(ctx context.Context, checkID string) (AuthResult, error) {
	phone, err := s.requireConfirmedPhone(ctx, checkID)
	if err != nil {
		return AuthResult{}, err
	}

	user, err := s.users.GetUserByPhone(ctx, phone)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return AuthResult{}, ErrPhoneUserNotFound
		}
		return AuthResult{}, err
	}

	token, err := s.issueToken(user.ID)
	if err != nil {
		return AuthResult{}, err
	}
	return AuthResult{Token: token, User: sanitizeUser(user)}, nil
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (AuthResult, error) {
	if err := domain.ValidatePassword(input.Password); err != nil {
		return AuthResult{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return AuthResult{}, err
	}

	user, err := domain.NewUser(domain.RegisterUserInput{
		ID:           uuid.New(),
		Email:        input.Email,
		Phone:        input.Phone,
		Name:         input.Name,
		PasswordHash: string(hash),
	})
	if err != nil {
		return AuthResult{}, err
	}

	if s.phones.Available() {
		confirmedPhone, err := s.requireConfirmedPhone(ctx, input.PhoneCheckID)
		if err != nil {
			return AuthResult{}, err
		}
		if confirmedPhone != user.Phone {
			return AuthResult{}, ErrPhoneVerificationNotConfirmed
		}
	}

	if user.Email != "" {
		if _, err := s.users.GetUserByEmail(ctx, user.Email); err == nil {
			return AuthResult{}, domain.ErrDuplicateEmail
		} else if !errors.Is(err, domain.ErrNotFound) {
			return AuthResult{}, err
		}
	}

	if _, err := s.users.GetUserByPhone(ctx, user.Phone); err == nil {
		return AuthResult{}, domain.ErrDuplicatePhone
	} else if !errors.Is(err, domain.ErrNotFound) {
		return AuthResult{}, err
	}

	created, err := s.users.CreateUser(ctx, user)
	if err != nil {
		return AuthResult{}, err
	}

	token, err := s.issueToken(created.ID)
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{Token: token, User: sanitizeUser(created)}, nil
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (AuthResult, error) {
	login := strings.TrimSpace(input.Login)
	if login == "" || strings.TrimSpace(input.Password) == "" {
		return AuthResult{}, domain.ErrInvalidCredentials
	}

	var user domain.User
	var err error
	if strings.Contains(login, "@") {
		user, err = s.users.GetUserByEmail(ctx, strings.ToLower(login))
	} else {
		user, err = s.users.GetUserByPhone(ctx, domain.NormalizePhone(login))
	}
	if err != nil {
		return AuthResult{}, domain.ErrInvalidCredentials
	}
	if strings.TrimSpace(user.PasswordHash) == "" {
		return AuthResult{}, domain.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return AuthResult{}, domain.ErrInvalidCredentials
	}

	token, err := s.issueToken(user.ID)
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{Token: token, User: sanitizeUser(user)}, nil
}

func (s *AuthService) ParseToken(token string) (uuid.UUID, error) {
	if strings.TrimSpace(token) == "" {
		return uuid.Nil, domain.ErrInvalidCredentials
	}

	parsed, err := jwt.ParseWithClaims(token, &AuthClaims{}, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, domain.ErrInvalidCredentials
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return uuid.Nil, domain.ErrInvalidCredentials
	}

	claims, ok := parsed.Claims.(*AuthClaims)
	if !ok || !parsed.Valid || claims.UserID == uuid.Nil {
		return uuid.Nil, domain.ErrInvalidCredentials
	}

	return claims.UserID, nil
}

func (s *AuthService) GetUser(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	user, err := s.users.GetUserByID(ctx, userID)
	if err != nil {
		return domain.User{}, err
	}
	return sanitizeUser(user), nil
}

func (s *AuthService) ListMyBookings(ctx context.Context, userID uuid.UUID, pagination ports.Pagination) (ports.BookingList, error) {
	return s.bookings.ListBookings(ctx, ports.BookingFilters{UserID: &userID}, pagination)
}

type OAuthLoginInput struct {
	Provider string
	Subject  string
	Email    string
	Name     string
}

func (s *AuthService) OAuthLogin(ctx context.Context, input OAuthLoginInput) (AuthResult, error) {
	user, err := s.users.GetUserByOAuth(ctx, input.Provider, input.Subject)
	if errors.Is(err, domain.ErrNotFound) {
		user, err = domain.NewOAuthUser(domain.OAuthUserInput{
			ID:       uuid.New(),
			Provider: input.Provider,
			Subject:  input.Subject,
			Email:    input.Email,
			Name:     input.Name,
		})
		if err != nil {
			return AuthResult{}, err
		}
		user, err = s.users.CreateUser(ctx, user)
		if err != nil {
			return AuthResult{}, err
		}
	} else if err != nil {
		return AuthResult{}, err
	}

	token, err := s.issueToken(user.ID)
	if err != nil {
		return AuthResult{}, err
	}
	return AuthResult{Token: token, User: sanitizeUser(user)}, nil
}

func (s *AuthService) requireConfirmedPhone(ctx context.Context, checkID string) (string, error) {
	if !s.phones.Available() {
		return "", ErrPhoneVerificationUnavailable
	}
	checkID = strings.TrimSpace(checkID)
	if checkID == "" {
		return "", ErrPhoneVerificationRequired
	}

	record, ok := s.lookupChallenge(checkID)
	if !ok {
		return "", ErrPhoneVerificationNotConfirmed
	}

	status, err := s.phones.Status(ctx, checkID)
	if err != nil {
		if errors.Is(err, ports.ErrPhoneVerifierNotConfigured) {
			return "", ErrPhoneVerificationUnavailable
		}
		return "", err
	}
	switch status {
	case ports.PhoneCheckConfirmed:
		s.forgetChallenge(checkID)
		return record.Phone, nil
	case ports.PhoneCheckPending:
		return "", ErrPhoneVerificationNotConfirmed
	case ports.PhoneCheckExpired:
		s.forgetChallenge(checkID)
		return "", ErrPhoneVerificationNotConfirmed
	default:
		return "", ErrPhoneVerificationNotConfirmed
	}
}

func (s *AuthService) lookupChallenge(checkID string) (phoneChallengeRecord, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pruneChallengesLocked(time.Now().UTC())
	record, ok := s.challenges[checkID]
	return record, ok
}

func (s *AuthService) forgetChallenge(checkID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.challenges, checkID)
}

func (s *AuthService) pruneChallengesLocked(now time.Time) {
	for id, record := range s.challenges {
		if now.After(record.ExpiresAt) {
			delete(s.challenges, id)
		}
	}
}

func (s *AuthService) issueToken(userID uuid.UUID) (string, error) {
	now := time.Now().UTC()
	claims := AuthClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func sanitizeUser(user domain.User) domain.User {
	user.PasswordHash = ""
	return user
}
