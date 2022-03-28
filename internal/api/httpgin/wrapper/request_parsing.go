package wrapper

import (
	"github.com/GeneralKenobi/mailman/internal/api"
	"github.com/GeneralKenobi/mailman/pkg/util"
	"github.com/gin-gonic/gin"
	"strconv"
)

// WithBoundRequestBody binds JSON request body to an instance of T and calls the given function if the binding was successful.
func WithBoundRequestBody[T any](request *gin.Context, todo func(requestBody T) error) error {
	_, err := WithBoundRequestBodyReturningV(request, func(requestBody T) (any, error) {
		return nil, todo(requestBody)
	})
	return err
}

// WithBoundRequestBodyReturningV binds JSON request body to an instance of T and calls the given function if the binding was successful.
func WithBoundRequestBodyReturningV[T, V any](request *gin.Context, todo func(requestBody T) (V, error)) (V, error) {
	var requestBody T
	err := request.ShouldBindJSON(&requestBody)
	if err != nil {
		return util.ZeroValue[V](), api.StatusBadInput.ErrorWithCause(err, "request body is invalid")
	}

	return todo(requestBody)
}

// WithRequiredIntPathParam finds parameter paramName in the path and tries to convert it to an integer. If it succeeds then it calls the
// given function with the converted value.
// For example, paramName for path "/api/customer/:id" is "id".
func WithRequiredIntPathParam(request *gin.Context, paramName string, todo func(param int) error) error {
	_, err := WithRequiredIntPathParamReturningV(request, paramName, func(param int) (any, error) {
		return nil, todo(param)
	})
	return err
}

// WithRequiredIntPathParamReturningV finds parameter paramName in the path and tries to convert it to an integer. If it succeeds then it
// calls the given function with the converted value.
// For example, paramName for path "/api/customer/:id" is "id".
func WithRequiredIntPathParamReturningV[V any](request *gin.Context, paramName string, todo func(param int) (V, error)) (V, error) {
	return WithRequiredPathParamReturningV[V](request, paramName, func(param string) (V, error) {
		paramAsInt, err := strconv.Atoi(param)
		if err != nil {
			return util.ZeroValue[V](), api.StatusBadInput.Error("path parameter %s has to be an integer", paramName)
		}
		return todo(paramAsInt)
	})
}

// WithRequiredPathParam finds parameter paramName in the path and calls the given function with it if it's not an empty string.
// For example, paramName for path "/api/customer/:email" is "email".
func WithRequiredPathParam(request *gin.Context, paramName string, todo func(param string) error) error {
	_, err := WithRequiredPathParamReturningV(request, paramName, func(param string) (any, error) {
		return nil, todo(param)
	})
	return err
}

// WithRequiredPathParamReturningV finds parameter paramName in the path and calls the given function with it if it's not an empty string.
// For example, paramName for path "/api/customer/:email" is "email".
func WithRequiredPathParamReturningV[V any](request *gin.Context, paramName string, todo func(param string) (V, error)) (V, error) {
	param := request.Param(paramName)
	if param == "" {
		return util.ZeroValue[V](), api.StatusBadInput.Error("path parameter %s is required", paramName)
	}

	return todo(param)
}
