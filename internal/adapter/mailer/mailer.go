package mailer

import "context"

type Mailer struct {
	ctx   context.Context
	relay Relay
	mail  chan Email
}

func New(ctx context.Context, relay Relay, bufferSize int) *Mailer {
	m := &Mailer{
		ctx:   ctx,
		relay: relay,
		mail:  make(chan Email, bufferSize),
	}

	go func() {
		for {
			select {
			case email, ok := <-m.mail:
				if !ok {
					return
				}

				m.relay.Send(ctx, email)

			case <-m.ctx.Done():
				close(m.mail)
			}
		}
	}()

	return m
}

func (m *Mailer) Send(email Email) {
	m.mail <- email
}
