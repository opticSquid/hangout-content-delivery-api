package telemetry

import (
	"context"
	"runtime"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

// startProcessMetrics initializes and starts the process metrics collection
// It uses OpenTelemetry to collect memory usage, CPU percentage, and goroutine count
func StartProcessMetrics() {
	meter := otel.GetMeterProvider().Meter("hangout.content-delivery.metrics")
	heapMemUsage, _ := meter.Float64ObservableGauge("go_heap_memory_usage")
	stackMemUsage, _ := meter.Float64ObservableGauge("go_stack_memory_usage")
	goRoutineCount, _ := meter.Int64ObservableGauge("go_goroutines_count")
	gcCount, _ := meter.Int64ObservableGauge("go_gc_cycle_count")
	gcPause, _ := meter.Float64ObservableGauge("go_gc_all_stop_pause_time_sum")

	_, err := meter.RegisterCallback(
		func(ctx context.Context, o metric.Observer) error {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			o.ObserveFloat64(heapMemUsage, float64(m.Alloc))
			o.ObserveFloat64(stackMemUsage, float64(m.StackInuse))
			o.ObserveInt64(goRoutineCount, int64(runtime.NumGoroutine()))
			o.ObserveFloat64(gcPause, float64(m.PauseTotalNs))
			o.ObserveInt64(gcCount, int64(m.NumGC))

			return nil
		},
		heapMemUsage, stackMemUsage, goRoutineCount, gcCount, gcPause,
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to register metrics callback")
	}
}
