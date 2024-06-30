package port

import (
	"context"

	"github.com/dagulv/ticker"
)

type Ticker interface {
	ticker.Ticker
	Insert(context.Context, ticker.Tick) error
}
