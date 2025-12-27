package routing

import (
	"fmt"
	"sync"

	"github.com/lugnitdgp/TDOC_Routrix/internal/core"
)

type WeightedRoundRobinRouter struct {
	current int
	mu      sync.Mutex
}

func NewWeightedRoundRobinRouter() *WeightedRoundRobinRouter {
	return &WeightedRoundRobinRouter{
		current: 0,
	}
}
func (wrr *WeightedRoundRobinRouter) GetNextAvaliableServer(
	backends []*core.Backend,
) *core.Backend {
	wrr.mu.Lock()
	defer wrr.mu.Unlock()

	n := len(backends)
	totalWeight := 0
	var schedule []int
	if n == 0 {
		fmt.Println("No Servers Present")
		return nil
	}
	for _, b := range backends {
		b.Mutex.Lock()
		totalWeight += b.Weight
		b.Mutex.Unlock()
	}
	for _, b := range backends {
		b.Mutex.Lock()
		b.RelativeWeight = float64(b.Weight) / float64(totalWeight)
		b.Mutex.Unlock()
	}
	
	// scheduling all the backends according to there weights
	for i, b := range backends{
		b.Mutex.Lock()
		slots := int(b.RelativeWeight * 10)
		for j := 0; j < slots; j++ {
			schedule = append(schedule, i)
		}
		b.Mutex.Unlock()
	}

	if len(schedule) == 0 {
		fmt.Println("No Servers in Schedule")
		return nil
	}

	for i := 0; i < len(schedule); i++ {
		slotIdx := (wrr.current + i) % len(schedule)
		idx := schedule[slotIdx]
		backend := backends[idx]
		backend.Mutex.Lock()
		alive := backend.Alive
		backend.Mutex.Unlock()

		if alive {
			wrr.current = (slotIdx + 1) % len(schedule) //to ensure circular logic
			return backend
		}
	}

	return nil
}
func (rr *WeightedRoundRobinRouter) Name() string {
	return "Weighted Round Robin"
}
