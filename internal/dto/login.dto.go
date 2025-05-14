package dto

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/saifoelloh/ranger/internal/constant"
)

type LoginRequest struct {
	Email       *string               `json:"email"`
	Password    *string               `json:"password"`
	SSOID       *string               `json:"sso_id"`
	SSOPlatform *constant.SSOPlatform `json:"sso_platform"`
	Device      string                `json:"device"`
	MacAddress  string                `json:"mac_address"`
	PublicKey   string                `json:"public_key"`
}

type LoginInput struct {
	Email         *string               `json:"email"`
	Password      *string               `json:"password"`
	SSOID         *string               `json:"sso_id"`
	SSOPlatform   *constant.SSOPlatform `json:"sso_platform"`
	Device        string                `json:"device"`
	MacAddress    string                `json:"mac_address"`
	PublicKey     string                `json:"public_key"`
	UserAgent     string                `json:"user_agent"`
	IP            string                `json:"ip"`
	Location      string                `json:"location"`
	ClientVersion string                `json:"client_version"`
}

type LoginResponse struct {
	SessionID    string `json:"session_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AppClaims struct {
	UserID    string `json:"user_id"`
	UserType  string `json:"user_type"`
	UserToken string `json:"user_token"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	SessionID string `json:"session_id"`
	jwt.RegisteredClaims
}

type UserAgent struct {
	Device     string `json:"device"`
	Os         string `json:"os"`
	Raw        string `json:"raw"`
	RedisLabel string `json:"redis_label"`
}
