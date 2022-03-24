package mock

import (
	"context"
	"github.com/GeneralKenobi/mailman/internal/email"
	"github.com/GeneralKenobi/mailman/pkg/mdctx"
)

type EmailService struct{}

var _ email.EmailService = (*EmailService)(nil)

func New() *EmailService {
	return &EmailService{}
}

func (m *EmailService) Send(ctx context.Context, emailAddress, title, content string) error {
	mdctx.Infof(ctx, "Mock email service: Sending email titled %q to %q with content %q", title, emailAddress, content)
	return nil
}
