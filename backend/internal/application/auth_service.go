package application

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

const phoneChallengeTTL = 5 * time.Minute
const passwordResetTTL = time.Hour
const passwordResetPurpose = "password_reset"

const unavailableMessage = "Пока что недоступно, используйте другой вариант."

type AuthService struct {
	users         ports.UserRepository
	bookings      ports.BookingRepository
	phones        ports.PhoneVerifier
	mailer        ports.Mailer
	social        SocialAuthConfig
	jwtSecret     []byte
	tokenTTL      time.Duration
	publicSiteURL string

	mu         sync.Mutex
	challenges map[string]phoneChallengeRecord
}

type SocialAuthConfig struct {
	YandexConfigured   bool
	VKConfigured       bool
	MaxConfigured      bool
	TelegramConfigured bool
	TelegramBotUsername string
}

type phoneChallengeRecord struct {
	Phone     string
	ExpiresAt time.Time
}

func NewAuthService(
	users ports.UserRepository,
	bookings ports.BookingRepository,
	phones ports.PhoneVerifier,
	mailer ports.Mailer,
	social SocialAuthConfig,
	jwtSecret string,
	tokenTTL time.Duration,
	publicSiteURL string,
) *AuthService {
	if tokenTTL <= 0 {
		tokenTTL = 7 * 24 * time.Hour
	}
	if phones == nil {
		phones = unavailablePhoneVerifier{}
	}
	if mailer == nil {
		mailer = unavailableMailer{}
	}
	return &AuthService{
		users:         users,
		bookings:      bookings,
		phones:        phones,
		mailer:        mailer,
		social:        social,
		jwtSecret:     []byte(jwtSecret),
		tokenTTL:      tokenTTL,
		publicSiteURL: strings.TrimRight(strings.TrimSpace(publicSiteURL), "/"),
		challenges:    make(map[string]phoneChallengeRecord),
	}
}

type unavailableMailer struct{}

func (unavailableMailer) Configured() bool { return false }
func (unavailableMailer) Send(context.Context, ports.MailMessage) error {
	return ports.ErrMailerNotConfigured
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

type passwordResetClaims struct {
	UserID      uuid.UUID `json:"uid"`
	Purpose     string    `json:"purpose"`
	Fingerprint string    `json:"fp"`
	jwt.RegisteredClaims
}

type AuthMethods struct {
	Password  bool
	PhoneCall AuthMethodStatus
	Yandex    AuthMethodStatus
	VK        AuthMethodStatus
	Max       AuthMethodStatus
	Telegram  AuthMethodStatus
	Mail      AuthMethodStatus
}

type AuthMethodStatus struct {
	Available bool
	Message   string
	Username  string
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
	socialStatus := func(ok bool) AuthMethodStatus {
		if ok {
			return AuthMethodStatus{Available: true}
		}
		return AuthMethodStatus{Available: false, Message: unavailableMessage}
	}
	mail := AuthMethodStatus{Available: false, Message: unavailableMessage}
	if s.mailer != nil && s.mailer.Configured() {
		mail = AuthMethodStatus{Available: true}
	}
	telegram := socialStatus(s.social.TelegramConfigured)
	telegram.Username = s.social.TelegramBotUsername
	return AuthMethods{
		Password:  true,
		PhoneCall: phone,
		Yandex:    socialStatus(s.social.YandexConfigured),
		VK:        socialStatus(s.social.VKConfigured),
		Max:       socialStatus(s.social.MaxConfigured),
		Telegram:  telegram,
		Mail:      mail,
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

	s.sendRegistrationMail(ctx, created)

	token, err := s.issueToken(created.ID)
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{Token: token, User: sanitizeUser(created)}, nil
}

func (s *AuthService) sendRegistrationMail(ctx context.Context, user domain.User) {
	if s.mailer == nil || !s.mailer.Configured() || strings.TrimSpace(user.Email) == "" {
		return
	}
	_ = s.mailer.Send(ctx, ports.MailMessage{
		To:      []string{user.Email},
		Subject: "Регистрация на сайте паломнической службы",
		Text:    "Здравствуйте, " + user.Name + "!\n\nВаш аккаунт создан. Войти можно на сайте в разделе «Кабинет».\n",
	})
}

// RequestPasswordReset sends a reset link when mailer is configured and the email
// matches a password-based account. Always succeeds for unknown emails (no enumeration).
func (s *AuthService) RequestPasswordReset(ctx context.Context, email string) error {
	if s.mailer == nil || !s.mailer.Configured() {
		return ErrPasswordResetUnavailable
	}
	normalized := strings.TrimSpace(strings.ToLower(email))
	if normalized == "" || !strings.Contains(normalized, "@") {
		return domain.ErrInvalidEmail
	}

	user, err := s.users.GetUserByEmail(ctx, normalized)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil
		}
		return err
	}
	if strings.TrimSpace(user.PasswordHash) == "" {
		return nil
	}

	token, err := s.issuePasswordResetToken(user)
	if err != nil {
		return err
	}
	link := s.publicSiteURL + "/account/reset-password?token=" + url.QueryEscape(token)
	_ = s.mailer.Send(ctx, ports.MailMessage{
		To:      []string{user.Email},
		Subject: "Восстановление пароля",
		Text: "Здравствуйте" + nameSuffix(user.Name) + "!\n\n" +
			"Чтобы задать новый пароль, откройте ссылку (действует около часа):\n" +
			link + "\n\nЕсли вы не запрашивали восстановление, просто проигнорируйте это письмо.\n",
	})
	return nil
}

func (s *AuthService) ResetPassword(ctx context.Context, token, newPassword string) error {
	if err := domain.ValidatePassword(newPassword); err != nil {
		return err
	}
	claims, err := s.parsePasswordResetToken(token)
	if err != nil {
		return err
	}
	user, err := s.users.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return ErrInvalidPasswordResetToken
	}
	if passwordFingerprint(user.PasswordHash) != claims.Fingerprint {
		return ErrInvalidPasswordResetToken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.users.UpdateUserPassword(ctx, user.ID, string(hash))
}

func nameSuffix(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}
	return ", " + name
}

func passwordFingerprint(passwordHash string) string {
	sum := sha256.Sum256([]byte(passwordHash))
	return hex.EncodeToString(sum[:8])
}

func (s *AuthService) issuePasswordResetToken(user domain.User) (string, error) {
	now := time.Now().UTC()
	claims := passwordResetClaims{
		UserID:      user.ID,
		Purpose:     passwordResetPurpose,
		Fingerprint: passwordFingerprint(user.PasswordHash),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(passwordResetTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   user.ID.String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) parsePasswordResetToken(token string) (passwordResetClaims, error) {
	if strings.TrimSpace(token) == "" {
		return passwordResetClaims{}, ErrInvalidPasswordResetToken
	}
	parsed, err := jwt.ParseWithClaims(token, &passwordResetClaims{}, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidPasswordResetToken
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return passwordResetClaims{}, ErrInvalidPasswordResetToken
	}
	claims, ok := parsed.Claims.(*passwordResetClaims)
	if !ok || !parsed.Valid || claims.UserID == uuid.Nil || claims.Purpose != passwordResetPurpose {
		return passwordResetClaims{}, ErrInvalidPasswordResetToken
	}
	return *claims, nil
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
