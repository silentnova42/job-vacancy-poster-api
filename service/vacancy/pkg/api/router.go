package router

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/silentnova42/job_vacancy_poster/service/vacancy/pkg/model"
)

type OfferStorage interface {
	GetAllAvailableVacancy(ctx context.Context) ([]*model.VacancyGet, error)
	GetVacancyById(ctx context.Context, vacancyId uint) (*model.VacancyGetWithResponses, error)
	AddVacancy(ctx context.Context, vacancy *model.VacancyCreate, email string) error
	UpdateVacancyByIdAndEmail(ctx context.Context, vacancy *model.VacancyUpdate, vacancyId uint, email string) error
	CloseVacancyByIdAndEmail(ctx context.Context, id uint, email string) error
	GetResponsesByVacancyId(ctx context.Context, id uint) ([]model.ResponseGet, error)
	AddResponseByIdAndEmail(ctx context.Context, id uint, email string) error
	DeleteResponseByIdAndEmail(ctx context.Context, vacancyId uint, email string) error
}

type Handler struct {
	client   OfferStorage
	validate *validator.Validate
}

func NewHandler(db OfferStorage) *Handler {
	return &Handler{
		client:   db,
		validate: validator.New(),
	}
}

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", h.GetAllAvailableVacancy)
	r.GET("/:id", h.GetVacancyById)
	r.POST("/", h.AddVacancy)
	r.PATCH("/:id", h.UpdateVacancyByIdAndEmail)
	r.DELETE("/:id", h.CloseVacancyByIdAndEmail)

	responses := r.Group("/responses")
	{
		responses.GET("/:id", h.GetResponsesByVacancyId)
		responses.PATCH("/apply/:id", h.AddResponseByIdAndEmail)
		responses.DELETE("/disapply/:id", h.DeleteResponseByIdAndEmail)
	}

	return r
}
