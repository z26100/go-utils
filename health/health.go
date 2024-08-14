package health

import (
	"github.com/gin-gonic/gin"
	"github.com/z26100/log-go"
)

func HealthEndpoint() func(ctx *gin.Context) {
	log.Infoln("Adding health endpoint")
	return func(ctx *gin.Context) {
		_, _ = ctx.Writer.WriteString("OK")
	}
}
