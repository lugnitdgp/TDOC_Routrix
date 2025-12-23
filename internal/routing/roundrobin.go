package routing

import (
	"fmt"
	"sync"

	"github.com/lugnitdgp/TDOC_Routrix/internal/core"
)

type RoundRobinRouter struct {
	current int
	mu      sync.Mutex
}

func NewRoundRobinRouter() *RoundRobinRouter {
	return &RoundRobinRouter{
		current: 0,
	}
}
func (rr *RoundRobinRouter) GetNextAvaliableServer(
	backends []*core.Backend,
) *core.Backend {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	n := len(backends)
	if n == 0 {
		fmt.Println("No Servers Present")
		return nil
	}

	for i := 0; i < n; i++ {
		idx := (rr.current + i) % n
		backend := backends[idx]
		backend.Mutex.Lock()
		alive := backend.Alive
		backend.Mutex.Unlock()

		if alive {
			rr.current = (idx + 1) % n //to ensure circular logic
			return backend
		}
	}
	return nil
}
func (rr *RoundRobinRouter) Name() string {
	return "Round Robin"
}
