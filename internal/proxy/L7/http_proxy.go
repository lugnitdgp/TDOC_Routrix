package l7

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/lugnitdgp/TDOC_Routrix/internal/core"
	"github.com/lugnitdgp/TDOC_Routrix/internal/routing"
)

type HTTPProxy struct {
	Pool   []*core.Backend
	Router routing.Router
}

func (p *HTTPProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	backend := p.Router.GetNextAvaliableServer(p.Pool)
	if backend == nil {
		http.Error(w, "No backend available", http.StatusServiceUnavailable)
		return
	}

	backend.Mutex.Lock()
	backend.ActiveConns++
	backend.Mutex.Unlock()

	target, _ := url.Parse("http://" + backend.Address)
	proxy := httputil.NewSingleHostReverseProxy(target)

	// -------- ERROR PATH --------
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		backend.Mutex.Lock()

		backend.FailureCount++
		backend.LastFailure = time.Now()

		if backend.CircuitState == "HALF_OPEN" || backend.FailureCount >= 3 {
			backend.CircuitState = "OPEN"
		}

		backend.ActiveConns--
		backend.Mutex.Unlock()

		http.Error(w, "Backend error", http.StatusBadGateway)
	}

	proxy.ServeHTTP(w, r)

	// -------- SUCCESS PATH --------
	backend.Mutex.Lock()

	if backend.CircuitState == "HALF_OPEN" {
		backend.CircuitState = "CLOSED"
		backend.FailureCount = 0
	}

	backend.LastSuccess = time.Now()
	backend.Latency = time.Since(start)
	backend.ActiveConns--

	backend.Mutex.Unlock()
}
