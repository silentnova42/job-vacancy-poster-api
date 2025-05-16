package ginrouter

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/silentnova42/job_vacancy_poster/pkg/structs"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) GetProfileByEmailAndPassword(ctx *gin.Context) {
	var (
		customer structs.CheckCustomer
		err      error
	)

	if err = bindAndValdate(ctx, &customer, h.validate); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	getCustomer, err := h.client.GetProfileByEmailAndPassword(ctx.Request.Context(), customer)
	if err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	if getCustomer == nil {
		abortWithErr(ctx, http.StatusBadRequest, errors.New("could not get the customer"))
		return
	}

	if err = comperePasswordHash([]byte(getCustomer.Password), []byte(customer.Password)); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, getCustomer)
}

func (h *Handler) AddProfile(ctx *gin.Context) {
	var (
		customer structs.CreateCustomer
		err      error
	)

	if err = bindAndValdate(ctx, &customer, h.validate); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	customer.Password, err = getPasswordHash(customer.Password)
	if err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	if err = h.client.AddProfile(ctx.Request.Context(), customer); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, customer.Email)
}

func (h *Handler) UpdateProfile(ctx *gin.Context) {
	var (
		updateCustomer structs.UpdateCustomer
		err            error
	)

	if err = bindAndValdate(ctx, &updateCustomer, h.validate); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	getCustomer, err := h.client.GetProfileByEmailAndPassword(ctx, updateCustomer.Check)
	if err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	if err = h.client.UpdateProfile(ctx.Request.Context(), updateCustomer, *getCustomer); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, updateCustomer)
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

func getPasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(hash), err
}

func comperePasswordHash(hash []byte, password []byte) error {
	return bcrypt.CompareHashAndPassword(hash, password)
}
