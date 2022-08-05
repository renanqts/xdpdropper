package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var namespace = "xdpdropper"

type metrics struct {
	name      string
	subsystem string
	collector prometheus.Collector
}

var m []metrics

// Init initializes singleton metrics
func Init() {
	if m == nil {
		m = []metrics{}
	}
}

func Register() error {
	for _, m := range m {
		if err := prometheus.Register(m.collector); err != nil {
			return err
		}
	}
	return nil
}

func NewCounter(name, subsystem, help string) prometheus.Counter {
	c := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      name,
		Help:      help,
	})
	m = append(m, metrics{subsystem, name, c})
	return c
}

func NewGauge(name, subsystem, help string) prometheus.Gauge {
	g := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      name,
		Help:      help,
	})
	m = append(m, metrics{subsystem, name, g})
	return g
}

func NewGaugeVec(name, subsystem, help string, labelNames []string) *prometheus.GaugeVec {
	g := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      name,
		Help:      help,
	}, labelNames)
	m = append(m, metrics{subsystem, name, g})
	return g
}
