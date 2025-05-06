package handlers

import (
	"github.com/artemmad/go-metrics-collector/internal/Server/storage"
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

			r := httptest.NewRequest(tt.method, tt.path, nil)
			rw := httptest.NewRecorder()
			h := MetricCalc(store)
			h.ServeHTTP(rw, r)

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
