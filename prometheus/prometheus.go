package prometheus

import (
	"github.com/gin-gonic/gin"
	"github.com/z26100/log-go"
)
import "github.com/prometheus/client_golang/prometheus/promhttp"

var handler = promhttp.Handler()

func PrometheusEndpoint() func(ctx *gin.Context) {
	log.Infoln("Adding prometheus endpoint")
	return func(ctx *gin.Context) {
		handler.ServeHTTP(ctx.Writer, ctx.Request)
	}
}
