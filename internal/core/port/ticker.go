package port

import (
	"context"

	"github.com/dagulv/ticker"
)

type Ticker interface {
	// Insert(context.Context, ticker.Tick) error
	CopyFrom(context.Context, []ticker.Tick) error
}
