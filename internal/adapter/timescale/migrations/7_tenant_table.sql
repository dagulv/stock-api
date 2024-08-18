-- +migrate Up
-- +migrate StatementBegin
-- DO $$
-- BEGIN
--     IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'id_and_name') THEN
--         CREATE TYPE ID_AND_NAME AS (
-- 			"id" VARCHAR(20),
-- 			"name" TEXT
-- 		);
--     END IF;
-- END$$;
-- +migrate StatementEnd

CREATE TABLE IF NOT EXISTS "tenant" (
    "id" VARCHAR(20) PRIMARY KEY,
    "name" TEXT NOT NULL,
    "domain" TEXT NOT NULL,
    "active" BOOLEAN NOT NULL,
    "timeCreated" TIMESTAMP WITH TIME ZONE NOT NULL,
    "timeUpdated" TIMESTAMP WITH TIME ZONE NOT NULL,
    "ts" tsvector GENERATED ALWAYS AS (to_tsvector('custom_english', COALESCE("name", ''))) STORED
);
CREATE INDEX IF NOT EXISTS "tenant_domain" ON "tenant" ("domain");
CREATE INDEX IF NOT EXISTS "tenant_active" ON "tenant" ("active");
CREATE INDEX IF NOT EXISTS "tenant_ts" ON "tenant" USING GIN ("ts");

CREATE TABLE IF NOT EXISTS "user" (
    "id" VARCHAR(20) PRIMARY KEY,
    "tenantId" VARCHAR(20) NOT NULL,
    "firstName" TEXT NOT NULL,
    "lastName" TEXT,
    "email" TEXT NOT NULL,
    "active" BOOLEAN NOT NULL,
    "timeCreated" TIMESTAMP WITH TIME ZONE NOT NULL,
    "timeUpdated" TIMESTAMP WITH TIME ZONE NOT NULL,
    "ts" tsvector GENERATED ALWAYS AS (to_tsvector('custom_english', COALESCE("firstName", '') || ' ' || COALESCE("lastName", '') || ' ' || COALESCE("email", ''))) STORED
);
CREATE INDEX IF NOT EXISTS "user_email" ON "user" ("email");
CREATE INDEX IF NOT EXISTS "user_active" ON "user" ("active");
CREATE INDEX IF NOT EXISTS "user_ts" ON "user" USING GIN ("ts");

CREATE TABLE IF NOT EXISTS "credentials" (
    "userId" VARCHAR(20) PRIMARY KEY,
    "password" TEXT,
    "otpSecret" TEXT,
    "credentialId" VARCHAR(1023) NOT NULL UNIQUE,
    "publicKey" TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS "credentials_credentialId" ON "credentials" ("credentialId");

-- +migrate Down
DROP TABLE "user";
DROP TABLE "tenant";
DROP TABLE "credentials";