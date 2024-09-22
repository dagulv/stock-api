package mailer

import (
	"context"

	"github.com/mailersend/mailersend-go"
)

type Mailersend struct {
	From mailersend.Recipient
	Ms   *mailersend.Mailersend
}

func (m Mailersend) Send(ctx context.Context, email Email) (string, error) {
	message := m.Ms.Email.NewMessage()
	message.SetFrom(m.From)
	message.SetSubject(email.Subject)
	message.SetRecipients([]mailersend.Recipient{
		{
			Name:  email.FirstName + " " + email.LastName,
			Email: email.Email,
		},
	})
	message.SetHTML(*email.HTMLTemplate)

	resp, err := m.Ms.Email.Send(ctx, message)

	if err != nil {
		return "", err
	}

	return resp.Header.Get("X-Message-Id"), nil
}
