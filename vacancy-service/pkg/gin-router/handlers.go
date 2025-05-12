package ginrouter

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/silentnova42/job_vacancy_poster/pkg/structs"
)

func (h *Handler) GetAllAvailableVacancy(ctx *gin.Context) {
	if vacancys, err := h.client.GetAllAvailableVacancy(ctx.Request.Context()); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
	} else {
		ctx.IndentedJSON(http.StatusOK, vacancys)
	}
}

func (h *Handler) GetVacancyById(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	vacancy, err := h.client.GetVacancyById(ctx.Request.Context(), uint(id))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, vacancy)
}

func (h *Handler) AddVacancy(ctx *gin.Context) {
	var (
		vacancy structs.VacancyCreate
		err     error
	)

	if err = ctx.ShouldBindJSON(&vacancy); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	if err = h.validate.Struct(vacancy); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	if err = h.client.AddVacancy(ctx.Request.Context(), &vacancy); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, vacancy)
}

func (h *Handler) UpdateVacancyById(ctx *gin.Context) {
	var (
		newVacancy structs.VacancyUpdate
		err        error
	)

	if err = ctx.ShouldBindJSON(&newVacancy); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	if err = h.validate.Struct(newVacancy); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 0, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	if err = h.client.UpdateVacancyById(ctx.Request.Context(), &newVacancy, uint(id)); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, newVacancy)
}

func (h *Handler) AddResponseById(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	if err = h.client.AddResponseById(ctx.Request.Context(), uint(id)); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, id)
}

func (h *Handler) CloseVacancyById(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	if err = h.client.CloseVacancyById(ctx.Request.Context(), uint(id)); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, id)
}
