package proxy

import (
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

type AuthProxy struct {
	proxy *httputil.ReverseProxy
}

func NewAuthProxy(target string) (*AuthProxy, error) {
	url, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	return &AuthProxy{
		proxy: httputil.NewSingleHostReverseProxy(url),
	}, nil
}

func (ap *AuthProxy) RegisterRoutes(r *gin.Engine) {
	r.Any("/auth/*path", func(ctx *gin.Context) {
		ctx.Request.URL.Path = ctx.Param("path")
		ap.proxy.ServeHTTP(ctx.Writer, ctx.Request)
	})
}
