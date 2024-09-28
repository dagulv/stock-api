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

CREATE TABLE IF NOT EXISTS "groups" (
    "id" VARCHAR(20) PRIMARY KEY,
    "name" TEXT NOT NULL,
    "active" BOOLEAN NOT NULL,
    "timeCreated" TIMESTAMP WITH TIME ZONE NOT NULL,
    "timeUpdated" TIMESTAMP WITH TIME ZONE NOT NULL,
    "ts" tsvector GENERATED ALWAYS AS (to_tsvector('custom_english', COALESCE("name", ''))) STORED
);
CREATE INDEX IF NOT EXISTS "groups_active" ON "groups" ("active");
CREATE INDEX IF NOT EXISTS "groups_ts" ON "groups" USING GIN ("ts");

CREATE TABLE IF NOT EXISTS "users" (
    "id" VARCHAR(20) PRIMARY KEY,
    "firstName" TEXT,
    "lastName" TEXT,
    "email" TEXT UNIQUE NOT NULL,
    "timeCreated" TIMESTAMP WITH TIME ZONE NOT NULL,
    "timeUpdated" TIMESTAMP WITH TIME ZONE NOT NULL,
    "ts" tsvector GENERATED ALWAYS AS (to_tsvector('custom_english', COALESCE("firstName", '') || ' ' || COALESCE("lastName", '') || ' ' || COALESCE("email", ''))) STORED
);
CREATE INDEX IF NOT EXISTS "users_email" ON "users" ("email");
CREATE INDEX IF NOT EXISTS "users_ts" ON "users" USING GIN ("ts");

CREATE TABLE IF NOT EXISTS "group_user_relations" (
    "groupId" VARCHAR(20) NOT NULL,
    "userId" VARCHAR(20) NOT NULL,
    "role" TEXT NOT NULL,
    PRIMARY KEY ("groupId", "userId")
);
CREATE INDEX IF NOT EXISTS "group_user_relations_role" ON "group_user_relations" ("role");

CREATE TABLE IF NOT EXISTS "credentials" (
    "userId" VARCHAR(20) PRIMARY KEY,
    "password" TEXT,
    "otpSecret" TEXT,
    "credentialId" VARCHAR(1023),
    "publicKey" TEXT
);
CREATE INDEX IF NOT EXISTS "credentials_credentialId" ON "credentials" ("credentialId");

-- +migrate Down
DROP TABLE "users";
DROP TABLE "groups";
DROP TABLE "group_user_relations";
DROP TABLE "credentials";