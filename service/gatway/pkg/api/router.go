package router

import (
	"github.com/gin-gonic/gin"
)

type ReverseProxy interface {
	RegisterRoutes(r *gin.Engine)
}

type Handler struct {
	proxy []ReverseProxy
}

func NewHandler(p ...ReverseProxy) *Handler {
	return &Handler{
		proxy: p,
	}
}

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.Default()
	for _, v := range h.proxy {
		v.RegisterRoutes(r)
	}
	return r
}
