package memory

import (
	"context"
	"reflect"
	"sort"
	"sync"
	"testing"

	"github.com/ShvetsovYura/metrics-collector/internal/storage/metric"
	"github.com/stretchr/testify/assert"
)

func TestMemStorage_ToList(t *testing.T) {
	type fields struct {
		mx            sync.Mutex
		gaugeMetrics  map[string]metric.Gauge
		counterMetric map[string]metric.Counter
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{{
		name: "Get list of memory metrics",
		fields: fields{
			mx: sync.Mutex{},
			gaugeMetrics: map[string]metric.Gauge{
				"gMetric1": metric.Gauge(0.1935),
				"gMetric2": metric.Gauge(-123),
				"gMetirc3": metric.Gauge(3.45601),
			},
			counterMetric: map[string]metric.Counter{
				"counter": metric.Counter(999),
			},
		},
		args: args{
			ctx: context.Background(),
		},
		want:    []string{"0.1935", "-123", "3.45601", "999"},
		wantErr: false,
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				mx:            tt.fields.mx,
				gaugeMetrics:  tt.fields.gaugeMetrics,
				counterMetric: tt.fields.counterMetric,
			}
			got, err := m.ToList(tt.args.ctx)
			sort.Strings(got)
			sort.Strings(tt.want)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMemStorage_SetCounter_(t *testing.T) {
	type fields struct {
		mx            sync.Mutex
		gaugeMetrics  map[string]metric.Gauge
		counterMetric map[string]metric.Counter
	}
	type args struct {
		ctx  context.Context
		name string
		val  int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "set simple counter value",
			fields: fields{
				mx:            sync.Mutex{},
				gaugeMetrics:  make(map[string]metric.Gauge),
				counterMetric: make(map[string]metric.Counter),
			},
			args: args{
				ctx:  context.Background(),
				name: "counter1",
				val:  100500,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				mx:            tt.fields.mx,
				gaugeMetrics:  tt.fields.gaugeMetrics,
				counterMetric: tt.fields.counterMetric,
			}
			if err := m.SetCounter(tt.args.ctx, tt.args.name, tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("MemStorage.SetCounter() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, int64(m.counterMetric["counter1"]), tt.args.val)
		})
	}
}

func TestMemStorage_SetGauge(t *testing.T) {
	type fields struct {
		mx            sync.Mutex
		gaugeMetrics  map[string]metric.Gauge
		counterMetric map[string]metric.Counter
	}
	type args struct {
		ctx  context.Context
		name string
		val  float64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success set gauge metric",
			fields: fields{
				mx:            sync.Mutex{},
				gaugeMetrics:  make(map[string]metric.Gauge),
				counterMetric: make(map[string]metric.Counter),
			},
			args:    args{ctx: context.Background(), name: "gauge1", val: 123.456},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				mx:            tt.fields.mx,
				gaugeMetrics:  tt.fields.gaugeMetrics,
				counterMetric: tt.fields.counterMetric,
			}
			if err := m.SetGauge(tt.args.ctx, tt.args.name, tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("MemStorage.SetGauge() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, float64(m.gaugeMetrics["gauge1"]), tt.args.val)
		})
	}
}

func TestMemStorage_GetGauge(t *testing.T) {
	type fields struct {
		mx            sync.Mutex
		gaugeMetrics  map[string]metric.Gauge
		counterMetric map[string]metric.Counter
	}
	initFields := fields{
		mx:            sync.Mutex{},
		gaugeMetrics:  make(map[string]metric.Gauge),
		counterMetric: make(map[string]metric.Counter),
	}
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    metric.Gauge
		wantErr bool
	}{
		{name: "success get gauge",
			fields:  initFields,
			args:    args{ctx: context.TODO(), name: "counter1"},
			want:    metric.Gauge(123.456),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				mx:            tt.fields.mx,
				gaugeMetrics:  tt.fields.gaugeMetrics,
				counterMetric: tt.fields.counterMetric,
			}
			m.gaugeMetrics[tt.args.name] = tt.want
			got, err := m.GetGauge(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemStorage.GetGauge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemStorage.GetGauge() = %v, want %v", got, tt.want)
			}

		})
	}
}