package internal

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReportCounterMetric(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		wantPath := "/update/counter/mycounter/42"
		if r.URL.Path != wantPath {
			t.Errorf("unexpected path: got %v want %v", r.URL.Path, wantPath)
		}
	}))
	defer ts.Close()

	SetServerAddress(ts.URL)
	ReportCounterMetric("mycounter", 42)
	if !called {
		t.Error("handler was not called")
	}
}

func TestReportGaugeMetric(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if !strings.Contains(r.URL.Path, "/update/gauge/temperature/") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer ts.Close()

	SetServerAddress(ts.URL)
	ReportGaugeMetric("temperature", 36.6)
	if !called {
		t.Error("handler was not called")
	}
}

func TestReportMetric(t *testing.T) {
	called := false

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		expected := "/update/gauge/test/123"
		if r.URL.Path != expected {
			t.Errorf("unexpected URL path: got %s, want %s", r.URL.Path, expected)
		}
	}))
	defer ts.Close()

	SetServerAddress(ts.URL)
	ReportMetric("gauge", "test", 123)
	if !called {
		t.Error("handler was not called")
	}
}

func TestUpdateMetrics(t *testing.T) {
	prev := pollCount
	UpdateMetrics()
	if pollCount != prev+1 {
		t.Errorf("pollCount not incremented: got %d, want %d", pollCount, prev+1)
	}
	if randomValue == 0 {
		t.Error("randomValue not updated")
	}
}

func Test_setServerAddress(t *testing.T) {
	SetServerAddress("http://test:1234")
	if ServerAddress != "http://test:1234" {
		t.Errorf("ServerAddress not updated: got %s", ServerAddress)
	}
}

func TestReportMetrics(t *testing.T) {
	var counter int
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		counter++
	}))
	defer ts.Close()

	SetServerAddress(ts.URL)
	ReportMetrics()

	if counter == 0 {
		t.Error("no metrics were reported")
	}
}
