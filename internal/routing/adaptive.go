package routing

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/lugnitdgp/TDOC_Routrix/internal/core"
)

const (
	MaxFailures  = 3
	CooldownTime = 10 * time.Second
)

type AdaptiveRouter struct {
	pool        *core.ServerPool
	rr          *RoundRobinRouter
	lc          *LeastConnectionsRouter
	rn          *RandomRouter
	currentAlgo string
	reason      string
	lastPicked  string
}

// ---------------- Decision Log ----------------
type Decision struct {
	Time    time.Time `json:"time"`
	Algo    string    `json:"algo"`
	Reason  string    `json:"reason"`
	Backend string    `json:"backend"`
}

var (
	DecisionLog []Decision
	DecisionMu  sync.Mutex
)

// ---------------- Constructor ----------------
func NewAdaptiveRouter(pool *core.ServerPool) *AdaptiveRouter {
	return &AdaptiveRouter{
		pool:        pool,
		rr:          NewRoundRobinRouter(),
		lc:          NewLeastConnectionsRouter(),
		rn:          NewRandomRouter(),
		currentAlgo: "roundrobin",
		reason:      "normal_conditions",
	}
}

// ---------------- Core Logic ----------------
func (ar *AdaptiveRouter) Pick() *core.Backend {
	backends := ar.pool.GetServers()
	if len(backends) == 0 {
		log.Println("[adaptive] no backends in pool")
		return nil
	}

	var (
		totalConns   int64
		totalLatency int64
		totalErrors  int64
		maxConns     int64
		aliveCount   int
	)

	now := time.Now()

	for _, b := range backends {
		b.Mutex.Lock()

		// ---------- CIRCUIT BREAKER ----------
		if b.CircuitState == "OPEN" {
			if now.Sub(b.LastFailure) > CooldownTime {
				b.CircuitState = "HALF_OPEN"
				log.Printf("[circuit] %s → HALF_OPEN", b.Address)
			} else {
				b.Mutex.Unlock()
				continue // still OPEN
			}
		}

		// Too many failures → OPEN
		if b.FailureCount >= MaxFailures {
			b.CircuitState = "OPEN"
			b.LastFailure = now
			log.Printf("[circuit] %s → OPEN", b.Address)
			b.Mutex.Unlock()
			continue
		}

		// ---------- METRICS COLLECTION ----------
		// CLOSED and HALF_OPEN are both allowed
		if b.Alive {
			aliveCount++
			totalConns += b.ActiveConns
			totalLatency += int64(b.Latency)
			totalErrors += b.ErrorCount

			if b.ActiveConns > maxConns {
				maxConns = b.ActiveConns
			}
		}

		b.Mutex.Unlock()
	}

	if aliveCount == 0 {
		log.Println("[adaptive] no alive backends")
		return nil
	}

	avgConns := totalConns / int64(aliveCount)
	avgLatency := totalLatency / int64(aliveCount)
	avgLatencyMs := avgLatency / int64(time.Millisecond)
	errorRate := float64(totalErrors) / float64(totalConns+1)

	log.Printf(
		"[adaptive] algo=%s reason=%s avgConns=%d maxConns=%d avgLatencyMs=%d errorRate=%.2f",
		ar.currentAlgo, ar.reason, avgConns, maxConns, avgLatencyMs, errorRate,
	)

	var selected *core.Backend

	// ---------- ALGO SELECTION ----------
	if errorRate > 0.3 {
		ar.currentAlgo = "random"
		ar.reason = fmt.Sprintf("high_error_rate(%.2f)", errorRate)
		selected = ar.rn.GetNextAvaliableServer(backends)

	} else if maxConns > 3 {
		ar.currentAlgo = "leastconnections"
		ar.reason = "high_concurrency"
		selected = ar.lc.GetNextAvaliableServer(backends)

	} else {
		ar.currentAlgo = "roundrobin"
		ar.reason = "normal_conditions"
		selected = ar.rr.GetNextAvaliableServer(backends)
	}

	if selected != nil {
		ar.lastPicked = selected.Address

		DecisionMu.Lock()
		DecisionLog = append(DecisionLog, Decision{
			Time:    time.Now(),
			Algo:    ar.currentAlgo,
			Reason:  ar.reason,
			Backend: selected.Address,
		})
		DecisionMu.Unlock()
	}

	return selected
}

// ---------------- Interface Methods ----------------
func (ar *AdaptiveRouter) GetNextAvaliableServer(_ []*core.Backend) *core.Backend {
	return ar.Pick()
}

func (ar *AdaptiveRouter) Name() string {
	return "adaptive"
}

func (ar *AdaptiveRouter) CurrentAlgo() string {
	return ar.currentAlgo
}

func (ar *AdaptiveRouter) Reason() string {
	return ar.reason
}

func (ar *AdaptiveRouter) LastPicked() string {
	return ar.lastPicked
}
