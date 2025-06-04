package router

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/silentnova42/job_vacancy_poster/service/profile/pkg/model"
)

var accessKey = os.Getenv("ACCESS_TOKEN")

func (h *Handler) GetProfileByEmailAndPassword(ctx *gin.Context) {
	var (
		customer = &model.LoginRequest{}
		err      error
	)

	if err = bindAndValdate(ctx, customer, h.validate); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	getCustomer, err := h.dbClient.GetCustomerByEmailAndPassword(ctx.Request.Context(), customer)
	if err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, getCustomer)
}

func (h *Handler) GetProfileByEmail(ctx *gin.Context) {
	customer, err := h.dbClient.GetCustomerByEmail(ctx.Request.Context(), ctx.Param("email"))
	if err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, customer)
}

func (h *Handler) AddProfile(ctx *gin.Context) {
	var (
		newCustomer = &model.CreateCustomer{}
		err         error
	)

	if err = bindAndValdate(ctx, newCustomer, h.validate); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	if err = h.dbClient.AddCustomer(ctx.Request.Context(), newCustomer); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.IndentedJSON(http.StatusCreated, newCustomer.Email)
}

func (h *Handler) UpdateProfile(ctx *gin.Context) {
	var (
		updateCustomer = &model.UpdateCustomer{}
		err            error
	)

	if err = bindAndValdate(ctx, updateCustomer, h.validate); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	email, err := h.CheckAccessTokenAndGetEmailFromThere(ctx)
	if err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	if err = h.dbClient.UpdateCustomer(ctx.Request.Context(), updateCustomer, email); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, updateCustomer)
}

func (h *Handler) UpdatePassword(ctx *gin.Context) {
	var (
		passwordUpdate model.PasswordUpdateRequest
		err            error
	)

	if err = ctx.ShouldBindJSON(&passwordUpdate); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	email, err := h.CheckAccessTokenAndGetEmailFromThere(ctx)
	if err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	if err = h.dbClient.UpdatePassword(ctx.Request.Context(), &passwordUpdate, email); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, email)
}

func (h *Handler) DeleteProfileByEmailAndPassword(ctx *gin.Context) {
	var (
		deleteCustomer = &model.PasswordPayload{}
		err            error
	)

	if err = bindAndValdate(ctx, deleteCustomer, h.validate); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	email, err := h.CheckAccessTokenAndGetEmailFromThere(ctx)
	if err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	if err = h.dbClient.DeleteCustomerByEmailAndPassword(ctx.Request.Context(), deleteCustomer, email); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.IndentedJSON(http.StatusNoContent, deleteCustomer)
}

func (h *Handler) CheckAccessTokenAndGetEmailFromThere(ctx *gin.Context) (string, error) {
	authParts := strings.Split(ctx.GetHeader("Authorization"), " ")

	if len(authParts) < 2 && authParts[0] != "Bearer" {
		return "", errors.New("incorrect header")
	}
	authHash := authParts[1]
	accessToken, err := jwt.Parse(authHash, func(t *jwt.Token) (interface{}, error) {
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

func bindAndValdate[T any](ctx *gin.Context, obj *T, validator *validator.Validate) error {
	if err := ctx.ShouldBindJSON(obj); err != nil {
		return err
	}

	return validator.Struct(obj)
}

func abortWithErr(ctx *gin.Context, status int, err error) {
	ctx.AbortWithStatusJSON(status, gin.H{"error": err.Error()})
}
