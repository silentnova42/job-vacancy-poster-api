package router

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/silentnova42/job_vacancy_poster/service/vacancy/pkg/model"
)

var accessKey = os.Getenv("ACCESS_TOKEN")

func (h *Handler) GetAllAvailableVacancy(ctx *gin.Context) {
	if vacancys, err := h.client.GetAllAvailableVacancy(ctx.Request.Context()); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
	} else {
		ctx.IndentedJSON(http.StatusOK, vacancys)
	}
}

func (h *Handler) GetVacancyById(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 64)
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

	email, err := h.CheckAccessTokenAndGetEmailFromThere(ctx)
	if err != nil {
		abortWithErr(ctx, http.StatusUnauthorized, err)
		return
	}

	if err = bindAndValidate(ctx, &vacancy, h.validate); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
	}

	if err = h.client.AddVacancy(ctx.Request.Context(), &vacancy, email); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, vacancy)
}

func (h *Handler) UpdateVacancyByIdAndEmail(ctx *gin.Context) {
	var (
		newVacancy model.VacancyUpdate
		err        error
	)

	email, err := h.CheckAccessTokenAndGetEmailFromThere(ctx)
	if err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	if err = bindAndValidate(ctx, &newVacancy, h.validate); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 0, 64)
	if err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	if err = h.client.UpdateVacancyByIdAndEmail(ctx.Request.Context(), &newVacancy, uint(id), email); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, id)
}

func (h *Handler) CloseVacancyByIdAndEmail(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 64)
	if err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	email, err := h.CheckAccessTokenAndGetEmailFromThere(ctx)
	if err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	if err = h.client.CloseVacancyByIdAndEmail(ctx.Request.Context(), uint(id), email); err != nil {
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

func (h *Handler) AddResponseByIdAndEmail(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 64)
	if err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	email, err := h.CheckAccessTokenAndGetEmailFromThere(ctx)
	if err != nil {
		abortWithErr(ctx, http.StatusUnauthorized, err)
		return
	}

	if err = h.client.AddResponseByIdAndEmail(ctx.Request.Context(), uint(id), email); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, id)
}

func (h *Handler) DeleteResponseByIdAndEmail(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 64)
	if err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	email, err := h.CheckAccessTokenAndGetEmailFromThere(ctx)
	if err != nil {
		abortWithErr(ctx, http.StatusUnauthorized, err)
		return
	}

	if err = h.client.DeleteResponseByIdAndEmail(ctx.Request.Context(), uint(id), email); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.IndentedJSON(http.StatusNoContent, id)
}

func (h *Handler) CheckAccessTokenAndGetEmailFromThere(ctx *gin.Context) (string, error) {
	authParts := strings.Split(ctx.GetHeader("Authorization"), " ")

	if len(authParts) < 2 || authParts[0] != "Bearer" {
		return "", errors.New("incorrect header")
	}

	accessHash := strings.TrimSpace(authParts[1])
	accessToken, err := jwt.Parse(accessHash, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("signature error")
		}

		return []byte(accessKey), nil
	})

	if err != nil {
		return "", err
	}

	if !accessToken.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := accessToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("cannot parse claims")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", errors.New("email not found")
	}

	return email, nil
}

func bindAndValidate[T any](ctx *gin.Context, obj *T, validator *validator.Validate) error {
	if err := ctx.ShouldBindJSON(obj); err != nil {
		return err
	}

	return validator.Struct(obj)
}

func abortWithErr(ctx *gin.Context, status int, err error) {
	ctx.AbortWithStatusJSON(status, gin.H{"error": err.Error()})
}
