package ginrouter

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type (
	Storage interface{}
)

type Handler struct {
	client   Storage
	validate *validator.Validate
}

func NewHandler(db Storage) *Handler {
	return &Handler{
		client:   db,
		validate: validator.New(),
	}
}

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, "Hello world!")
	})
	return r
}
