package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/silentnova42/job_vacancy_poster/pkg/model"
)

type Authorization interface {
	GenerateAccessToken(customer *model.Customer) (string, error)
	GenerateRefreshToken(customer *model.Customer) (string, error)
	Secure(accessToken string) (*model.Customer, error)
	Refresh(refreshToken string) (*model.Customer, error)
}

type Handler struct {
	auth       Authorization
	validate   *validator.Validate
	expRefresh time.Duration
}

func NewHandler(auth Authorization, expRefresh time.Duration) (*Handler, error) {
	return &Handler{
		auth:       auth,
		validate:   validator.New(),
		expRefresh: expRefresh,
	}, nil
}

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/login", h.Login)
	r.POST("/refresh", h.Refresh)
	r.POST("/me", h.Secure)
	return r
}
