package wrapper

import (
	"context"
	"github.com/GeneralKenobi/mailman/internal/api/httpgin/request"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SimpleHandler interface {
	OnSuccess(onSuccess func(ctx context.Context)) SimpleHandler
	OnError(onError func(ctx context.Context, handlerErr error)) SimpleHandler
	Handle(handler func(ctx context.Context) error)
}

// ForRequest creates a convenience standard handler which writes status HTTP200 on success and writes an error response on error.
// SimpleHandler.OnSuccess and SimpleHandler.OnError can be used to override the default success/error handlers.
// The request is processed by calling SimpleHandler.Handle which also accepts the handler function.
func ForRequest(requestCtx *gin.Context) SimpleHandler {
	return &simpleHandler{
		request: requestCtx,
		onSuccess: func(_ context.Context) {
			requestCtx.Status(http.StatusOK)
		},
		onError: func(ctx context.Context, handlerErr error) {
			request.WriteErrorResponse(ctx, requestCtx, handlerErr)
		},
	}
}

type SimpleHandlerRetV[V any] interface {
	OnSuccess(onSuccess func(ctx context.Context, handlerResult V)) SimpleHandlerRetV[V]
	OnError(onError func(ctx context.Context, handlerErr error)) SimpleHandlerRetV[V]
	Handle(handler func(ctx context.Context) (V, error))
}

// ForRequestRetV creates a convenience standard handler which writes status HTTP200 on success and writes an error response on error.
// SimpleHandler.OnSuccess and SimpleHandler.OnError can be used to override the default success/error handlers.
// The request is processed by calling SimpleHandler.Handle which also accepts the handler function.
func ForRequestRetV[V any](requestCtx *gin.Context) SimpleHandlerRetV[V] {
	return &simpleHandlerRetV[V]{
		request: requestCtx,
		onSuccess: func(_ context.Context, handlerResult V) {
			requestCtx.JSON(http.StatusOK, handlerResult)
		},
		onError: func(ctx context.Context, handlerErr error) {
			request.WriteErrorResponse(ctx, requestCtx, handlerErr)
		},
	}
}

type simpleHandlerRetV[V any] struct {
	request   *gin.Context
	onSuccess func(ctx context.Context, handlerResult V)
	onError   func(ctx context.Context, handlerErr error)
}

var _ SimpleHandlerRetV[any] = (*simpleHandlerRetV[any])(nil) // Interface guard

func (handler *simpleHandlerRetV[V]) OnSuccess(onSuccess func(ctx context.Context, handlerResult V)) SimpleHandlerRetV[V] {
	handler.onSuccess = onSuccess
	return handler
}

func (handler *simpleHandlerRetV[V]) OnError(onError func(ctx context.Context, handlerErr error)) SimpleHandlerRetV[V] {
	handler.onError = onError
	return handler
}

func (handler *simpleHandlerRetV[V]) Handle(handlerFunc func(ctx context.Context) (V, error)) {
	ctx := request.Context(handler.request)
	value, err := handlerFunc(ctx)
	if err == nil {
		handler.onSuccess(ctx, value)
	} else {
		handler.onError(ctx, err)
	}
}

type simpleHandler struct {
	request   *gin.Context
	onSuccess func(ctx context.Context)
	onError   func(ctx context.Context, handlerErr error)
}

var _ SimpleHandler = (*simpleHandler)(nil) // Interface guard

func (handler *simpleHandler) OnSuccess(onSuccess func(ctx context.Context)) SimpleHandler {
	handler.onSuccess = onSuccess
	return handler
}

func (handler *simpleHandler) OnError(onError func(ctx context.Context, handlerErr error)) SimpleHandler {
	handler.onError = onError
	return handler
}

func (handler *simpleHandler) Handle(handlerFunc func(ctx context.Context) error) {
	ctx := request.Context(handler.request)
	err := handlerFunc(ctx)
	if err == nil {
		handler.onSuccess(ctx)
	} else {
		handler.onError(ctx, err)
	}
}
