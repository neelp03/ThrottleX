package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "net/http"
)

// TotalRequests counts the total number of requests handled by the rate limiter.
//
// This metric tracks every request processed by the rate limiter, whether it is allowed
// or denied. It is a Prometheus counter, meaning it only increments and is not reset
// over time.
//
// Prometheus Metric:
//   - Name: `throttlex_total_requests`
//   - Type: Counter
//   - Help: "Total number of requests handled by the rate limiter."
var TotalRequests = prometheus.NewCounter(
    prometheus.CounterOpts{
        Name: "throttlex_total_requests",
        Help: "Total number of requests handled by the rate limiter.",
    },
)

// DeniedRequests counts the total number of requests denied due to exceeding the rate limit.
//
// This metric tracks the number of requests that were rejected by the rate limiter because
// the request count exceeded the configured rate limit. It is a Prometheus counter, meaning
// it only increments and is not reset over time.
//
// Prometheus Metric:
//   - Name: `throttlex_denied_requests`
//   - Type: Counter
//   - Help: "Total number of requests denied due to rate limiting."
var DeniedRequests = prometheus.NewCounter(
    prometheus.CounterOpts{
        Name: "throttlex_denied_requests",
        Help: "Total number of requests denied due to rate limiting.",
    },
)

// RegisterMetrics registers the Prometheus metrics for tracking requests.
//
// This function registers the `TotalRequests` and `DeniedRequests` counters with Prometheus.
// It must be called during the application startup to ensure that these metrics are
// available and can be scraped by Prometheus.
//
// Example usage:
//   metrics.RegisterMetrics()
func RegisterMetrics() {
    prometheus.MustRegister(TotalRequests)
    prometheus.MustRegister(DeniedRequests)
}

// ExposeMetrics starts an HTTP server to expose Prometheus metrics at the `/metrics` endpoint.
//
// This function launches an HTTP server on port `2112` to expose the registered Prometheus
// metrics. Prometheus scrapers can fetch metrics from this endpoint to track the
// performance and behavior of the rate limiter.
//
// Endpoint:
//   - `/metrics`: Standard Prometheus metrics endpoint.
//
// Example usage:
//   metrics.ExposeMetrics()
//
// The `/metrics` endpoint can be scraped by Prometheus to collect metrics.
//
// Example Prometheus configuration for scraping the `/metrics` endpoint:
//
// scrape_configs:
//   - job_name: "throttlex_rate_limiter"
//     static_configs:
//       - targets: ['localhost:2112']
func ExposeMetrics() {
    http.Handle("/metrics", promhttp.Handler())
    go http.ListenAndServe(":2112", nil) // Start HTTP server on port 2112
}
