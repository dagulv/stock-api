-- +migrate Up
CREATE TABLE IF NOT EXISTS "session" (
    "tokenId" CHAR(20) PRIMARY KEY,
    "userId" CHAR(20) NOT NULL,
    "timeCreated" TIMESTAMPTZ NOT NULL
);
CREATE INDEX IF NOT EXISTS "session_userId" ON "session" ("userId");
-- +migrate Down
DROP TABLE "session";