package wrapper

import (
	"github.com/GeneralKenobi/mailman/pkg/util"
	"github.com/gin-gonic/gin"
)

// WithBoundRequestBody binds JSON request body to an instance of T and calls the given function if the binding was successful.
func WithBoundRequestBody[T any](ginCtx *gin.Context, todo func(requestBody T) error) error {
	_, err := WithBoundRequestBodyReturningV(ginCtx, func(requestBody T) (any, error) {
		return nil, todo(requestBody)
	})
	return err
}

// WithBoundRequestBodyReturningV binds JSON request body to an instance of T and calls the given function if the binding was successful.
func WithBoundRequestBodyReturningV[T, V any](request *gin.Context, todo func(requestBody T) (V, error)) (V, error) {
	var requestBody T
	err := request.BindJSON(&requestBody)
	if err != nil {
		return util.ZeroValue[V](), err
	}

	return todo(requestBody)
}
