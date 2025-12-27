package routing

import (
	"math"

	"github.com/lugnitdgp/TDOC_Routrix/internal/core"
)

type LeastResponseTimeRouter struct {
}

func NewLeastResponseTimeRouter() *LeastResponseTimeRouter {
	return &LeastResponseTimeRouter{}
}
func (lrt *LeastResponseTimeRouter) GetNextAvaliableServer(
	backends []*core.Backend,
) *core.Backend {
	var chosenOne *core.Backend
	var minLatency int64 = math.MaxInt64
	for _, backend := range backends {
		backend.Mutex.Lock()
		alive := backend.Alive
		avgLat := backend.Emalatency
		backend.Mutex.Unlock()
		if !alive {
			continue
		}
		if avgLat < minLatency {
			minLatency = avgLat
			chosenOne = backend
		}
	}
	return chosenOne
}
func (lrt *LeastResponseTimeRouter) Name() string {
	return "Least Response Time"
}
