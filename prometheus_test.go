package prometheusmiddleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/negroni"
)

func Test_InstrumentNegroni(t *testing.T) {
	recorder := httptest.NewRecorder()

	middleware := NewPrometheusMiddleware(Opts{})

	n := negroni.New(middleware)
	n.Use(middleware)
	r := http.NewServeMux()
	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})
	n.UseHandler(r)

	req1, err := http.NewRequest("GET", "http://localhost:3000/", nil)
	if err != nil {
		t.Error(err)
	}
	req2, err := http.NewRequest("GET", "http://localhost:3000/metrics", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(recorder, req1)
	n.ServeHTTP(recorder, req2)
	body := recorder.Body.String()
	if !strings.Contains(body, requestName) {
		t.Errorf("body does not contain request total entry '%s'", requestName)
	}
	if !strings.Contains(body, latencyName) {
		t.Errorf("body does not contain request duration entry '%s'", requestName)
	}
}
