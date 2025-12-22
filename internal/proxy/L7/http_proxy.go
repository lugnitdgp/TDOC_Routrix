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
		http.Error(w, "No Backend avaliable", http.StatusServiceUnavailable)
		return
	}

	backend.Mutex.Lock()
	backend.ActiveConns++
	backend.Mutex.Unlock()

	target, _ := url.Parse("http://" + backend.Address)

	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.ServeHTTP(w, r)

	backend.Mutex.Lock()
	backend.ActiveConns++
	backend.Latency = time.Since(start)
	backend.Mutex.Unlock()
}
