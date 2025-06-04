package router

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/silentnova42/job_vacancy_poster/service/profile/pkg/model"
)

type CustomerProfileStorage interface {
	GetCustomerByEmailAndPassword(ctx context.Context, credentials *model.LoginRequest) (*model.GetPrivateCustomer, error)
	GetCustomerByEmail(ctx context.Context, email string) (*model.GetPublicCustomer, error)
	AddCustomer(ctx context.Context, newCustomer *model.CreateCustomer) error
	UpdateCustomer(ctx context.Context, updateCustomer *model.UpdateCustomer, email string) error
	UpdatePassword(ctx context.Context, passwordUpdate *model.PasswordUpdateRequest, email string) error
	DeleteCustomerByEmailAndPassword(ctx context.Context, credentials *model.PasswordPayload, email string) error
}

type Handler struct {
	dbClient CustomerProfileStorage
	validate *validator.Validate
}

func NewHandler(db CustomerProfileStorage) *Handler {
	return &Handler{dbClient: db, validate: validator.New()}
}

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/:email", h.GetProfileByEmail)
	r.POST("/", h.GetProfileByEmailAndPassword)
	r.POST("/reg/", h.AddProfile)
	r.PATCH("/", h.UpdateProfile)
	r.PATCH("/password", h.UpdatePassword)
	r.DELETE("/", h.DeleteProfileByEmailAndPassword)

	return r
}
