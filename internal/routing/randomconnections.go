package routing

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/lugnitdgp/TDOC_Routrix/internal/core"
)

type RandomRouter struct {
	mu  sync.Mutex
	rng *rand.Rand
}

func NewRandomRouter() *RandomRouter {
	return &RandomRouter{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (rr *RandomRouter) GetNextAvaliableServer(
	backends []*core.Backend,
) *core.Backend {

	rr.mu.Lock()
	defer rr.mu.Unlock()

	n := len(backends)
	if n == 0 {
		fmt.Println("No Servers Present")
		return nil
	}

	//5 servers 4 severs are 2nd servers

	for i := 0; i < n; i++ {
		idx := rr.rng.Intn(n) //[0 to n)
		backend := backends[idx]

		backend.Mutex.Lock()
		alive := backend.Alive
		backend.Mutex.Unlock()
		if alive {
			return backend
		}
	}
	return nil
}
func (rn *RandomRouter) Name() string {
	return "Random"
}
