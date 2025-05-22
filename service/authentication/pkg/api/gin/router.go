package ginrouter

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	customer "github.com/silentnova42/job_vacancy_poster/pkg"
)

type Authorization interface {
	FindTokensFromEnv() error
	GenerateAccessToken(customer customer.Customer, exp time.Duration) (string, error)
	GenerateRefreshToken(customer customer.Customer, exp time.Duration) (string, error)
	Secure(accessToken string) (*customer.Customer, error)
	Refresh(refreshToken string) (string, string, error)
}

type Handler struct {
	auth     Authorization
	validate validator.Validate
}

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.Default()
	return r
}
