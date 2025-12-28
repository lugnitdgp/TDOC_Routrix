package api

import (
	"encoding/json"
	"net/http"

	"github.com/lugnitdgp/TDOC_Routrix/internal/core"
)

func MetricsHandler(pool []*core.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(pool)
	}
}