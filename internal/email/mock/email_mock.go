package mock

import (
	"context"
	"github.com/GeneralKenobi/mailman/internal/email"
	"github.com/GeneralKenobi/mailman/pkg/mdctx"
)

type Emailer struct{}

var _ email.Service = (*Emailer)(nil)

func NewEmailer() *Emailer {
	return &Emailer{}
}

func (emailer *Emailer) Send(ctx context.Context, emailAddress, title, content string) error {
	mdctx.Infof(ctx, "Mock email service: Sending email titled %q to %q with content %q", title, emailAddress, content)
	return nil
}
