package router

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/silentnova42/job_vacancy_poster/service/auth/pkg/model"
)

func (h *Handler) Login(ctx *gin.Context) {
	var customer model.Credentials
	if err := bindAndValidateCustomer(ctx, &customer, h.validate); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	buf, err := json.Marshal(customer)
	if err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	resp, err := http.Post(h.profileService, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		abortWithErr(ctx, http.StatusBadRequest, errors.New("profile not found"))
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	var getCustomer model.GetCustomer
	if err = json.Unmarshal(data, &getCustomer); err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	access, err := h.auth.GenerateAccessToken(&getCustomer)
	if err != nil {
		abortWithErr(ctx, http.StatusUnauthorized, err)
		return
	}

	refresh, err := h.auth.GenerateRefreshToken(&getCustomer)
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

	newRefreshAndAccessTokens, err := h.auth.Refresh(token)
	if err != nil {
		abortWithErr(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.SetCookie(
		"refresh_token",
		newRefreshAndAccessTokens.RefreshToken,
		int(h.expRefresh),
		"/",
		"",
		true,
		true,
	)

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"access_token:": newRefreshAndAccessTokens.AccessToken,
	})
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
