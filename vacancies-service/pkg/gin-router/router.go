package ginrouter

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Storage interface{}

type Handler struct {
	client Storage
}

func NewHandler(s Storage) *Handler {
	return &Handler{client: s}
}

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.New()
	r.GET("/", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, "Hello world!")
	})
	return r
}
