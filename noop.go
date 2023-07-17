package metricsotel

import "github.com/ipfs/go-metrics-interface"

// Note: this was copied from go-metrics-interface as it's not public there.

// Also implements the Counter interface
type noop struct{}

func (g *noop) Set(v float64) {
	// Noop
}

func (g *noop) Inc() {
	// Noop
}

func (g *noop) Dec() {
	// Noop
}

func (g *noop) Add(v float64) {
	// Noop
}

func (g *noop) Sub(v float64) {
	// Noop
}

func (g *noop) Observe(v float64) {
	// Noop
}

// Creator functions

func (g *noop) Counter() metrics.Counter {
	return g
}

func (g *noop) Gauge() metrics.Gauge {
	return g
}

func (g *noop) Histogram(buckets []float64) metrics.Histogram {
	return g
}

func (g *noop) Summary(opts metrics.SummaryOpts) metrics.Summary {
	return g
}
