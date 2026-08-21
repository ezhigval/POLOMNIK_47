package application

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"polomnik/internal/domain"
	"polomnik/internal/ports"
)

type AuthService struct {
	users     ports.UserRepository
	bookings  ports.BookingRepository
	jwtSecret []byte
	tokenTTL  time.Duration
}

func NewAuthService(
	users ports.UserRepository,
	bookings ports.BookingRepository,
	jwtSecret string,
	tokenTTL time.Duration,
) *AuthService {
	if tokenTTL <= 0 {
		tokenTTL = 7 * 24 * time.Hour
	}
	return &AuthService{
		users:     users,
		bookings:  bookings,
		jwtSecret: []byte(jwtSecret),
		tokenTTL:  tokenTTL,
	}
}

type RegisterInput struct {
	Email    string
	Phone    string
	Name     string
	Password string
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
