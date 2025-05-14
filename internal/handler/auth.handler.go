package handler

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/saifoelloh/ranger/internal/dto"
	service "github.com/saifoelloh/ranger/internal/services"
	"github.com/saifoelloh/ranger/internal/utils"
	"github.com/saifoelloh/ranger/pkg/errors"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequest(
			errors.WithScope("AuthHandler"),
			errors.WithLocation("Login.BindJSON"),
			errors.WithMessage("invalid request body"),
			errors.WithErrorCode("auth/invalid-json"),
			errors.WithDetail(err.Error()),
		))
		return
	}

	uniqueLable := utils.GetUniqueLabel(req.Email, req.SSOID)
	rawUserAgent := c.Request.UserAgent()
	userAgent := dto.UserAgent{
		Device:     utils.ParseUserAgent(rawUserAgent).Device,
		Os:         utils.ParseUserAgent(rawUserAgent).OS,
		Raw:        rawUserAgent,
		RedisLabel: uniqueLable,
	}
	formattedUserAgent, _ := json.Marshal(userAgent)

	resp, err := h.authService.Login(
		c.Request.Context(),
		dto.LoginInput{
			Email:         req.Email,
			Password:      req.Password,
			SSOID:         req.SSOID,
			SSOPlatform:   req.SSOPlatform,
			Device:        req.Device,
			MacAddress:    req.MacAddress,
			PublicKey:     req.PublicKey,
			UserAgent:     string(formattedUserAgent),
			IP:            c.ClientIP(),
			Location:      c.GetHeader("X-Location"),
			ClientVersion: c.GetHeader("x-client-version"),
		},
	)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, resp)
}
