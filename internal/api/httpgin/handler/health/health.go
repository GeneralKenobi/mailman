package health

import (
	"context"
	"github.com/GeneralKenobi/mailman/internal/api/httpgin/wrapper"
	"github.com/GeneralKenobi/mailman/pkg/mdctx"
	"github.com/gin-gonic/gin"
)

// HandlerFunc always responds with HTTP200 (if the server is up the application is healthy).
func HandlerFunc(request *gin.Context) {
	wrapper.ForRequest(request).Handle(func(ctx context.Context) error {
		mdctx.Debugf(ctx, "Health probe")
		return nil
	})
}
