package service

import (
	"context"
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/saifoelloh/ranger/internal/config"
	"github.com/saifoelloh/ranger/internal/dto"
	"github.com/saifoelloh/ranger/internal/model"
	"github.com/saifoelloh/ranger/internal/redis"
	repository "github.com/saifoelloh/ranger/internal/repositories"
	"github.com/saifoelloh/ranger/internal/utils"
	"github.com/saifoelloh/ranger/pkg/errors"
)

type AuthService struct {
	config           config.Config
	userRepo         *repository.UserRepository
	sessionRepo      *repository.SessionRepository
	rateLimiterRedis *redis.RateLimiterRepository
	tokenCacheRedis  *redis.TokenRepository
}

func NewAuthService(
	config config.Config,
	userRepo *repository.UserRepository,
	sessionRepo *repository.SessionRepository,
	rateLimiterRedis *redis.RateLimiterRepository,
	tokenCacheRedis *redis.TokenRepository,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		sessionRepo:      sessionRepo,
		config:           config,
		rateLimiterRedis: rateLimiterRedis,
		tokenCacheRedis:  tokenCacheRedis,
	}
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginInput) (*dto.LoginResponse, error) {
	uniqueLabel := utils.GetUniqueLabel(req.Email, req.SSOID)
	s.rateLimiterRedis.IsAllowed(ctx, uniqueLabel)

	var user *model.User
	var err error

	if req.SSOID != nil && *req.SSOID != "" {
		user, err = s.userRepo.FindBySSOID(*req.SSOID, *req.SSOPlatform)
		if err != nil {
			return nil, errors.Unauthorized(
				errors.WithScope("AuthService"),
				errors.WithLocation("Login.FindBySSO"),
				errors.WithMessage("invalid credentials"),
				errors.WithErrorCode("auth/invalid-credentials"),
			)
		}
	}

	if req.Email != nil && *req.Email != "" {
		hashedEmail := utils.CryptoHash(*req.Email)
		user, err = s.userRepo.FindByEmail(hashedEmail)
		if err != nil {
			return nil, errors.Unauthorized(
				errors.WithScope("AuthService"),
				errors.WithLocation("Login.FindByEmail"),
				errors.WithMessage("invalid credentials"),
				errors.WithErrorCode("auth/invalid-credentials"),
			)
		}
	}

	if user == nil {
		return nil, errors.Unauthorized(
			errors.WithScope("AuthService"),
			errors.WithLocation("Login.NoUserFound"),
			errors.WithMessage("no user found for given credentials"),
			errors.WithErrorCode("auth/user-not-found"),
		)
	}

	if req.Password == nil && req.SSOID == nil {
		return nil, errors.Unauthorized(
			errors.WithScope("AuthService"),
			errors.WithLocation("Login.NoUserFound"),
			errors.WithMessage("no user found for given credentials"),
			errors.WithErrorCode("auth/user-not-found"),
		)
	} else if req.Password != nil && *req.Password != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(*req.Password)); err != nil {
			return nil, errors.Unauthorized(
				errors.WithScope("AuthService"),
				errors.WithLocation("Login.ComparePassword"),
				errors.WithMessage("invalid email or password"),
				errors.WithErrorCode("auth/invalid-credentials"),
			)
		}
	}

	if err := s.sessionRepo.DeactivateSessionsByUserID(user.ID); err != nil {
		return nil, err
	}

	sessionID := uuid.New().String()
	session := &model.Session{
		ID:            sessionID,
		UserID:        user.ID,
		Device:        req.Device,
		MacAddress:    req.MacAddress,
		PublicKey:     req.PublicKey,
		Active:        true,
		IP:            sql.NullString{String: req.IP, Valid: true},
		UserAgent:     sql.NullString{String: req.UserAgent, Valid: true},
		Location:      sql.NullString{String: req.UserAgent, Valid: true},
		ClientVersion: req.ClientVersion,
	}

	if err := s.sessionRepo.CreateSession(session); err != nil {
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, dto.AppClaims{
		UserID:    user.ID,
		UserType:  "Investor",
		UserToken: user.InvestorType.String,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			Issuer:    s.config.JwtIssuer,
		},
	})

	signedToken, _ := token.SignedString([]byte(s.config.JwtSecret))
	refreshToken := uuid.New().String()

	s.tokenCacheRedis.SetAccessToken(ctx, user.ID, signedToken, 1*time.Hour)
	s.rateLimiterRedis.Reset(ctx, uniqueLabel)

	return &dto.LoginResponse{
		AccessToken:  signedToken,
		RefreshToken: refreshToken,
		SessionID:    sessionID,
	}, nil
}
