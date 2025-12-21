package core

import(
	"sync"
)
//serverpool to manage backends
type ServerPool struct{
	servers []*Backend
	mu sync.RWMutex
}
//newserver pool to make a new server pool
func NewServerPool() *ServerPool{
	return &ServerPool{}
}
func(p *ServerPool) AddServer(b *Backend){
	p.mu.Lock()//Lock() is sued in ADD
	defer p.mu.Unlock()
	p.servers =append(p.servers, b)
}
func (p *ServerPool) GetServers() []*Backend{
	p.mu.RLock()//RLock is sued in Get
	defer p.mu.RUnlock()

	return p.servers
}

