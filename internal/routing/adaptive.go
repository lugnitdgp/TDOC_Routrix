package routing

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/lugnitdgp/TDOC_Routrix/internal/core"
)

type AdaptiveRouter struct{
	pool *core.ServerPool
	rr *RoundRobinRouter
	lc *LeastConnectionsRouter
	rn *RandomRouter
	currentAlgo string
	reason string
	lastPicked string
}
//decision
type Decision struct{
	Time time.Time `json:"time"`
	Algo string `json:"algo"`
	Reason string `json:"reason"`
	Backend string `json:"backend"`
}
var(
	DecisionLog []Decision
	DecisionMu sync.Mutex
)

func NewAdaptiveRouter(pool *core.ServerPool) *AdaptiveRouter{
	return &AdaptiveRouter{
		pool: pool,
		rr: NewRoundRobinRouter(),
		lc: NewLeastConnectionsRouter(),
		rn: NewRandomRouter(),
		currentAlgo: "roundrobin",
		reason: "normal_conditions",
	}
}
func (ar *AdaptiveRouter) Pick() *core.Backend{
	backends := ar.pool.GetServers()
	if len(backends)==0{
		log.Println("[adaptive] no backends in pool")
		return nil
	}
	var totalConns int64
	var totalLatency int64
	var totalErrors int64
	var maxConns int64
	aliveCount :=0

	for _, b := range backends{
		b.Mutex.Lock()
		if b.Alive{
			aliveCount++
			totalConns+=b.ActiveConns
			totalLatency += int64(b.Latency)
			totalErrors += b.ErrorCount

			if b.ActiveConns > maxConns{
				maxConns=b.ActiveConns
			}
		}
		b.Mutex.Unlock()

	}

	if aliveCount == 0{
		log.Println("[adaptive] no alive backends")
		return nil
	}

	avgConns := totalConns/int64(aliveCount)
	avgLatency := totalLatency/int64(aliveCount)
	avgLatencyMs := avgLatency/int64(time.Millisecond)
	errorRate :=float64(totalErrors)/float64(totalConns+1)



log.Printf("Adaptive algo=%s reason=%s picked=%s avgavgConns=%d maxConns=%d avgLatencyMs=%d errorRate=%.2f",
 ar.currentAlgo, ar.reason, ar.lastPicked, avgConns,maxConns, avgLatencyMs,errorRate,)

 var selected *core.Backend

 //decision log
  if errorRate>0.3{
	ar.currentAlgo="random"
	ar.reason=fmt.Sprintf("high_error_rate(%.2f)",errorRate)
	selected=ar.rn.GetNextAvaliableServer(backends)
  }else if maxConns>3{
	ar.currentAlgo="leastconnections"
	ar.reason="high_concurrency"
	selected=ar.lc.GetNextAvaliableServer(backends)
  }else{
	ar.currentAlgo="roundrobin"
	ar.reason="normal_conditions"
	selected=ar.rr.GetNextAvaliableServer(backends)
  }
if selected != nil{
	ar.lastPicked=selected.Address
	DecisionMu.Lock()
	DecisionLog= append(DecisionLog, Decision{
		Time: time.Now(),
		Algo: ar.currentAlgo,
		Reason: ar.reason,
		Backend: selected.Address,
	})
	DecisionMu.Unlock()
}
return selected
}

func (ar *AdaptiveRouter) GetNextAvaliableServer(_ []*core.Backend) *core.Backend{
	return ar.Pick()
}
func(ar *AdaptiveRouter)Name()string{
	return"adaptive"
}
func(ar *AdaptiveRouter)CurrentAlgo() string {return ar.currentAlgo}
func(ar *AdaptiveRouter) Reason() string{return ar.reason}
func(ar *AdaptiveRouter) LastPicked() string{ return ar.lastPicked}