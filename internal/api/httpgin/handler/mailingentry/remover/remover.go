package remover

import (
	"context"
	"github.com/GeneralKenobi/mailman/internal/api/httpgin/wrapper"
	"github.com/GeneralKenobi/mailman/internal/persistence"
	"github.com/GeneralKenobi/mailman/internal/service/mailingentry/remover"
	"github.com/GeneralKenobi/mailman/pkg/mdctx"
	"github.com/gin-gonic/gin"
)

func NewHandler(transactioner persistence.Transactioner) *Handler {
	return &Handler{transactioner: transactioner}
}

type Handler struct {
	transactioner persistence.Transactioner
}

func (handler *Handler) HandlerFunc(request *gin.Context) {
	wrapper.ForRequest(request).Handle(func(ctx context.Context) error {
		ctx = mdctx.WithOperationName(ctx, "delete mailing entry with ID")
		return wrapper.WithRequiredIntPathParam(request, "id", func(id int) error {
			return persistence.WithinTransaction(ctx, handler.transactioner, func(transactionalRepository persistence.Repository) error {
				mailingEntryRemover := remover.New(transactionalRepository)
				return mailingEntryRemover.Remove(ctx, id)
			})
		})
	})
}
