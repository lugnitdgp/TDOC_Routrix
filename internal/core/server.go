package core

import (
	"time"
	"sync"
)

type Backend struct {
	Address      string
	Weight       int
	Alive        bool
	ActiveConns  int64
	Latency      time.Duration
	ErrorCount   int64

	// Circuit Breaker fields
	CircuitState string    // CLOSED, OPEN, HALF_OPEN
	FailureCount int64
	LastFailure  time.Time
	LastSuccess  time.Time

	Mutex sync.Mutex
}
