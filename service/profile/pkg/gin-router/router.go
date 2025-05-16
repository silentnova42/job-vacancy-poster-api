package ginrouter

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/silentnova42/job_vacancy_poster/pkg/structs"
)

type CustomerProfileStorage interface {
	GetProfileByEmailAndPassword(ctx context.Context, checkCustomer structs.CheckCustomer) (*structs.GetCustomer, error)
	AddProfile(ctx context.Context, customer structs.CreateCustomer) error
	UpdateProfile(ctx context.Context, updateCustomer structs.UpdateCustomer, getCustomer structs.GetCustomer) error
}

type Handler struct {
	client   CustomerProfileStorage
	validate *validator.Validate
}

func NewHandler(db CustomerProfileStorage) *Handler {
	return &Handler{client: db, validate: validator.New()}
}

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.Default()

	profiles := r.Group("/profiles")
	{
		profiles.GET("/", h.GetProfileByEmailAndPassword)
		profiles.POST("/", h.AddProfile)
		profiles.PATCH("/", h.UpdateProfile)
	}

	return r
}
