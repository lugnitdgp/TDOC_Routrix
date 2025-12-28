package routing

import "github.com/lugnitdgp/TDOC_Routrix/internal/core"

type Router interface {
	GetNextAvaliableServer(backends []*core.Backend) *core.Backend
}
