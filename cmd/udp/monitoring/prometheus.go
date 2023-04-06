package monitoring

import (
	dbMonitoring "github.com/jacobbrewer1/reverse-proxy/pkg/dataacess/monitoring"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	TotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "udpproxy_http_total_requests",
			Help: "Total number of events received",
		},
		[]string{"udpproxy"},
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "udpproxy_http_request_duration",
			Help: "Duration of proxy requests",
		},
		[]string{"udpproxy"},
	)

	TotalSystemErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "udpproxy_system_errors",
			Help: "Number of system errors",
		},
		[]string{"udpproxy"},
	)

	RedisLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "udpproxy_redis_latency",
			Help: "Duration of redis requests",
		},
		[]string{"collection"},
	)
)

func Register() error {
	if err := prometheus.Register(TotalRequests); err != nil {
		return err
	}

	if err := prometheus.Register(RequestDuration); err != nil {
		return err
	}

	if err := prometheus.Register(TotalSystemErrors); err != nil {
		return err
	}

	dbMonitoring.RedisLatency = RedisLatency
	if err := prometheus.Register(RedisLatency); err != nil {
		return err
	}

	return nil
}

func Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	}
}
