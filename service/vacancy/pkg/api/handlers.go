package router

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/silentnova42/job_vacancy_poster/pkg/model"
)

func (h *Handler) GetAllAvailableVacancy(ctx *gin.Context) {
	if vacancys, err := h.client.GetAllAvailableVacancy(ctx.Request.Context()); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
	} else {
		ctx.IndentedJSON(http.StatusOK, vacancys)
	}
}

func (h *Handler) GetVacancyById(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 32)
	if err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	vacancy, err := h.client.GetVacancyById(ctx.Request.Context(), uint(id))
	if err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, vacancy)
}

func (h *Handler) AddVacancy(ctx *gin.Context) {
	var (
		vacancy model.VacancyCreate
		err     error
	)

	if err = bindAndValdate(ctx, &vacancy, h.validate); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
	}

	if err = h.client.AddVacancy(ctx.Request.Context(), &vacancy); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, vacancy)
}

func (h *Handler) UpdateVacancyById(ctx *gin.Context) {
	var (
		newVacancy model.VacancyUpdate
		err        error
	)

	if err = bindAndValdate(ctx, &newVacancy, h.validate); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 0, 32)
	if err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	if err = h.client.UpdateVacancyById(ctx.Request.Context(), &newVacancy, uint(id)); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, newVacancy)
}

func (h *Handler) AddResponseById(ctx *gin.Context) {
	var (
		response model.ResponseCreate
		err      error
	)

	if err = bindAndValdate(ctx, &response, h.validate); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
	}

	if err = h.client.AddResponseById(ctx.Request.Context(), uint(response.VacancyId), response.Email); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, response.Email)
}

func (h *Handler) CloseVacancyById(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 32)
	if err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	if err = h.client.CloseVacancyById(ctx.Request.Context(), uint(id)); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, id)
}

func (h *Handler) GetResponsesByVacancyId(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 32)
	if err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	responses, err := h.client.GetResponsesByVacancyId(ctx, uint(id))
	if err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, responses)
}

func bindAndValdate[T any](ctx *gin.Context, obj *T, validator *validator.Validate) error {
	if err := ctx.ShouldBindJSON(obj); err != nil {
		return err
	}

	return validator.Struct(obj)
}

func abortWithErr(ctx *gin.Context, status int, err error) {
	ctx.AbortWithStatusJSON(status, gin.H{"error": err.Error()})
}
