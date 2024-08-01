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
	ohlcv []ticker.Ohlcv
}

func (s *Ticker) Spawn(ctx context.Context) (err error) {
	var avanzaIds []int

	if avanzaIds, err = s.Store.GetAvanzaIds(ctx, avanzaIds); err != nil {
		return
	}

	t := ticker.New[ticker.Tick](s, avanzaIds...)

	t.FetchHistoric(ctx)

	// if err = t.StartWebsocket(ctx); err != nil {
	// 	return
	// }

	timeTicker := time.NewTicker(time.Second * 60)

	go func() {
		for {
			select {
			case <-ctx.Done():
				timeTicker.Stop()
				return

			case <-timeTicker.C:
				if err = s.Store.InsertHistoricOhlcv(ctx, s.ohlcv); err != nil {
					log.Println(err)
					//Continue
				}

				s.ohlcv = s.ohlcv[:0]
			}
		}
	}()

	return
}

func (s *Ticker) ExposeTick(ctx context.Context, tick ticker.Tick) (err error) {
	s.ticks = append(s.ticks, tick)

	return
}

func (s *Ticker) Push(m ticker.Method) {
	log.Println(m)
	log.Println(m.(ticker.Ohlcv))
	s.ohlcv = append(s.ohlcv, (m).(ticker.Ohlcv))
}
