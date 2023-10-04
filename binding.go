package metricsotel

import (
	"context"
	"math"
	"strings"
	"sync/atomic"

	logging "github.com/ipfs/go-log/v2"
	metrics "github.com/ipfs/go-metrics-interface"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"

	metricsprometheus "github.com/MichaelMure/go-metrics-otel/prometheus"
	"github.com/MichaelMure/go-metrics-otel/split"
)

var log logging.EventLogger = logging.Logger("metrics-otel")

func Inject() error {
	return metrics.InjectImpl(NewCreator)
}

func InjectOtelAndPrometheus() error {
	return metrics.InjectImpl(split.NewSplit(NewCreator, metricsprometheus.NewCreator))
}

func NewCreator(name, helptext string) metrics.Creator {
	// Note: prometheus and OTEL convention are fairy different. The code below is quite simplistic, and ideally
	// the go-metrics-interface creator API would be reworked to account for it.

	var meterName, instrumentName string
	i := strings.LastIndex(name, ".")
	if i < 0 {
		meterName = "default"
		instrumentName = name
	} else {
		meterName, instrumentName = name[:i], name[i+1:]
	}

	return &creator{
		meterName:      meterName,
		instrumentName: instrumentName,
		helptext:       helptext,
	}
}

var _ metrics.Creator = &creator{}

type creator struct {
	meterName      string
	instrumentName string
	helptext       string
}

func (c *creator) Counter() metrics.Counter {
	res, err := otel.Meter(c.meterName).Float64Counter(c.instrumentName, metric.WithDescription(c.helptext))
	if err != nil {
		log.Warnf("error creating otel counter: %s", err.Error())
		return (&noop{}).Counter()
	}
	return &counterWrapper{res}
}

var _ metrics.Counter = &counterWrapper{}

type counterWrapper struct {
	fc metric.Float64Counter
}

func (c *counterWrapper) Inc() {
	c.fc.Add(context.Background(), 1)
}

func (c *counterWrapper) Add(f float64) {
	c.fc.Add(context.Background(), f)
}

func (c *creator) Gauge() metrics.Gauge {
	m := otel.Meter(c.meterName)

	// current OTEL doesn't have a sync instrument with a set() function, so we have to use the async one ...
	// which complicate quite a lot our code
	res, err := m.Float64ObservableGauge(c.instrumentName, metric.WithDescription(c.helptext))
	if err != nil {
		log.Warnf("error creating otel gauge: %s", err.Error())
		return (&noop{}).Gauge()
	}

	wrapper := &gaugeWrapper{}

	_, err = m.RegisterCallback(func(_ context.Context, o metric.Observer) error {
		o.ObserveFloat64(res, wrapper.Get())
		return nil
	})
	if err != nil {
		log.Warnf("error registering observable on otel gauge: %s", err.Error())
		return (&noop{}).Gauge()
	}

	return &gaugeWrapper{}
}

var _ metrics.Gauge = &gaugeWrapper{}

type gaugeWrapper struct {
	v atomic.Uint64
}

func (g *gaugeWrapper) Set(f float64) {
	g.v.Store(math.Float64bits(f))
}

func (g *gaugeWrapper) Get() float64 {
	return math.Float64frombits(g.v.Load())
}

func (g *gaugeWrapper) Inc() {
	for {
		old := g.v.Load()
		if g.v.CompareAndSwap(old, math.Float64bits(math.Float64frombits(old)+1)) {
			return
		}
	}
}

func (g *gaugeWrapper) Dec() {
	for {
		old := g.v.Load()
		if g.v.CompareAndSwap(old, math.Float64bits(math.Float64frombits(old)-1)) {
			return
		}
	}
}

func (g *gaugeWrapper) Add(f float64) {
	for {
		old := g.v.Load()
		if g.v.CompareAndSwap(old, math.Float64bits(math.Float64frombits(old)+f)) {
			return
		}
	}
}

func (g *gaugeWrapper) Sub(f float64) {
	for {
		old := g.v.Load()
		if g.v.CompareAndSwap(old, math.Float64bits(math.Float64frombits(old)-f)) {
			return
		}
	}
}

func (c *creator) Histogram(buckets []float64) metrics.Histogram {
	res, err := otel.Meter(c.meterName).Float64Histogram(c.instrumentName, metric.WithDescription(c.helptext))
	if err != nil {
		log.Warnf("error creating otel histogram: %s", err.Error())
		return (&noop{}).Histogram(buckets)
	}

	// Note: I understand that the bucketing may be done on the exporter side, but not when declaring that histogram.

	return &histogramWrapper{res}
}

var _ metrics.Histogram = &histogramWrapper{}

type histogramWrapper struct {
	fh metric.Float64Histogram
}

func (h *histogramWrapper) Observe(f float64) {
	h.fh.Record(context.Background(), f)
}

func (c *creator) Summary(opts metrics.SummaryOpts) metrics.Summary {
	res, err := otel.Meter(c.meterName).Float64Histogram(c.instrumentName, metric.WithDescription(c.helptext))
	if err != nil {
		log.Warnf("error creating otel summary: %s", err.Error())
		return (&noop{}).Summary(opts)
	}

	// Note: it seems that OTEL has support for a Summary metric type specifically to migrate old prometheus summaries,
	// but I can't find support for that in the golang libraries. Histogram is the recommended replacement.

	return &summaryWrapper{res}
}

var _ metrics.Summary = &summaryWrapper{}

type summaryWrapper struct {
	fh metric.Float64Histogram
}

func (s *summaryWrapper) Observe(f float64) {
	s.fh.Record(context.Background(), f)
}
