package email

import "context"

type Service interface {
	Send(ctx context.Context, emailAddress, title, content string) error
}
