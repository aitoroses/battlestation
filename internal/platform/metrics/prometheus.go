package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Request metrics
	RequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "battlestation_request_duration_seconds",
		Help:    "Time taken to process attack requests",
		Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1},
	}, []string{"protocol"})

	RequestTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "battlestation_requests_total",
		Help: "Total number of attack requests",
	}, []string{"protocol", "status"})

	// Target selection metrics
	TargetSelectionDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "battlestation_target_selection_duration_seconds",
		Help:    "Time taken to select targets",
		Buckets: []float64{.001, .005, .01, .025, .05, .1},
	}, []string{"protocol"})

	TargetsProcessed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "battlestation_targets_processed_total",
		Help: "Total number of targets processed",
	}, []string{"protocol", "type"})

	// Ion cannon metrics
	CannonAvailability = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "battlestation_ion_cannon_available",
		Help: "Availability status of ion cannons",
	}, []string{"generation"})

	CannonFireDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "battlestation_ion_cannon_fire_duration_seconds",
		Help:    "Time taken to fire ion cannons",
		Buckets: []float64{.01, .025, .05, .1, .25, .5, 1},
	}, []string{"generation"})

	CannonFireTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "battlestation_ion_cannon_fire_total",
		Help: "Total number of ion cannon fires",
	}, []string{"generation", "status"})

	// Error metrics
	ErrorTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "battlestation_errors_total",
		Help: "Total number of errors",
	}, []string{"type", "operation"})
)

// RecordRequestDuration records the duration of an attack request
func RecordRequestDuration(protocol string, duration float64) {
	RequestDuration.WithLabelValues(protocol).Observe(duration)
}

// RecordRequestComplete records a completed request
func RecordRequestComplete(protocol string, status string) {
	RequestTotal.WithLabelValues(protocol, status).Inc()
}

// RecordTargetSelection records target selection metrics
func RecordTargetSelection(protocol string, duration float64) {
	TargetSelectionDuration.WithLabelValues(protocol).Observe(duration)
}

// RecordTargetProcessed records a processed target
func RecordTargetProcessed(protocol, targetType string) {
	TargetsProcessed.WithLabelValues(protocol, targetType).Inc()
}

// UpdateCannonAvailability updates ion cannon availability
func UpdateCannonAvailability(generation string, available float64) {
	CannonAvailability.WithLabelValues(generation).Set(available)
}

// RecordCannonFire records an ion cannon fire
func RecordCannonFire(generation string, duration float64, status string) {
	CannonFireDuration.WithLabelValues(generation).Observe(duration)
	CannonFireTotal.WithLabelValues(generation, status).Inc()
}

// RecordError records an error
func RecordError(errorType, operation string) {
	ErrorTotal.WithLabelValues(errorType, operation).Inc()
}
