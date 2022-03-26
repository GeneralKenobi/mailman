package email

import "context"

type EmailService interface {
	Send(ctx context.Context, emailAddress, title, content string) error
}