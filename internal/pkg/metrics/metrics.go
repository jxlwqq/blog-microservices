package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"strconv"
)

type Metrics interface {
	IncHits(status int, method, path string)
	ObserveResponseTime(status int, method, path string, observeTime float64)
}

type metrics struct {
	HitsTotal prometheus.Counter
	Hits      *prometheus.CounterVec
	Times     *prometheus.HistogramVec
}

func (m *metrics) IncHits(status int, method, path string) {
	m.HitsTotal.Inc()
	m.Hits.WithLabelValues(strconv.Itoa(status), method, path).Inc()
}

func (m *metrics) ObserveResponseTime(status int, method, path string, observeTime float64) {
	m.Times.WithLabelValues(strconv.Itoa(status), method, path).Observe(observeTime)
}

func New(name string) (Metrics, error) {
	var m metrics
	m.HitsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: name + "_server_hits_total",
		Help: "Total number of requests",
	})

	err := prometheus.Register(m.HitsTotal)
	if err != nil {
		return nil, err
	}

	m.Hits = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: name + "_server_hits",
		Help: "Number of requests per status code",
	}, []string{"status", "method", "path"})

	err = prometheus.Register(m.Hits)
	if err != nil {
		return nil, err
	}

	m.Times = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: name + "_server_times",
		Help: "Response time per status code",
	}, []string{"status", "method", "path"})

	err = prometheus.Register(m.Times)
	if err != nil {
		return nil, err
	}

	err = prometheus.Register(collectors.NewBuildInfoCollector())
	if err != nil {
		return nil, err
	}

	return &m, nil
}
