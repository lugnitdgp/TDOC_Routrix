package health

import (
	"log"
	"net"
	"time"

	"github.com/lugnitdgp/TDOC_Routrix/internal/core"
)

const (
	CircuitCooldown = 10 * time.Second
)

type Checker struct {
	Pool     *core.ServerPool
	Interval time.Duration
	Timeout  time.Duration
}

func (c *Checker) Start() {
	ticker := time.NewTicker(c.Interval)

	go func() {
		for range ticker.C {
			backends := c.Pool.GetServers()
			for _, backend := range backends {
				go c.checkBackend(backend)
			}
		}
	}()
}

func (c *Checker) checkBackend(b *core.Backend) {
	start := time.Now()

	conn, err := net.DialTimeout("tcp", b.Address, c.Timeout)

	b.Mutex.Lock()
	defer b.Mutex.Unlock()

	// ---------------- CIRCUIT BREAKER COOLDOWN ----------------
	if b.CircuitState == "OPEN" {
		if time.Since(b.LastFailure) > CircuitCooldown {
			b.CircuitState = "HALF_OPEN"
			log.Printf("[circuit] %s → HALF_OPEN", b.Address)
		}
	}

	// ---------------- HEALTH CHECK ----------------
	if err != nil {
		// Do NOT flip Alive here — circuit breaker handles failures
		log.Printf("[health] backend unreachable: %s", b.Address)
		return
	}

	_ = conn.Close()

	b.Alive = true
	b.Latency = time.Since(start)
}
