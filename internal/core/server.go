package core

import (
	"sync"
	"time"
)

type Backend struct {
	Address     string
	Weight      int
	Alive       bool
	ActiveConns int64
	Latency     time.Duration
	ErrorCount  int64
	Mutex       sync.Mutex
	Emalatency  int64
}
