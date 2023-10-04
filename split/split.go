package split

import "github.com/ipfs/go-metrics-interface"

func NewSplit(constructors ...metrics.InternalNew) metrics.InternalNew {
	return func(name, helptext string) metrics.Creator {
		res := make(splitCreator, len(constructors))
		for i, constructor := range constructors {
			res[i] = constructor(name, helptext)
		}
		return res
	}
}

type splitCreator []metrics.Creator

func (sc splitCreator) Counter() metrics.Counter {
	res := make(splitCounter, len(sc))
	for i, creator := range sc {
		res[i] = creator.Counter()
	}
	return res
}

func (sc splitCreator) Gauge() metrics.Gauge {
	res := make(splitGauge, len(sc))
	for i, creator := range sc {
		res[i] = creator.Gauge()
	}
	return res
}

func (sc splitCreator) Histogram(buckets []float64) metrics.Histogram {
	res := make(splitHistogram, len(sc))
	for i, creator := range sc {
		res[i] = creator.Histogram(buckets)
	}
	return res
}

func (sc splitCreator) Summary(opts metrics.SummaryOpts) metrics.Summary {
	res := make(splitSummary, len(sc))
	for i, creator := range sc {
		res[i] = creator.Summary(opts)
	}
	return res
}

var _ metrics.Counter = splitCounter{}

type splitCounter []metrics.Counter

func (sc splitCounter) Inc() {
	for _, counter := range sc {
		counter.Inc()
	}
}

func (sc splitCounter) Add(f float64) {
	for _, counter := range sc {
		counter.Add(f)
	}
}

var _ metrics.Gauge = splitGauge{}

type splitGauge []metrics.Gauge

func (sg splitGauge) Set(f float64) {
	for _, gauge := range sg {
		gauge.Set(f)
	}
}

func (sg splitGauge) Inc() {
	for _, gauge := range sg {
		gauge.Inc()
	}
}

func (sg splitGauge) Dec() {
	for _, gauge := range sg {
		gauge.Dec()
	}
}

func (sg splitGauge) Add(f float64) {
	for _, gauge := range sg {
		gauge.Add(f)
	}
}

func (sg splitGauge) Sub(f float64) {
	for _, gauge := range sg {
		gauge.Sub(f)
	}
}

var _ metrics.Histogram = splitHistogram{}

type splitHistogram []metrics.Histogram

func (sh splitHistogram) Observe(f float64) {
	for _, histogram := range sh {
		histogram.Observe(f)
	}
}

var _ metrics.Summary = splitSummary{}

type splitSummary []metrics.Summary

func (ss splitSummary) Observe(f float64) {
	for _, summary := range ss {
		summary.Observe(f)
	}
}
