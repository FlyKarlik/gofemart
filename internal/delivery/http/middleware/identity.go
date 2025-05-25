package middleware

import (
	"context"
	"net/http"

	"github.com/FlyKarlik/gofemart/internal/delivery/http/response"
	"github.com/FlyKarlik/gofemart/internal/errs"
	"github.com/FlyKarlik/gofemart/internal/model"
	"github.com/FlyKarlik/gofemart/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (m *Middleware) Identity(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		m.logger.Error("middleware", "Identity", "Failed to get auth header", errs.ErrEmptyAuthHeader)
		response.New[any](c, http.StatusUnauthorized, false, nil, errs.ErrUnauthorized)
		return
	}

	token, err := jwt.GetClearToken(authHeader)
	if err != nil {
		m.logger.Error("middleware", "Identity", "Failed to get validated token", errs.ErrInvalidToken)
		response.New[any](c, http.StatusUnauthorized, false, nil, errs.ErrUnauthorized)
		return
	}

	claims, err := jwt.ParseToken(token, m.cfg.AppGofemart.JWTSecret)
	if err != nil {
		m.logger.Error("middleware", "Identity", "Failed to parse token", errs.ErrInvalidToken)
		response.New[any](c, http.StatusUnauthorized, false, nil, errs.ErrUnauthorized)
		return
	}

	user, err := m.usecase.GetUserByID(c.Request.Context(), uuid.MustParse(claims.UserID))
	if err != nil {
		m.logger.Error("middleware", "Identity", "Failed to get user by id", err)
		response.New[any](c, http.StatusUnauthorized, false, nil, errs.ErrUnauthorized)
		return
	}

	ctx := context.WithValue(c.Request.Context(), model.ContextKeyEnumUserID, *user.ID)
	c.Request = c.Request.WithContext(ctx)
	c.Next()
}
