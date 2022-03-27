package sender

import (
	"context"
	"fmt"
	"github.com/GeneralKenobi/mailman/internal/api/httpgin/wrapper"
	"github.com/GeneralKenobi/mailman/internal/persistence"
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

func (handler *Handler) HandlerFunc(request *gin.Context) {
	wrapper.ForRequest(request, func(ctx context.Context) error {
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
				return mailer.Send(ctx, mailingRequest)
			})
			if err != nil {
				return fmt.Errorf("error sending mailing entries with ID %d: %w", mailingRequest.MailingId, err)
			}

			return nil
		})
	}).Do()
}
