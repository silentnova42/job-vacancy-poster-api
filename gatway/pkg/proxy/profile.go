package proxy

import (
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

type ProfileProxy struct {
	proxy *httputil.ReverseProxy
}

func NewProfileProxy(target string) (*ProfileProxy, error) {
	url, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	return &ProfileProxy{
		proxy: httputil.NewSingleHostReverseProxy(url),
	}, nil
}

func (pp ProfileProxy) RegisterRoutes(r *gin.Engine) {
	r.Any("/profiles/*path", func(ctx *gin.Context) {
		ctx.Request.URL.Path = ctx.Param("path")
		pp.proxy.ServeHTTP(ctx.Writer, ctx.Request)
	})
}
