package router

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type ReverseProxy interface {
	RegisterRoutes(r *gin.Engine)
}

type Handler struct {
	proxy []ReverseProxy
}

func NewHandler(p ...ReverseProxy) (*Handler, error) {
	if len(p) == 0 {
		return nil, errors.New("proxy not found")
	}

	return &Handler{
		proxy: p,
	}, nil
}

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.Default()

	for _, v := range h.proxy {
		v.RegisterRoutes(r)
	}
	return r
}
