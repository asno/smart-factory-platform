package simulation

import (
	"sync"

	"smart-factory/machine-simulator/internal/machines"
)

type Registry struct {
	machines []machines.RuntimeMachine
	mutex    sync.RWMutex
}

func NewRegistry() *Registry {

	return &Registry{
		machines: []machines.RuntimeMachine{},
	}
}

func (r *Registry) Add(
	m machines.RuntimeMachine,
) {

	r.mutex.Lock()

	defer r.mutex.Unlock()

	r.machines = append(
		r.machines,
		m,
	)
}

func (r *Registry) GetMachines() []machines.RuntimeMachine {

	r.mutex.RLock()

	defer r.mutex.RUnlock()

	return r.machines
}
