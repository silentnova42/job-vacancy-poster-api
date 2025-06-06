package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/silentnova42/job_vacancy_poster/service/auth/pkg/model"
)

type Authorization interface {
	GenerateAccessToken(customer *model.GetCustomer) (string, error)
	GenerateRefreshToken(customer *model.GetCustomer) (string, error)
	Refresh(refreshToken string) (*model.TokenPair, error)
}

type Handler struct {
	auth           Authorization
	validate       *validator.Validate
	expRefresh     int64
	profileService string
}

func NewHandler(auth Authorization, expRefresh int64, profileService string) (*Handler, error) {
	return &Handler{
		auth:           auth,
		validate:       validator.New(),
		expRefresh:     expRefresh,
		profileService: profileService,
	}, nil
}

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/login", h.Login)
	r.POST("/logout", h.Logout)
	r.POST("/refresh", h.Refresh)
	return r
}
