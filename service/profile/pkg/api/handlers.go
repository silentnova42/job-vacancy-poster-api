package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/silentnova42/job_vacancy_poster/pkg/model"
)

func (h *Handler) GetProfileByEmailAndPassword(ctx *gin.Context) {
	var (
		customer model.Credentials
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

	ctx.IndentedJSON(http.StatusOK, getCustomer)
}

func (h *Handler) GetProfileByEmail(ctx *gin.Context) {
	customer, err := h.client.GetCustomerByEmail(ctx.Request.Context(), ctx.Param("email"))
	if err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, customer)
}

func (h *Handler) AddProfile(ctx *gin.Context) {
	var (
		newCustomer model.CreateCustomer
		err         error
	)

	if err = bindAndValdate(ctx, &newCustomer, h.validate); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	if err = h.client.AddProfile(ctx.Request.Context(), newCustomer); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.IndentedJSON(http.StatusCreated, newCustomer.Email)
}

func (h *Handler) UpdateProfile(ctx *gin.Context) {
	var (
		updateCustomer model.UpdateCustomer
		err            error
	)

	if err = bindAndValdate(ctx, &updateCustomer, h.validate); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	if err = h.client.UpdateProfile(ctx.Request.Context(), updateCustomer); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, updateCustomer)
}

func (h *Handler) DeleteProfileByEmailAndPassword(ctx *gin.Context) {
	var (
		deleteCustomer model.Credentials
		err            error
	)

	if err = bindAndValdate(ctx, &deleteCustomer, h.validate); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	if err = h.client.DeleteProfileByEmailAndPassword(ctx.Request.Context(), deleteCustomer); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.IndentedJSON(http.StatusNoContent, deleteCustomer)
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
