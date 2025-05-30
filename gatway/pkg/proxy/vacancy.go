package proxy

import (
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

type VacancyProxy struct {
	proxy *httputil.ReverseProxy
}

func NewVacancyProxy(target string) (*VacancyProxy, error) {
	url, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	return &VacancyProxy{
		proxy: httputil.NewSingleHostReverseProxy(url),
	}, nil
}

func (of *VacancyProxy) RegisterRoutes(r *gin.Engine) {
	r.Any("/vacancies/*path", func(ctx *gin.Context) {
		ctx.Request.URL.Path = ctx.Param("path")
		of.proxy.ServeHTTP(ctx.Writer, ctx.Request)
	})
}
