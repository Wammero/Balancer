package health

import (
	"net/http"
	"sync"
	"time"

	"github.com/Wammero/Balancer/internal/balancer"
)

func HealthChecker(b *balancer.Balancer) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			var wg sync.WaitGroup
			for _, backend := range b.GetAllBackends() {
				wg.Add(1)
				go func(be *balancer.Backend) {
					defer wg.Done()
					checkHealth(be)
				}(backend)
			}
			wg.Wait()
		}
	}
}

func checkHealth(b *balancer.Backend) {
	timeout := time.Second
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(b.URL.String() + "/health")
	if err != nil || resp.StatusCode != 200 {
		b.SetAlive(false)
		return
	}
	b.SetAlive(true)
}
