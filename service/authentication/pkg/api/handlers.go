package router

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/silentnova42/job_vacancy_poster/service/auth/pkg/model"
)

func (h *Handler) Login(ctx *gin.Context) {
	var customer model.Customer
	if err := bindAndValidateCustomer(ctx, &customer, h.validate); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	access, err := h.auth.GenerateAccessToken(&customer)
	log.Println(access)
	if err != nil {
		abortWithErr(ctx, http.StatusUnauthorized, err)
		return
	}

	refresh, err := h.auth.GenerateRefreshToken(&customer)
	log.Println(access)
	if err != nil {
		abortWithErr(ctx, http.StatusUnauthorized, err)
		return
	}

	ctx.SetCookie(
		"refresh_token",
		refresh,
		int(h.expRefresh),
		"/",
		"",
		true,
		true,
	)

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"access_token:": access,
	})
}

func (h *Handler) Refresh(ctx *gin.Context) {
	cookie, err := ctx.Request.Cookie("refresh_token")
	if err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	token := cookie.Value
	if strings.TrimSpace(token) == "" {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	customer, err := h.auth.Refresh(token)
	if err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	access, err := h.auth.GenerateAccessToken(customer)
	if err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	refresh, err := h.auth.GenerateRefreshToken(customer)
	if err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.SetCookie(
		"refresh_token",
		refresh,
		int(h.expRefresh),
		"/",
		"",
		true,
		true,
	)

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"access_token:": access,
	})
}

func (h *Handler) Secure(ctx *gin.Context) {
	access := ctx.GetHeader("Authorization")
	if access == "" {
		abortWithErr(ctx, http.StatusBadRequest, errors.New("incorrect token"))
		return
	}

	parts := strings.Split(access, "Bearer")
	if parts[0] != "" || len(parts) < 2 {
		abortWithErr(ctx, http.StatusUnauthorized, errors.New("incorrect token"))
		return
	}

	customer, err := h.auth.Secure(strings.TrimSpace(parts[1]))
	if err != nil {
		abortWithErr(ctx, http.StatusUnauthorized, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, customer)
}

func bindAndValidateCustomer(ctx *gin.Context, obj interface{}, v *validator.Validate) error {
	var err error
	if err = ctx.ShouldBindJSON(obj); err != nil {
		return err
	}

	if err = v.Struct(obj); err != nil {
		return err
	}

	return nil
}

func abortWithErr(ctx *gin.Context, status int, err error) {
	ctx.AbortWithStatusJSON(status, gin.H{"error": err.Error()})
}
