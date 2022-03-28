package creator

import (
	"context"
	"github.com/GeneralKenobi/mailman/internal/api/httpgin/wrapper"
	"github.com/GeneralKenobi/mailman/internal/persistence"
	customercreator "github.com/GeneralKenobi/mailman/internal/service/customer/creator"
	mailingentrycreator "github.com/GeneralKenobi/mailman/internal/service/mailingentry/creator"
	"github.com/GeneralKenobi/mailman/pkg/api/model"
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
		ctx = mdctx.WithOperationName(ctx, "create mailing entry")
		return wrapper.WithBoundRequestBody(request, func(mailingEntry model.MailingEntryDto) error {
			return persistence.WithinTransaction(ctx, handler.transactioner, func(transactionalRepository persistence.Repository) error {
				customerCreator := customercreator.New(transactionalRepository)
				mailingEntryCreator := mailingentrycreator.New(transactionalRepository, customerCreator)
				_, err := mailingEntryCreator.CreateFromDto(ctx, mailingEntry)
				return err
			})
		})
	})
}
