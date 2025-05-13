package ginrouter

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomerStorage interface {
}

type Handler struct {
	client CustomerStorage
}

func NewHandler(db CustomerStorage) *Handler {
	return &Handler{client: db}
}

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.Default()
	return r
}

func (h *Handler) Home(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, "Hello word!")
}
