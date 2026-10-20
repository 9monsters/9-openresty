// Package metrics exposes Prometheus metrics for the hub.
package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "openresty_admin_http_requests_total",
		Help: "API requests by method, path, and status.",
	}, []string{"method", "path", "status"})

	httpRequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "openresty_admin_http_request_duration_seconds",
		Help:    "API request latency histogram.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path"})

	agentsOnline = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "openresty_admin_agents_online",
		Help: "Number of connected agents.",
	})

	deployTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "openresty_admin_deploy_total",
		Help: "Config releases by status.",
	}, []string{"status"})

	deployDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "openresty_admin_deploy_duration_seconds",
		Help:    "End-to-end deploy latency per release.",
		Buckets: []float64{0.1, 0.5, 1, 2, 5, 10, 30, 60},
	})

	agentRPCDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "openresty_admin_agent_rpc_duration_seconds",
		Help:    "gRPC call latency per agent and method.",
		Buckets: prometheus.DefBuckets,
	}, []string{"agent", "method"})

	agentRPCErrorsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "openresty_admin_agent_rpc_errors_total",
		Help: "gRPC errors per agent and method.",
	}, []string{"agent", "method"})
)

func init() {
	prometheus.MustRegister(
		httpRequestsTotal,
		httpRequestDuration,
		agentsOnline,
		deployTotal,
		deployDuration,
		agentRPCDuration,
		agentRPCErrorsTotal,
	)
}

// Handler returns an http.Handler serving Prometheus metrics.
var Handler = promhttp.Handler()

// RecordHTTPRequest increments the request counter and observes latency.
func RecordHTTPRequest(method, path, status string, duration time.Duration) {
	httpRequestsTotal.WithLabelValues(method, path, status).Inc()
	httpRequestDuration.WithLabelValues(method, path).Observe(duration.Seconds())
}

// SetAgentsOnline sets the gauge of connected agents.
func SetAgentsOnline(n int) {
	agentsOnline.Set(float64(n))
}

// RecordDeploy records a deploy attempt.
func RecordDeploy(status string, duration time.Duration) {
	deployTotal.WithLabelValues(status).Inc()
	deployDuration.Observe(duration.Seconds())
}

// RecordAgentRPC records an agent RPC result.
func RecordAgentRPC(agent, method string, duration time.Duration, isError bool) {
	agentRPCDuration.WithLabelValues(agent, method).Observe(duration.Seconds())
	if isError {
		agentRPCErrorsTotal.WithLabelValues(agent, method).Inc()
	}
}