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
	row := s.db.QueryRow(
		ctx,
		`SELECT
			"credentials"."userId",
			"credentials"."password",
			"credentials"."otpSecret",
			"credentials"."credentialId",
			"credentials"."publicKey"
		FROM "users"
		INNER JOIN "credentials" ON "credentials"."userId" = "users"."id"
		WHERE "users"."email" = $1`, email,
	)

	if err = row.Scan(&credentials.UserId, &credentials.Password, &credentials.OtpSecret, &credentials.CredentialId, &credentials.PublicKey); err != nil {
		return
	}

	return nil
}

func (s authStore) InsertCredentials(ctx context.Context, creds domain.Credentials) (err error) {
	_, err = s.db.Exec(
		ctx,
		`INSERT INTO "credentials" (
			"userId",
			"password",
			"otpSecret",
			"credentialId",
			"publicKey"
		) VALUES ($1, $2, $3, $4, $5)`, creds.UserId, creds.Password, creds.OtpSecret, creds.CredentialId, creds.PublicKey,
	)

	return
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

func (s authStore) LazyGetSessionUser(ctx context.Context, sessionId xid.ID) (sessionUser *domain.SessionUser, err error) {
	var exists bool

	if sessionUser, exists = s.cache.Get(sessionId); exists {
		return
	}

	return s.getSessionUser(ctx, sessionId)
}

func (s authStore) getSessionUser(ctx context.Context, sessionId xid.ID) (_ *domain.SessionUser, err error) {
	var sessionUser domain.SessionUser

	row := s.db.QueryRow(
		ctx,
		`SELECT
			"users"."id",
			"users"."firstName",
			"users"."lastName",
			"users"."email"
		FROM "session"
		INNER JOIN "users" ON "session"."userId" = "users"."id"
		WHERE "session"."id" = $1`, sessionId,
	)

	if err = row.Scan(&sessionUser.Id, &sessionUser.FirstName, &sessionUser.LastName, &sessionUser.Email); err != nil {
		return
	}

	s.cache.Put(sessionUser.Id, sessionUser)

	return &sessionUser, nil
}

func (s authStore) GetSession(ctx context.Context, sessionId xid.ID, session *domain.Session) (err error) {
	row, err := s.db.Query(
		ctx,
		`SELECT
			"id",
			"userId",
			"scope",
			"timeExpired"
		FROM "session"
		WHERE "id" = $1`, sessionId,
	)

	if err != nil {
		return
	}

	defer row.Close()

	return row.Scan(&session.Id, &session.UserId, &session.Scope, &session.TimeExpired)
}

func (s authStore) InsertSession(ctx context.Context, session domain.Session) (err error) {
	_, err = s.db.Exec(
		ctx,
		`INSERT INTO "session" (
			"id",
			"userId",
			"scope",
			"timeExpired"
		) VALUES ($1, $2, $3, $4)`, session.Id, session.UserId, session.Scope, session.TimeExpired,
	)

	return
}

func (s authStore) DeleteSession(ctx context.Context, sessionId xid.ID) (err error) {
	_, err = s.db.Exec(
		ctx,
		`DELETE FROM "session" WHERE "id" = $1`, sessionId,
	)

	if err != nil {
		return
	}

	s.cache.Delete(sessionId)

	return
}
