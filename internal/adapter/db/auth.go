package db

import (
	"context"
	"math"

	"github.com/dagulv/stock-api/internal/adapter/cache"
	"github.com/dagulv/stock-api/internal/core/domain"
	"github.com/dagulv/stock-api/internal/core/port"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"
)

type authStore struct {
	db    *pgxpool.Pool
	cache *cache.Lru[xid.ID, domain.SessionUser]
}

func NewAuth(db *pgxpool.Pool) port.Auth {
	s := authStore{
		db:    db,
		cache: cache.NewLru[xid.ID, domain.SessionUser](math.MaxUint8),
	}

	return s
}

func (s authStore) GetCredentialsByEmail(ctx context.Context, email string, credentials *domain.Credentials) (err error) {
	row, err := s.db.Query(
		ctx,
		`SELECT
			"credentials"."userId",
			"credentials"."password",
			"credentials"."otpSecret",
			"credentials"."credentialId",
			"credentials"."publicKey"
		FROM "user"
		LEFT JOIN "credentials" ON "credentials"."userId" = "user"."id"
		WHERE "user"."email" = $1`, email,
	)

	if err != nil {
		return
	}

	defer row.Close()

	return row.Scan(&credentials.UserId, &credentials.Password, &credentials.OtpSecret, &credentials.CredentialId, &credentials.PublicKey)
}

func (s authStore) UpdatePassword(ctx context.Context, userId xid.ID, hashedPassword []byte) (err error) {
	_, err = s.db.Exec(
		ctx,
		`UPDATE "credentials"
		SET "password" = $1
		WHERE "userId" = $2`,
		hashedPassword, userId,
	)

	return
}

func (s authStore) LazyGetSessionUser(ctx context.Context, sessionUserId xid.ID) (sessionUser *domain.SessionUser, err error) {
	var exists bool

	if sessionUser, exists = s.cache.Get(sessionUserId); exists {
		return
	}

	if err = s.getSessionUser(ctx, sessionUserId, sessionUser); err != nil {
		return nil, err
	}

	return
}

func (s authStore) getSessionUser(ctx context.Context, sessionUserId xid.ID, sessionUser *domain.SessionUser) (err error) {
	row, err := s.db.Query(
		ctx,
		`SELECT
			"id",
			"tenantId",
			"firstName",
			"lastName",
			"email"
		FROM "user"
		WHERE "id" = $1`, sessionUserId,
	)

	if err != nil {
		return
	}

	defer row.Close()

	if err = row.Scan(&sessionUser.Id, &sessionUser.TenantId, &sessionUser.FirstName, &sessionUser.LastName, &sessionUser.Email); err != nil {
		return
	}

	s.cache.Put(sessionUser.Id, *sessionUser)

	return
}

func (s authStore) GetSession(ctx context.Context, sessionId xid.ID, session *domain.Session) (err error) {
	row, err := s.db.Query(
		ctx,
		`SELECT
			"id",
			"userId",
			"timeExpired"
		FROM "session"
		WHERE "id" = $1`, sessionId,
	)

	if err != nil {
		return
	}

	defer row.Close()

	return row.Scan(&session.Id, &session.UserId, &session.TimeExpired)
}

func (s authStore) InsertSession(ctx context.Context, session domain.Session) (err error) {
	_, err = s.db.Exec(
		ctx,
		`INSERT INTO "session" (
			"id",
			"userId",
			"timeExpired"
		) VALUES ($1, $2, $3)`, session.Id, session.UserId, session.TimeExpired,
	)

	return
}

func (s authStore) DeleteSession(ctx context.Context, sessionId xid.ID) (err error) {
	_, err = s.db.Exec(
		ctx,
		`DELETE FROM "session" WHERE "id" = $1`, sessionId,
	)

	return
}
