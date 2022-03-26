package wrapper

import (
	"context"
	"github.com/GeneralKenobi/mailman/internal/api/httpgin/request"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SimpleHandler interface {
	OnSuccess(onSuccess func()) SimpleHandler
	OnError(onError func(error)) SimpleHandler
	Do()
}

// ForRequest is a convenience standard handler. Handler func is called and, depending on its return value, on success or on error func is
// called.
//
// Default on success writes status HTTP200.
// Default on error writes error response from the error returned from the handler.
//
// To run the handler Do has to be called.
func ForRequest(request *gin.Context, handler func(ctx context.Context) error) SimpleHandler {
	return &simpleHandler{
		handlerReturningNil: ForRequestReturningV(request, func(ctx context.Context) (any, error) {
			err := handler(ctx)
			return nil, err
		}),
	}
}

type SimpleHandlerReturningV[V any] interface {
	OnSuccess(onSuccessFunc func(V)) SimpleHandlerReturningV[V]
	OnError(onError func(error)) SimpleHandlerReturningV[V]
	Do()
}

// ForRequestReturningV is a convenience standard handler. Handler func is called and, depending on its return values, on success or on
// error func is called.
//
// Default on success converts V to JSON and writes it together with status HTTP200.
// Default on error writes error response from the error returned from the handler.
//
// To run the handler Do has to be called.
func ForRequestReturningV[V any](requestCtx *gin.Context, handler func(ctx context.Context) (V, error)) SimpleHandlerReturningV[V] {
	ctx := request.Context(requestCtx)
	return &simpleHandlerReturningV[V]{
		handler: func() (V, error) {
			return handler(ctx)
		},
		onSuccess: func(result V) {
			requestCtx.JSON(http.StatusOK, result)
		},
		onError: func(err error) {
			request.WriteErrorResponse(ctx, requestCtx, err)
		},
	}
}

type simpleHandlerReturningV[V any] struct {
	handler   func() (V, error)
	onSuccess func(V)
	onError   func(error)
}

var _ SimpleHandlerReturningV[any] = (*simpleHandlerReturningV[any])(nil) // Interface guard

func (handler *simpleHandlerReturningV[V]) OnSuccess(onSuccess func(V)) SimpleHandlerReturningV[V] {
	handler.onSuccess = onSuccess
	return handler
}

func (handler *simpleHandlerReturningV[V]) OnError(onError func(error)) SimpleHandlerReturningV[V] {
	handler.onError = onError
	return handler
}

func (handler *simpleHandlerReturningV[V]) Do() {
	value, err := handler.handler()
	if err == nil {
		handler.onSuccess(value)
	} else {
		handler.onError(err)
	}
}

type simpleHandler struct {
	handlerReturningNil SimpleHandlerReturningV[any]
}

var _ SimpleHandler = (*simpleHandler)(nil) // Interface guard

func (handler *simpleHandler) OnSuccess(onSuccess func()) SimpleHandler {
	handler.handlerReturningNil.OnSuccess(func(any) {
		onSuccess()
	})
	return handler
}

func (handler *simpleHandler) OnError(onError func(error)) SimpleHandler {
	handler.handlerReturningNil.OnError(onError)
	return handler
}

func (handler *simpleHandler) Do() {
	handler.handlerReturningNil.Do()
}
