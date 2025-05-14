package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/saifoelloh/ranger/internal/dto"
	"github.com/saifoelloh/ranger/internal/model"
	repository "github.com/saifoelloh/ranger/internal/repositories"
	"github.com/saifoelloh/ranger/internal/utils"
	"github.com/saifoelloh/ranger/pkg/errors"
)

type AuthService struct {
	userRepo    *repository.UserRepository
	sessionRepo *repository.SessionRepository
	jwtSecret   string
}

func NewAuthService(userRepo *repository.UserRepository, sessionRepo *repository.SessionRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		jwtSecret:   jwtSecret,
	}
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginInput) (*dto.LoginResponse, error) {
	var user *model.User
	var err error

	if req.SSOID != nil && *req.SSOID != "" {
		user, err = s.userRepo.FindBySSOID(*req.SSOID)
		if err != nil {
			return nil, errors.Unauthorized(
				errors.WithScope("AuthService"),
				errors.WithLocation("Login.FindByEmail"),
				errors.WithMessage("invalid email or password"),
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
				errors.WithMessage("invalid email or password"),
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

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(*req.Password)); err != nil {
		return nil, errors.Unauthorized(
			errors.WithScope("AuthService"),
			errors.WithLocation("Login.ComparePassword"),
			errors.WithMessage("invalid email or password"),
			errors.WithErrorCode("auth/invalid-credentials"),
		)
	}

	fmt.Println(user.ID)
	if err := s.sessionRepo.DeactivateSessionsByUserID(user.ID); err != nil {
		return nil, err
	}

	sessionID := uuid.New().String()
	session := &model.Session{
		ID:         sessionID,
		UserID:     user.ID,
		Device:     req.Device,
		MacAddress: req.MacAddress,
		PublicKey:  req.PublicKey,
		Active:     true,
		IP:         sql.NullString{String: req.IP},
		UserAgent:  sql.NullString{String: req.UserAgent},
		Location:   sql.NullString{String: req.UserAgent},
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
		},
	})

	signedToken, _ := token.SignedString([]byte(s.jwtSecret))

	refreshToken := uuid.New().String()

	return &dto.LoginResponse{
		AccessToken:  signedToken,
		RefreshToken: refreshToken,
		SessionID:    sessionID,
	}, nil
}
