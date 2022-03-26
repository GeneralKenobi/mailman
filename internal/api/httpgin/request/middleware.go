package request

import (
	"github.com/GeneralKenobi/mailman/pkg/mdctx"
	"github.com/gin-gonic/gin"
)

// ContextMiddleware creates and saves in gin's context an MDC-enhanced context for the request.
func ContextMiddleware(request *gin.Context) {
	ctx := mdctx.New()
	if correlationId := request.GetHeader("X-Correlation-ID"); correlationId != "" {
		ctx = mdctx.WithCorrelationId(ctx, correlationId)
	}
	ctx = mdctx.WithRequestMethod(ctx, request.Request.Method)
	ctx = mdctx.WithRequestUri(ctx, request.Request.RequestURI)
	ctx = mdctx.WithClientIp(ctx, request.ClientIP())
	request.Set(requestContextKey, ctx)
	request.Next()
}

func LogRequestProcessingMiddleware(request *gin.Context) {
	ctx := Context(request)
	mdctx.Infof(ctx, "Begin processing")
	request.Next()
	mdctx.Infof(ctx, "End processing")
}
