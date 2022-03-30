package request

import (
	"context"
	"errors"
	"github.com/GeneralKenobi/mailman/internal/api"
	"github.com/GeneralKenobi/mailman/pkg/api/apimodel"
	"github.com/GeneralKenobi/mailman/pkg/mdctx"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Context extracts context to pass downstream from a gin context.
func Context(request *gin.Context) context.Context {
	value, found := request.Get(requestContextKey)
	if !found {
		newCtx := mdctx.New()
		mdctx.Warnf(newCtx, "Missing request context - creating an empty one")
		return newCtx
	}

	ctx, ok := value.(context.Context)
	if !ok {
		newCtx := mdctx.New()
		mdctx.Warnf(newCtx, "Unexpected request context type (%T) - creating an empty one", value)
		return newCtx
	}

	return ctx
}

const requestContextKey = "requestContext"

// WriteErrorResponse looks for a wrapped api.StatusError in err, if it's found it writes the response based on it, if it's not found then
// it writes a generic internal server error response. The error is also logged.
func WriteErrorResponse(ctx context.Context, request *gin.Context, err error) {
	var apiError api.StatusError
	if !errors.As(err, &apiError) {
		apiError = api.StatusInternalError.WithMessageAndCause(err, "Request processing failed")
	}

	mdctx.Errorf(ctx, "Error processing request: %v", apiError)

	errorDto := apimodel.Error{
		Status:      apiStatusToHttpStatus(apiError.Status()),
		Message:     apiError.Message(),
		OperationId: mdctx.OperationId(ctx),
	}
	request.JSON(errorDto.Status, errorDto)
}

func apiStatusToHttpStatus(status api.Status) int {
	switch status {
	case api.StatusBadInput:
		return http.StatusBadRequest
	case api.StatusUnauthorized:
		return http.StatusUnauthorized
	case api.StatusNotFound:
		return http.StatusNotFound
	case api.StatusInternalError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
