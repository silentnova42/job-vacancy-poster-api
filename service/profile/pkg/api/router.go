package router

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/silentnova42/job_vacancy_poster/service/profile/pkg/model"
)

type CustomerProfileStorage interface {
	GetProfileByEmailAndPassword(ctx context.Context, checkCustomer model.Credentials) (*model.GetPrivateCustomer, error)
	GetCustomerByEmail(ctx context.Context, email string) (*model.GetPublicCustomer, error)
	AddProfile(ctx context.Context, customer model.CreateCustomer) error
	UpdateProfile(ctx context.Context, updateCustomer model.UpdateCustomer) error
	DeleteProfileByEmailAndPassword(ctx context.Context, check model.Credentials) error
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
	r.GET("/", h.GetProfileByEmailAndPassword)
	r.GET("/:email", h.GetProfileByEmail)
	r.POST("/", h.AddProfile)
	r.PATCH("/", h.UpdateProfile)
	r.DELETE("/", h.DeleteProfileByEmailAndPassword)

	return r
}
