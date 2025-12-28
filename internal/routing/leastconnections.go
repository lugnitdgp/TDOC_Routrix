package routing

import (
	"fmt"
	"sync"

	"github.com/lugnitdgp/TDOC_Routrix/internal/core"
)

type LeastConnectionsRouter struct {
	mu sync.Mutex
}

func NewLeastConnectionsRouter() *LeastConnectionsRouter {
	return &LeastConnectionsRouter{}
}

func (lc *LeastConnectionsRouter) GetNextAvaliableServer(
	backends []*core.Backend,
) *core.Backend {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	n := len(backends)
	if n == 0 {
		fmt.Println("No Servers Present")
		return nil
	}

	var selected *core.Backend
	//cin>> (>> operator & | ) 1101 >>1 0110 <<1  1010
	minConns := int64(^uint64(0) >> 1) //initialises minConns to the maximum number possible

	for _, backend := range backends {
		backend.Mutex.Lock()
		alive := backend.Alive
		active := backend.ActiveConns
		backend.Mutex.Unlock()

		if alive && active < minConns {
			minConns = active
			selected = backend
		}
	}
	if selected != nil {
		selected.Mutex.Lock()
		selected.ActiveConns++
		selected.Mutex.Unlock()
	}
	return selected
}

func (lc *LeastConnectionsRouter) Name() string {
	return "Least Connections"
}