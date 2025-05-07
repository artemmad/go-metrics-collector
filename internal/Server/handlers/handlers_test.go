package handlers

import (
	"context"
	"github.com/artemmad/go-metrics-collector/internal"
	"github.com/artemmad/go-metrics-collector/internal/Server/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMetricCalc(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantBody   string
		prepare    func(store storage.Storage)
		verify     func(t *testing.T, store storage.Storage)
	}{
		{
			name:       "valid gauge",
			method:     http.MethodPost,
			path:       "/update/gauge/temp/36.6",
			wantStatus: http.StatusOK,
			verify: func(t *testing.T, store storage.Storage) {
				val := store.GetGauges()["temp"]
				assert.Equal(t, 36.6, val)
			},
		},
		{
			name:       "valid counter",
			method:     http.MethodPost,
			path:       "/update/counter/requests/42",
			wantStatus: http.StatusOK,
			verify: func(t *testing.T, store storage.Storage) {
				val := store.GetCounters()["requests"]
				assert.Equal(t, int64(42), val)
			},
		},
		{
			name:       "invalid type",
			method:     http.MethodPost,
			path:       "/update/unknown/type/123",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid value",
			method:     http.MethodPost,
			path:       "/update/gauge/temp/abc",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := storage.NewMemStorage()
			if tt.prepare != nil {
				tt.prepare(store)
			}

			r := chi.NewRouter()
			r.Post("/update/{metricType}/{metricName}/{value}", MetricCalc(store))

			req := httptest.NewRequest(http.MethodPost, tt.path, nil)
			rw := httptest.NewRecorder()
			r.ServeHTTP(rw, req)

			res := rw.Result()
			assert.Equal(t, tt.wantStatus, res.StatusCode)
			if tt.wantBody != "" {
				body, _ := io.ReadAll(res.Body)
				defer res.Body.Close()
				assert.Contains(t, string(body), tt.wantBody)
			}
			if tt.verify != nil {
				tt.verify(t, store)
			}
		})
	}
}

func TestMetricList(t *testing.T) {
	tests := []struct {
		name       string
		prepare    func(store storage.Storage)
		wantStatus int
		wantBody   []string
	}{
		{
			name: "non-empty list",
			prepare: func(store storage.Storage) {
				store.SetGauge("temp", 36.6)
				store.SetCounter("hits", 123)
			},
			wantStatus: http.StatusOK,
			wantBody: []string{
				"GAUGE:",
				"\ttemp: 36.600000",
				"COUNTER:",
				"\thits: 123",
			},
		},
		{
			name:       "empty list",
			wantStatus: http.StatusOK,
			wantBody:   []string{"GAUGE:", "COUNTER:"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := storage.NewMemStorage()
			if tt.prepare != nil {
				tt.prepare(store)
			}

			r := httptest.NewRequest(http.MethodGet, "/", nil)
			rw := httptest.NewRecorder()
			MetricList(store).ServeHTTP(rw, r)

			res := rw.Result()
			assert.Equal(t, tt.wantStatus, res.StatusCode)

			body, _ := io.ReadAll(res.Body)
			defer res.Body.Close()
			bodyStr := string(body)
			for _, expected := range tt.wantBody {
				assert.Contains(t, bodyStr, expected)
			}
		})
	}
}

func TestGetOneMetric(t *testing.T) {
	tests := []struct {
		name           string
		metricType     string
		metricName     string
		prepare        func(store *storage.MemStorage)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:       "existing gauge",
			metricType: internal.GaugeType,
			metricName: "temp",
			prepare: func(s *storage.MemStorage) {
				s.SetGauge("temp", 36.6)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "36.6",
		},
		{
			name:       "existing counter",
			metricType: internal.CounterType,
			metricName: "requests",
			prepare: func(s *storage.MemStorage) {
				s.SetCounter("requests", 42)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "42",
		},
		{
			name:           "non-existent metric",
			metricType:     internal.GaugeType,
			metricName:     "missing",
			prepare:        func(s *storage.MemStorage) {},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "metric not found\n",
		},
		{
			name:           "unknown metric type",
			metricType:     "unknown",
			metricName:     "test",
			prepare:        func(s *storage.MemStorage) {},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "unknown metric type\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := storage.NewMemStorage()
			if tt.prepare != nil {
				tt.prepare(store)
			}

			r := httptest.NewRequest(http.MethodGet, "/value/"+tt.metricType+"/"+tt.metricName, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("metricType", tt.metricType)
			rctx.URLParams.Add("metricName", tt.metricName)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			rw := httptest.NewRecorder()
			GetOneMetric(store).ServeHTTP(rw, r)

			res := rw.Result()
			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			body := rw.Body.String()
			assert.Equal(t, tt.expectedBody, body)
			rw.Result().Body.Close()
		})
	}
}
