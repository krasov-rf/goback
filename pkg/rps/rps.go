package rps

import (
	"context"
	"sync"

	"golang.org/x/time/rate"
)

type IRps interface {
	Wait(ctx context.Context, clientID string) error
}

type Rps struct {
	// карта лимитов для каждого клиента
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rps      rate.Limit
	burst    int
}

func NewRps(rps rate.Limit, burst int) *Rps {
	return &Rps{
		limiters: make(map[string]*rate.Limiter),
		rps:      rps,
		burst:    burst,
	}
}

func (c *Rps) getLimiter(clientID string) *rate.Limiter {
	c.mu.Lock()
	defer c.mu.Unlock()

	limiter, exists := c.limiters[clientID]
	if !exists {
		limiter = rate.NewLimiter(c.rps, c.burst)
		c.limiters[clientID] = limiter
	}

	return limiter
}

// ожидаем пока
func (c *Rps) Wait(ctx context.Context, clientID string) error {
	limiter := c.getLimiter(clientID)
	if err := limiter.Wait(ctx); err != nil {
		return err
	}
	return nil
}
