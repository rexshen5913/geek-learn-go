package prometheus

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rexshen5913/geek-learn-go/geektime-go /orm/homework3"
	"time"
)

type MiddlewareBuilder struct {
	Name        string
	Subsystem   string
	ConstLabels map[string]string
	Help        string
}

func (m *MiddlewareBuilder) Build() orm.Middleware {
	summaryVec := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:        m.Name,
		Subsystem:   m.Subsystem,
		ConstLabels: m.ConstLabels,
		Help:        m.Help,
	}, []string{"type", "table"})

	return func(next orm.HandleFunc) orm.HandleFunc {
		return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
			startTime := time.Now()
			defer func() {
				endTime := time.Now()
				summaryVec.WithLabelValues(qc.Type, qc.Model.TableName).
					Observe(float64(endTime.Sub(startTime).Milliseconds()))
			}()
			return next(ctx, qc)
		}
	}
}
