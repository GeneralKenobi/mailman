package api

import (
	"context"
	"github.com/GeneralKenobi/mailman/pkg/mdctx"
	"github.com/gin-gonic/gin"
)

// RequestContext extracts context to pass downstream from a gin context.
func RequestContext(ginCtx *gin.Context) context.Context {
	value, found := ginCtx.Get(requestContextKey)
	if !found {
		newCtx := mdctx.New()
		mdctx.Warnf(newCtx, "Missing request context")
		return newCtx
	}

	ctx, ok := value.(context.Context)
	if !ok {
		newCtx := mdctx.New()
		mdctx.Warnf(newCtx, "Unexpected request context type (%T)", value)
		return newCtx
	}

	return ctx
}

// requestContextMiddleware creates and saves in gin's context an MDC-enhanced context for the request.
func requestContextMiddleware(ginCtx *gin.Context) {
	ctx := mdctx.New()
	if correlationId := ginCtx.GetHeader("X-Correlation-ID"); correlationId != "" {
		ctx = mdctx.WithCorrelationId(ctx, correlationId)
	}
	ctx = mdctx.WithRequestMethod(ctx, ginCtx.Request.Method)
	ctx = mdctx.WithRequestUri(ctx, ginCtx.Request.RequestURI)
	ctx = mdctx.WithClientIp(ctx, ginCtx.ClientIP())
	ginCtx.Set(requestContextKey, ctx)
	ginCtx.Next()
}

const requestContextKey = "requestContext"

func logRequestProcessingMiddleware(ginCtx *gin.Context) {
	ctx := RequestContext(ginCtx)
	mdctx.Infof(ctx, "Begin processing")
	ginCtx.Next()
	mdctx.Infof(ctx, "End processing")
}
