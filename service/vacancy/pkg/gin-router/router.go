package ginrouter

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/silentnova42/job_vacancy_poster/pkg/structs"
)

type VacancyStorage interface {
	GetAllAvailableVacancy(ctx context.Context) ([]*structs.VacancyGet, error)
	GetVacancyById(ctx context.Context, id uint) (*structs.VacancyGet, error)
	AddVacancy(ctx context.Context, vacancy *structs.VacancyCreate) error
	UpdateVacancyById(ctx context.Context, vacancy *structs.VacancyUpdate, id uint) error
	AddResponseById(ctx context.Context, id uint, email string) error
	CloseVacancyById(ctx context.Context, id uint) error
	GetResponsesByOwnerId(ctx context.Context, id uint) ([]structs.ResponseGet, error)
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
	response := r.Group("/response")
	{
		response.GET("/:id", h.GetResponsesByOwnerId)
	}
	return r
}
