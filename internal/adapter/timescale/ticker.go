package timescale

import (
	"context"

	"github.com/dagulv/stock-api/internal/core/port"
	"github.com/dagulv/ticker"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"
)

type tickerStore struct {
	db *pgxpool.Pool
}

func NewTicker(db *pgxpool.Pool) port.Ticker {
	s := tickerStore{
		db: db,
	}

	return s
}

// func (p timescale) ExposeTick(ctx context.Context, tick ticker.Tick) (err error) {
// 	_, err = p.db.Exec(ctx, `
// 	INSERT INTO "stocks_data" (
// 		"time",
// 		"symbol",
// 		"price",
// 		"dayVolume"
// 	)
// 	VALUES (
// 		$1, $2, $3, $4
// 	)
// 	`, tick.Time, tick.Symbol, tick.Price, tick.DayVolume)
// 	log.Printf("error %s", err)
// 	log.Println(tick)
// 	return err
// }

func (s tickerStore) CopyFrom(ctx context.Context, ticks []ticker.Tick) (err error) {
	_, err = s.db.CopyFrom(ctx, pgx.Identifier{"stocks_data"}, []string{"time", "symbol", "price", "dayVolume"}, pgx.CopyFromSlice(len(ticks), func(i int) ([]any, error) {
		return []any{ticks[i].Time, ticks[i].Symbol, ticks[i].Price, ticks[i].DayVolume}, nil
	}))

	return
}

func (s tickerStore) InsertHistoricOhlcv(ctx context.Context, ohlcv []ticker.Ohlcv) (err error) {
	_, err = s.db.CopyFrom(ctx, pgx.Identifier{"ohlcv"}, []string{"id", "open", "high", "low", "close", "volume", "time"}, pgx.CopyFromSlice(len(ohlcv), func(i int) ([]any, error) {
		return []any{
			ohlcv[i].Id,
			ohlcv[i].Open,
			ohlcv[i].High,
			ohlcv[i].Low,
			ohlcv[i].Close,
			ohlcv[i].Volume,
			ohlcv[i].Time}, nil
	}))

	return
}

func (s tickerStore) GetAvanzaIds(ctx context.Context, ids []xid.ID, avanzaIds []int) (_ []xid.ID, _ []int, err error) {
	rows, err := s.db.Query(ctx, `
		SELECT
			"id",
			"avanza_id"
		FROM "company"
	`)

	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var id xid.ID
		var avanzaId pgtype.Int4

		if err = rows.Scan(&id, &avanzaId); err != nil {
			return
		}

		ids = append(ids, id)
		avanzaIds = append(avanzaIds, int(avanzaId.Int32))
	}

	return ids, avanzaIds, nil
}
