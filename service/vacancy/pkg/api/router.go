package router

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/silentnova42/job_vacancy_poster/service/vacancy/pkg/model"
)

type OfferStorage interface {
	GetAllAvailableVacancy(ctx context.Context) ([]*model.VacancyGet, error)
	GetVacancyById(ctx context.Context, id uint) (*model.VacancyGet, error)
	AddVacancy(ctx context.Context, vacancy *model.VacancyCreate) error
	UpdateVacancyById(ctx context.Context, vacancy *model.VacancyUpdate, id uint) error
	AddResponseById(ctx context.Context, id uint, email string) error
	CloseVacancyById(ctx context.Context, id uint) error
	GetResponsesByVacancyId(ctx context.Context, id uint) ([]model.ResponseGet, error)
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
	r.PATCH("/:id", h.UpdateVacancyById)
	r.DELETE("/:id", h.CloseVacancyById)

	response := r.Group("/responses")
	{
		response.PATCH("/apply/", h.AddResponseById)
		response.GET("/:id", h.GetResponsesByVacancyId)
	}

	return r
}
