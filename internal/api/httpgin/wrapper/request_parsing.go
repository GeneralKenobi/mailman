package wrapper

import (
	"github.com/GeneralKenobi/mailman/internal/api"
	"github.com/GeneralKenobi/mailman/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"strconv"
	"strings"
)

// WithBoundRequestBody binds JSON request body to an instance of T and validates it. If both operations were successful calls the
// given function.
func WithBoundRequestBody[T any](request *gin.Context, todo func(requestBody T) error) error {
	_, err := WithBoundRequestBodyReturningV(request, func(requestBody T) (any, error) {
		return nil, todo(requestBody)
	})
	return err
}

// WithBoundRequestBodyReturningV binds JSON request body to an instance of T and validates it. If both operations were successful calls the
// given function.
func WithBoundRequestBodyReturningV[T, V any](request *gin.Context, todo func(requestBody T) (V, error)) (V, error) {
	var requestBody T
	err := request.ShouldBindJSON(&requestBody)
	if err != nil {
		return util.ZeroValue[V](), api.StatusBadInput.WithMessageAndCause(err, "malformed request body")
	}
	err = validateRequestBody(requestBody)
	if err != nil {
		return util.ZeroValue[V](), err
	}

	return todo(requestBody)
}

// validateRequestBody validates a request body. If there were validation errors it converts them into a api.StatusBadInput error.
func validateRequestBody(toValidate any) error {
	err := validate.Struct(toValidate)
	if err == nil {
		return nil
	}

	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return api.StatusBadInput.WithMessageAndCause(err, "invalid request body")
	}

	validationMessages := make([]string, len(validationErrs))
	for i, validationErr := range validationErrs {
		validationMessages[i] = validationErr.Namespace() + ": " + validationErr.Tag()
	}
	return api.StatusBadInput.WithMessageAndCause(err, "invalid request body: %s", strings.Join(validationMessages, ", "))
}

// Use a single instance of Validate, it caches struct info.
var validate = validator.New()

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
			return util.ZeroValue[V](), api.StatusBadInput.WithMessage("path parameter %s has to be an integer", paramName)
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
		return util.ZeroValue[V](), api.StatusBadInput.WithMessage("path parameter %s is required", paramName)
	}

	return todo(param)
}
