package prommiddleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

func Test_InstrumentGorillaMux(t *testing.T) {
	recorder := httptest.NewRecorder()

	reg := prometheus.NewRegistry()

	middleware, err := New(Opts{Registry: reg})
	if err != nil {
		t.Errorf("error initializing middleware %s", err)
	}

	r := mux.NewRouter()
	r.Handle("/metrics", middleware.Handler())
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	r.Use(middleware.InstrumentHandlerDuration)

	ts := httptest.NewServer(r)
	defer ts.Close()

	req1, err := http.NewRequest("GET", ts.URL+"/?version=0.1", nil)
	if err != nil {
		t.Error(err)
	}
	req2, err := http.NewRequest("GET", ts.URL+"/metrics", nil)
	if err != nil {
		t.Error(err)
	}

	r.ServeHTTP(recorder, req1)
	r.ServeHTTP(recorder, req2)
	body := recorder.Body.String()
	if !strings.Contains(body, requestName) {
		t.Errorf("body does not contain request total entry '%s'", requestName)
	}
	if !strings.Contains(body, latencyName) {
		t.Errorf("body does not contain request duration entry '%s'", requestName)
	}

    if !strings.Contains(body, `http_request_duration_seconds_count{method="get",path="/",status="200",version="0.1"}`) {
        t.Errorf("body does not contain expected version value '%s'", `http_request_duration_seconds_count{method="get",path="/",status="200",version="0.1"}`)
    }
}
