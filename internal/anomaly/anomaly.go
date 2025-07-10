package anomaly

import (
	"context"
	"log"
	"time"

	"github.com/c4erries/Sentry/internal/model"
)

const (
	redisTTLBuffer = 30 * time.Second
)

type Detector interface {
	ID() string
	Process(ctx context.Context, event *model.Event) (*model.Alert, error)
}

type DetectorRegistry struct {
	detectors []Detector
}

func (r *DetectorRegistry) Registry(d Detector) {
	r.detectors = append(r.detectors, d)
}

func (r *DetectorRegistry) ProcessAll(ctx context.Context, e *model.Event) []*model.Alert {
	var result []*model.Alert
	for _, d := range r.detectors {
		alert, err := d.Process(ctx, e)
		if err != nil {
			log.Printf("[Detector-%s] error: %v", d.ID(), err)
			continue
		}
		if alert == nil {
			continue
		}
		result = append(result, alert)
	}
	//log.Printf("Return result len:%v", len(result))
	return result
}
