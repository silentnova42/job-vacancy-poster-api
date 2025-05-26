package router

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/silentnova42/job_vacancy_poster/pkg/model"
)

type VacancyStorage interface {
	GetAllAvailableVacancy(ctx context.Context) ([]*model.VacancyGet, error)
	GetVacancyById(ctx context.Context, id uint) (*model.VacancyGet, error)
	AddVacancy(ctx context.Context, vacancy *model.VacancyCreate) error
	UpdateVacancyById(ctx context.Context, vacancy *model.VacancyUpdate, id uint) error
	AddResponseById(ctx context.Context, id uint, email string) error
	CloseVacancyById(ctx context.Context, id uint) error
	GetResponsesByVacancyId(ctx context.Context, id uint) ([]model.ResponseGet, error)
}

type Handler struct {
	client   VacancyStorage
	validate *validator.Validate
}

func NewHandler(db VacancyStorage) *Handler {
	return &Handler{
		client:   db,
		validate: validator.New(),
	}
}

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.Default()
	vacancys := r.Group("/vacancies")
	{
		vacancys.GET("/", h.GetAllAvailableVacancy)
		vacancys.GET("/:id", h.GetVacancyById)
		vacancys.POST("/", h.AddVacancy)
		vacancys.PATCH("/:id", h.UpdateVacancyById)
		vacancys.PATCH("/apply/", h.AddResponseById)
		vacancys.DELETE("/:id", h.CloseVacancyById)
	}
	response := r.Group("/responses")
	{
		response.GET("/:id", h.GetResponsesByVacancyId)
	}
	return r
}
