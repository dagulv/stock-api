-- +migrate Up
SELECT create_hypertable('stocks_data', by_range('time'));
CREATE INDEX IF NOT EXISTS "stocks_data_symbolTime" ON "stocks_data" ("symbol", "time" DESC);
-- +migrate Down