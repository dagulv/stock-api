package service

import (
	"context"
	"log"
	"time"

	"github.com/dagulv/stock-api/internal/core/port"
	"github.com/dagulv/ticker"
)

type Ticker struct {
	Store port.Ticker
	ticks []ticker.Tick
}

func (s *Ticker) SpawnBatcher(ctx context.Context) (err error) {
	t := time.NewTicker(time.Second * 30)

	go func() {
		for {
			select {
			case <-ctx.Done():
				t.Stop()
				return

			case <-t.C:
				if err = s.Store.CopyFrom(ctx, s.ticks); err != nil {
					log.Println(err)
					//Continue
				}

				s.ticks = s.ticks[:0]
			}
		}
	}()

	return
}

func (s *Ticker) ExposeTick(ctx context.Context, tick ticker.Tick) (err error) {
	s.ticks = append(s.ticks, tick)

	return
}
