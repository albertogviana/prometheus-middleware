package prommiddleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	defaultBuckets = []float64{0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5}
	defaultLabels  = []string{"status", "method", "path"}
)

const (
	requestName = "http_requests_total"
	latencyName = "http_request_duration_seconds"
)

// Opts specifies options how to create new PrometheusMiddleware.
type Opts struct {
	// Buckets specifies an custom buckets to be used in request histograpm.
	Buckets []float64
	// Labels specifies the label names that will be used
	Labels []string
}

// PromMiddleware specifies the metrics that is going to be generated
type PromMiddleware struct {
	request *prometheus.CounterVec
	latency *prometheus.HistogramVec
}

// New creates a new PrometheusMiddleware instance
func New(opts Opts) (*PromMiddleware, error) {
	var prometheusMiddleware PromMiddleware

	labels := opts.Labels
	if len(labels) == 0 {
		labels = defaultLabels
	}

	prometheusMiddleware.request = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: requestName,
			Help: "How many HTTP requests processed, partitioned by status code, method and HTTP path.",
		},
		labels,
	)

	if err := prometheus.Register(prometheusMiddleware.request); err != nil {
		return nil, fmt.Errorf("could not register request metric %w", err)
	}

	buckets := opts.Buckets
	if len(buckets) == 0 {
		buckets = defaultBuckets
	}

	prometheusMiddleware.latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    latencyName,
		Help:    "How long it took to process the request, partitioned by status code, method and HTTP path.",
		Buckets: buckets,
	},
		labels,
	)

	if err := prometheus.Register(prometheusMiddleware.latency); err != nil {
		return nil, fmt.Errorf("could not register latency metric %w", err)
	}

	return &prometheusMiddleware, nil
}

// InstrumentHandlerDuration is a middleware that wraps the http.Handler and it record
// how long the handler took to run, which path was called, and the status code.
// This method is going to be used with gorilla/mux.
func (p *PromMiddleware) InstrumentHandlerDuration(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()

		delegate := &responseWriterDelegator{ResponseWriter: w}
		rw := delegate

		next.ServeHTTP(rw, r) // call original

		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		code := sanitizeCode(delegate.status)
		method := sanitizeMethod(r.Method)

		go p.request.WithLabelValues(
			code,
			method,
			path,
		).Inc()

		go p.latency.WithLabelValues(
			code,
			method,
			path,
		).Observe(float64(time.Since(begin)) / float64(time.Second))
	})
}

type responseWriterDelegator struct {
	http.ResponseWriter
	status      int
	written     int64
	wroteHeader bool
}

func (r *responseWriterDelegator) WriteHeader(code int) {
	r.status = code
	r.wroteHeader = true
	r.ResponseWriter.WriteHeader(code)
}

func (r *responseWriterDelegator) Write(b []byte) (int, error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	n, err := r.ResponseWriter.Write(b)
	r.written += int64(n)
	return n, err
}

func sanitizeMethod(m string) string {
	return strings.ToLower(m)
}

func sanitizeCode(s int) string {
	return strconv.Itoa(s)
}
