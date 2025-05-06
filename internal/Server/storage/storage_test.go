package storage

import (
	"reflect"
	"testing"
)

func TestMemStorage_GetCounters(t *testing.T) {
	tests := []struct {
		name   string
		preset map[string]int64
		want   map[string]int64
	}{
		{
			name:   "one counter",
			preset: map[string]int64{"hits": 10},
			want:   map[string]int64{"hits": 10},
		},
		{
			name:   "empty map",
			preset: map[string]int64{},
			want:   map[string]int64{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				counters: tt.preset,
				gauges:   map[string]float64{},
			}
			got := s.GetCounters()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCounters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_GetGauges(t *testing.T) {
	tests := []struct {
		name   string
		preset map[string]float64
		want   map[string]float64
	}{
		{
			name:   "one gauge",
			preset: map[string]float64{"temp": 36.6},
			want:   map[string]float64{"temp": 36.6},
		},
		{
			name:   "empty map",
			preset: map[string]float64{},
			want:   map[string]float64{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				gauges:   tt.preset,
				counters: map[string]int64{},
			}
			got := s.GetGauges()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetGauges() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_SetCounter(t *testing.T) {
	s := NewMemStorage()
	s.SetCounter("requests", 42)

	if got := s.GetCounters()["requests"]; got != 42 {
		t.Errorf("SetCounter() = %v, want %v", got, 42)
	}
}

func TestMemStorage_SetGauge(t *testing.T) {
	s := NewMemStorage()
	s.SetGauge("load", 0.75)

	if got := s.GetGauges()["load"]; got != 0.75 {
		t.Errorf("SetGauge() = %v, want %v", got, 0.75)
	}
}

func TestNewMemStorage(t *testing.T) {
	s := NewMemStorage()
	if s == nil {
		t.Fatal("NewMemStorage() returned nil")
	}
	if len(s.gauges) != 0 || len(s.counters) != 0 {
		t.Errorf("expected empty maps, got gauges: %v, counters: %v", s.gauges, s.counters)
	}
}
