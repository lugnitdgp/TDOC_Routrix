package api

import (
	"encoding/json"
	"net/http"

	"github.com/lugnitdgp/TDOC_Routrix/internal/core"
	"github.com/lugnitdgp/TDOC_Routrix/internal/routing"
)

type StatusResponse struct{
	CurrentAlgo string `json:"current_algo"`
	AdaptiveReason string `json:"adaptive_reason"`
	SelectedBackend string `json:"selected_backend"`
	Backends []*core.Backend `json:"backends"`
	DecisionLog []routing.Decision `json:"decision_log"`
}

func StatusHandler (router *routing.AdaptiveRouter, getBackends func() []*core.Backend) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		routing.DecisionMu.Lock()
		logs:= make([]routing.Decision, len(routing.DecisionLog))
		copy(logs, routing.DecisionLog)
		routing.DecisionMu.Unlock()

		//construct response object
		resp:= StatusResponse{
			CurrentAlgo: router.CurrentAlgo(),
			AdaptiveReason: router.Reason(),
			SelectedBackend: router.LastPicked(),
			Backends: getBackends(),
			DecisionLog: logs,
		}

		w.Header().Set("Content-Type","application/json")
		json.NewEncoder(w).Encode(resp)
	})


}