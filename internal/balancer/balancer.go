package balancer

import (
	"net/url"
	"sync"
	"time"
)

type Backend struct {
	URL         *url.URL
	Alive       bool
	ConnCount   int
	mu          sync.RWMutex
	LastChecked time.Time
}

type Balancer struct {
	backends []*Backend
	mu       sync.Mutex
	index    int
}

func NewBalancer(urls []string) *Balancer {
	b := &Balancer{}
	for _, addr := range urls {
		u, _ := url.Parse(addr)
		b.backends = append(b.backends, &Backend{URL: u, Alive: true})
	}
	return b
}

func (b *Balancer) GetNextBackend() *Backend {
	b.mu.Lock()
	defer b.mu.Unlock()

	n := len(b.backends)
	if n == 0 {
		return nil
	}

	for i := 0; i < n; i++ {
		backend := b.backends[b.index]
		b.index = (b.index + 1) % n

		if backend.IsAlive() {
			return backend
		}
	}

	return nil
}

func (b *Balancer) GetAllBackends() []*Backend {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.backends
}

func (b *Backend) SetAlive(alive bool) {
	b.mu.Lock()
	b.Alive = alive
	b.mu.Unlock()
}

func (b *Backend) IsAlive() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Alive
}
