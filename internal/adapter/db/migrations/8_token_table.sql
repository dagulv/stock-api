-- +migrate Up
CREATE TABLE IF NOT EXISTS "session" (
    "id" CHAR(20) PRIMARY KEY,
    "userId" CHAR(20) NOT NULL,
    "timeExpired" TIMESTAMP NOT NULL
);
CREATE INDEX IF NOT EXISTS "session_userId" ON "session" ("userId");
-- +migrate Down
DROP TABLE "session";