# prometheus-middleware [![Build Status](https://travis-ci.org/albertogviana/prometheus-middleware.svg?branch=master)](https://travis-ci.org/albertogviana/prometheus-middleware) [![Go Report Card](https://goreportcard.com/badge/github.com/albertogviana/prometheus-middleware)](https://goreportcard.com/report/github.com/albertogviana/prometheus-middleware)

[Prometheus](http://prometheus.io) middleware that supports [negroni](https://github.com/urfave/negroni) and [gorilla/mux](https://github.com/gorilla/mux).

## Installation

```bash
go get -u github.com/albertogviana/prometheus-middleware
```

## What you will get

You will get:

- the HTTP request duration in seconds (`http_request_duration_seconds`) histogram partitioned by status code, method and HTTP path.

```
...
# HELP http_request_duration_seconds How long it took to process the request, partitioned by status code, method and HTTP path.
# TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{code="200",method="get",path="/",le="0.3"} 2
http_request_duration_seconds_bucket{code="200",method="get",path="/",le="1"} 2
http_request_duration_seconds_bucket{code="200",method="get",path="/",le="2.5"} 2
http_request_duration_seconds_bucket{code="200",method="get",path="/",le="5"} 2
http_request_duration_seconds_bucket{code="200",method="get",path="/",le="+Inf"} 2
http_request_duration_seconds_sum{code="200",method="get",path="/"} 3.5256e-05
http_request_duration_seconds_count{code="200",method="get",path="/"} 2
http_request_duration_seconds_bucket{code="200",method="get",path="/metrics",le="0.3"} 2
http_request_duration_seconds_bucket{code="200",method="get",path="/metrics",le="1"} 2
http_request_duration_seconds_bucket{code="200",method="get",path="/metrics",le="2.5"} 2
http_request_duration_seconds_bucket{code="200",method="get",path="/metrics",le="5"} 2
http_request_duration_seconds_bucket{code="200",method="get",path="/metrics",le="+Inf"} 2
http_request_duration_seconds_sum{code="200",method="get",path="/metrics"} 0.001261767
http_request_duration_seconds_count{code="200",method="get",path="/metrics"} 2
...
```

- HTTP request total (`http_requests_total`) partitioned by status code, method, and HTTP path.

```
...
# HELP http_requests_total How many HTTP requests processed, partitioned by status code, method and HTTP path.
# TYPE http_requests_total counter
http_requests_total{code="200",method="get",path="/"} 2
http_requests_total{code="200",method="get",path="/metrics"} 2
...
```

## How to use it

- Negroni

```go
middleware := NewPrometheusMiddleware(Opts{})

n := negroni.New(middleware)
r := http.NewServeMux()
r.Handle("/metrics", promhttp.Handler())
r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, "ok")
})
n.UseHandler(r)
```

- Gorilla/Mux

```go
middleware := NewPrometheusMiddleware(Opts{})

r := mux.NewRouter()
r.Handle("/metrics", promhttp.Handler())
r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, "ok")
})

r.Use(middleware.InstrumentHandlerDuration)
```
