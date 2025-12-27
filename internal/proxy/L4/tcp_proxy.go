package l4

import (
	"io"
	"net"
	"time"

	"github.com/lugnitdgp/TDOC_Routrix/internal/core"
	"github.com/lugnitdgp/TDOC_Routrix/internal/routing"
)

type TCPProxy struct {
	Pool   []*core.Backend
	Router routing.Router
}

func (p *TCPProxy) Start(listenAddr string) error {
	listen, err := net.Listen("tcp", listenAddr)

	if err != nil {
		return err
	}

	for {
		clientConn, err := listen.Accept()
		if err != nil {
			continue
		}

		go p.handleConnection(clientConn)
	}
}

func (p *TCPProxy) handleConnection(clientConn net.Conn) {
	start := time.Now()

	backend := p.Router.GetNextAvaliableServer(p.Pool)

	if backend == nil {
		clientConn.Close()
		return
	}

	backend.Mutex.Lock()
	backend.ActiveConns++
	backend.Mutex.Unlock()

	serverConn, err := net.Dial("tcp", backend.Address)

	if err != nil {
		backend.Mutex.Lock()
		backend.ActiveConns--
		backend.ErrorCount++
		backend.Mutex.Unlock()
		clientConn.Close()
		return
	}

	go io.Copy(serverConn, clientConn)
	io.Copy(clientConn, serverConn)

	clientConn.Close()
	serverConn.Close()
	duration := time.Since(start)
	latency := duration.Milliseconds()

	backend.Mutex.Lock()
	backend.ActiveConns--
	backend.Latency = duration
	if backend.Emalatency == 0 {
		backend.Emalatency = latency
	} else {
		backend.Emalatency = int64(0.5*float64(latency) + 0.5*float64(backend.Emalatency))
	}

	backend.Mutex.Unlock()

}
