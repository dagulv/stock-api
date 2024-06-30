package timescale

import (
	"context"
	"log"

	"github.com/dagulv/stock-api/internal/core/port"
	"github.com/dagulv/ticker"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewTicker(db *pgxpool.Pool) port.Ticker {
	return timescale{
		db: db,
	}
}

func (p timescale) ExposeTick(ctx context.Context, tick ticker.Tick) (err error) {
	_, err = p.db.Exec(ctx, `
	INSERT INTO "stocks_data" (
		"time",
		"symbol",
		"price",
		"dayVolume"
	)
	VALUES (
		$1, $2, $3, $4
	)
	`, tick.Time, tick.Symbol, tick.Price, tick.DayVolume)
	log.Printf("error %s", err)
	log.Println(tick)
	return err
}
