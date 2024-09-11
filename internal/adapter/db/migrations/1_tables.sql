-- +migrate Up
CREATE TABLE IF NOT EXISTS "company" (
    "symbol" TEXT PRIMARY KEY,
    "name" TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS "stocks_data" (
  "time" TIMESTAMPTZ NOT NULL,
  "symbol" TEXT NOT NULL,
  "price" DOUBLE PRECISION NULL,
  "dayVolume" BIGINT NULL
);

-- +migrate Down
DROP TABLE "company";
DROP TABLE "stocks_data";