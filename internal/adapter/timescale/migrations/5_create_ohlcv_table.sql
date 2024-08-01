-- +migrate Up
CREATE TABLE IF NOT EXISTS "ohlcv" (
    "id" CHAR(20) NOT NULL,
    "open" DOUBLE PRECISION NOT NULL,
    "high" DOUBLE PRECISION NOT NULL,
    "low" DOUBLE PRECISION NOT NULL,
    "close" DOUBLE PRECISION NOT NULL,
    "volume" INT NOT NULL,
    "time" TIMESTAMPTZ NOT NULL
);

SELECT create_hypertable('ohlcv', by_range('time'));
CREATE INDEX id_time_ohlcv ON "ohlcv" ("id", "time" DESC);

-- +migrate Down

DROP TABLE "ohlcv";