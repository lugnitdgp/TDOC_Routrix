package core

import (
	"sync"
	"time"
)

// Backend represents a single backend server
type Backend struct {
	Address     string
	Weight      int
	Alive       bool
	ActiveConns int64
	Latency     time.Duration
	ErrorCount  int64

	// Circuit Breaker fields
	CircuitState string    // CLOSED, OPEN, HALF_OPEN
	FailureCount int64
	LastFailure  time.Time
	LastSuccess  time.Time

	Mutex sync.Mutex
}

// UpdateCircuitState updates the circuit breaker state of the backend
func (b *Backend) UpdateCircuitState() {
	now := time.Now()

	switch b.CircuitState {

	case "CLOSED":
		// Too many failures -> open the circuit
		if b.FailureCount >= 3 {
			b.CircuitState = "OPEN"
			b.LastFailure = now
		}

	case "OPEN":
		// Cooldown period before allowing retry
		if now.Sub(b.LastFailure) > 10*time.Second {
			b.CircuitState = "HALF_OPEN"
		}

	case "HALF_OPEN":
		// Recovery handled by success/failure paths
	}
}
