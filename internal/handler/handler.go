package handler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/c4erries/Sentry/internal/alert"
	"github.com/c4erries/Sentry/internal/anomaly"
	"github.com/c4erries/Sentry/internal/model"
	"github.com/c4erries/Sentry/internal/redis"
)

type Processor struct {
	dupGuard   *anomaly.DuplicateGuard
	registry   *anomaly.DetectorRegistry
	dispatcher *alert.Dispatcher
	// db *sql.DB
}

func NewProcessor(cache redis.RedisClient, dispatcher *alert.Dispatcher) (*Processor, error) {
	dupGuard := anomaly.NewDuplicateGuard(cache, 15*time.Second)
	reg := &anomaly.DetectorRegistry{}
	reg.Registry(anomaly.NewLoginStormDetector(cache, 15*time.Minute, 5))
	reg.Registry(anomaly.NewGeoSwitchingDetector(cache, 5*time.Minute))
	return &Processor{
		dupGuard:   dupGuard,
		registry:   reg,
		dispatcher: dispatcher,
	}, nil
}

func (p *Processor) Process(ctx context.Context, e *model.Event) error {
	isDup, err := p.dupGuard.IsDuplicate(ctx, e)
	if err != nil {
		return fmt.Errorf("dupGuard ISDUPLICATE error: %v", err)
	}

	if isDup {
		return nil
	}
	alerts := p.registry.ProcessAll(ctx, e)
	for _, alert := range alerts {
		select {
		case p.dispatcher.Chan <- alert:
		default:
			log.Panicln("alert channel full, dropped alert:", alert)
		}
	}
	return nil
}
