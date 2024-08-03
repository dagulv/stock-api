package port

import (
	"context"

	"github.com/dagulv/ticker"
	"github.com/rs/xid"
)

type Ticker interface {
	// Insert(context.Context, ticker.Tick) error
	CopyFrom(context.Context, []ticker.Tick) error
	InsertHistoricOhlcv(context.Context, []ticker.Ohlcv) error
	GetAvanzaIds(context.Context, []xid.ID, []int) ([]xid.ID, []int, error)
}
