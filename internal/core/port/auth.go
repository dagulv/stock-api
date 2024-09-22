package port

import (
	"context"

	"github.com/dagulv/stock-api/internal/core/domain"
	"github.com/rs/xid"
)

type Auth interface {
	GetCredentialsByEmail(ctx context.Context, email string) (credentials *domain.Credentials, err error)
	UpdatePassword(ctx context.Context, userId xid.ID, hashedPassword []byte) error

	LazyGetSessionUser(ctx context.Context, sessionId xid.ID) (sessionUser *domain.SessionUser, err error)

	GetSession(ctx context.Context, sessionId xid.ID, session *domain.Session) error
	InsertSession(ctx context.Context, session domain.Session) error
	DeleteSession(ctx context.Context, sessionId xid.ID) error
}
