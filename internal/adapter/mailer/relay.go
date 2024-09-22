package mailer

import "context"

type Relay interface {
	Send(ctx context.Context, email Email) (string, error)
}
