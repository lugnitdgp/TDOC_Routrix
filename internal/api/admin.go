package api

import (
	"net/http"

	"github.com/lugnitdgp/TDOC_Routrix/internal/core"
)

func AddServerHandler(pool *core.ServerPool) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		addr := r.URL.Query().Get("address")

		if addr == " " {
			http.Error(w, "Backend address is not found", http.StatusBadRequest)
			return
		}

		pool.AddServer(&core.Backend{Address: addr, Alive: true})

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("server added"))
	}
}
