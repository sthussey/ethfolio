package internal

import (
	"net/http"
	"log"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)


type MetricDefinitions struct {
	txns_processed				prom.Counter
	blocks_processed				prom.Counter
	last_block_completed		prom.Gauge
}

func GetMetricDefinitions() MetricDefinitions {
	return MetricDefinitions{txns_processed: prom.NewCounter(prom.CounterOpts{Name: "txns_processed", Help: "Number of Ethereum transactions checked for lot qualification"}),
		blocks_processed: prom.NewCounter(prom.CounterOpts{Name: "blocks_processed", Help: "Number of Ethereum blocks where the full set of txns have been checked for lot qualification"}),
		last_block_completed: prom.NewGauge(prom.GaugeOpts{Name: "last_block_completed", Help: "The number of the last block fully processed"})}
}

func ServeMetrics(registry *prom.Registry) {
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry, EnableOpenMetrics: true})
	http.Handle("/metrics", h)
	log.Fatal(http.ListenAndServe("127.0.0.1:5000", nil))
}

func InitializeMetrics() MetricDefinitions {
	metrics := GetMetricDefinitions()
	r := prom.NewRegistry()
	r.MustRegister(metrics.txns_processed)
	r.MustRegister(metrics.blocks_processed)
	r.MustRegister(metrics.last_block_completed)
	go ServeMetrics(r)
	return metrics
}
