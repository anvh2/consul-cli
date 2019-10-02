package consul

import (
	"sync"

	"go.uber.org/zap"
)

// Health -
type Health struct {
	mu   sync.Mutex
	endp []*Endpoint

	logger zap.Logger
}

// Endpoint -
type Endpoint struct {
}
