CREATE TABLE IF NOT EXISTS "users" (
    "id" CHAR(20) NOT NULL PRIMARY KEY,
    "active" BOOLEAN NOT NULL,
    "name" VARCHAR(128) NOT NULL,
    "email" VARCHAR(255) NOT NULL,
    "password" TEXT,
    "timeCreated" TIMESTAMP WITH TIME ZONE NOT NULL,
    "timeUpdated" TIMESTAMP WITH TIME ZONE NOT NULL
);
CREATE INDEX IF NOT EXISTS "users_active" ON "users" ("active");
CREATE INDEX IF NOT EXISTS "users_email" ON "users" ("email");