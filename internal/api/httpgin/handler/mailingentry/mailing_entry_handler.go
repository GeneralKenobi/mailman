package mailingentry

import (
	"context"
	"fmt"
	"github.com/GeneralKenobi/mailman/internal/api/httpgin/wrapper"
	"github.com/GeneralKenobi/mailman/internal/persistence"
	customercreator "github.com/GeneralKenobi/mailman/internal/service/customer/creator"
	mailingentrycreator "github.com/GeneralKenobi/mailman/internal/service/mailingentry/creator"
	"github.com/GeneralKenobi/mailman/internal/service/mailingentry/remover"
	"github.com/GeneralKenobi/mailman/internal/service/mailingentry/sender"
	"github.com/GeneralKenobi/mailman/internal/service/mailingentry/staleremover"
	"github.com/GeneralKenobi/mailman/pkg/api/model"
	"github.com/GeneralKenobi/mailman/pkg/mdctx"
	"github.com/gin-gonic/gin"
)

func NewHandler(transactioner persistence.Transactioner, emailer sender.Emailer) *Handler {
	return &Handler{
		transactioner: transactioner,
		emailer:       emailer,
	}
}

type Handler struct {
	transactioner persistence.Transactioner
	emailer       sender.Emailer
}

func (handler *Handler) CreateHandlerFunc(request *gin.Context) {
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

func (handler *Handler) DeleteHandlerFunc(request *gin.Context) {
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

func (handler *Handler) SendMailingIdHandlerFunc(request *gin.Context) {
	wrapper.ForRequest(request).Handle(func(ctx context.Context) error {
		ctx = mdctx.WithOperationName(ctx, "send mailing entries with mailing ID")

		return wrapper.WithBoundRequestBody(request, func(mailingRequest model.MailingRequestDto) error {
			mdctx.Debugf(ctx, "Sending mailing entries with mailing ID %d", mailingRequest.MailingId)

			// Use a separate transaction for stale entry cleanup because it can be committed even if sending fails later on.
			err := persistence.WithinTransaction(ctx, handler.transactioner, func(transactionalRepository persistence.Repository) error {
				staleEntryRemover := staleremover.New(transactionalRepository)
				return staleEntryRemover.RemoveByMailingId(ctx, mailingRequest.MailingId)
			})
			if err != nil {
				return fmt.Errorf("can't proceed with sending mailing entries with ID %d - error cleaning up stale entries: %w",
					mailingRequest.MailingId, err)
			}

			err = persistence.WithinTransaction(ctx, handler.transactioner, func(transactionalRepository persistence.Repository) error {
				mailer := sender.New(transactionalRepository, handler.emailer)
				return mailer.SendMailingRequest(ctx, mailingRequest)
			})
			if err != nil {
				return fmt.Errorf("error sending mailing entries with ID %d: %w", mailingRequest.MailingId, err)
			}

			return nil
		})
	})
}
