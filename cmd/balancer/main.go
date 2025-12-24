package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/lugnitdgp/TDOC_Routrix/internal/api"
	"github.com/lugnitdgp/TDOC_Routrix/internal/core"
	"github.com/lugnitdgp/TDOC_Routrix/internal/health"
	l4 "github.com/lugnitdgp/TDOC_Routrix/internal/proxy/L4"
	l7 "github.com/lugnitdgp/TDOC_Routrix/internal/proxy/L7"
	"github.com/lugnitdgp/TDOC_Routrix/internal/routing"
)

func startDummyBackend(port string) {
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello from backend %s", port)
		})

		log.Printf("[DUMMY] starting backend on :%s\n", port)
		if err := http.ListenAndServe(":"+port, mux); err != nil {
			log.Printf("[DUMMY %s] crashed: %v", port, err)
		}
	}()
}

func main() {
	mode := os.Getenv("LB_MODE")
	if mode == "" {
		mode = "L7"
	}

	// Start dummy backends
	startDummyBackend("9001")
	startDummyBackend("9002")
	startDummyBackend("9003")

	// Backend pool (THREAD SAFE)
	// --------------------------------------------------
	pool := core.NewServerPool()
	pool.AddServer(&core.Backend{Address: "localhost:9001", Alive: true})
	pool.AddServer(&core.Backend{Address: "localhost:9002", Alive: true})
	pool.AddServer(&core.Backend{Address: "localhost:9003", Alive: true})

	router:= routing.NewAdaptiveRouter(pool)

	mux := http.NewServeMux()

	//metrics and admin
	mux.Handle("/metrics",api.MetricsHandler(pool.GetServers()))
	mux.Handle("/admin/add", api.AddServerHandler(pool))

	//status endpoint
	mux.Handle("/status",api.StatusHandler(router,pool.GetServers))

	//health chcker
	checker := &health.Checker{
		Pool: pool,
		Interval: 5*time.Second,
		Timeout: 2 *time.Second,
	}
	checker.Start()

	//L4 mode
	if mode =="L4"{
		log.Println("[MAIN] Starting L4 TCP Load Balancer on :8080")
		tcpProxy := &l4.TCPProxy{ 
			Pool:pool.GetServers(),
			Router: router,
		}
		log.Fatal(tcpProxy.Start("8080"))
		return
	}
	//L7
	log.Println("[MAIN] Starting L7 HTTP Load Balancer on :8080")
	httpProxy := &l7.HTTPProxy{
		Pool: pool.GetServers(),
		Router: router,
	}
	mux.Handle("/", httpProxy)

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET","POST","OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),	
	)(mux)
log.Fatal(http.ListenAndServe(":8080",corsHandler))

}