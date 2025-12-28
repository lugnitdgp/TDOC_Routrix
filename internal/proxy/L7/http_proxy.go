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
		http.Error(w, "No Backend available", http.StatusServiceUnavailable)
		return
	}

	backend.Mutex.Lock()
	backend.ActiveConns++
	backend.Mutex.Unlock()

	target, _ := url.Parse("http://" + backend.Address)
	proxy := httputil.NewSingleHostReverseProxy(target)

	// Track proxy errors
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		backend.Mutex.Lock()
		backend.FailureCount++
		backend.LastFailure = time.Now()
		backend.ActiveConns--
		backend.Mutex.Unlock()

		http.Error(w, "Backend error", http.StatusBadGateway)
	}

	proxy.ServeHTTP(w, r)

	// Successful proxying
	backend.Mutex.Lock()
	backend.FailureCount = 0
	backend.CircuitState = "CLOSED"
	backend.LastSuccess = time.Now()
	backend.ActiveConns--
	backend.Latency = time.Since(start)
	backend.Mutex.Unlock()
}
