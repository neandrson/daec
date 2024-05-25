package orchestrator

import (
	"context"
	"time"
)

type Orchestrator struct {
	cfg Config
}

func New(cfg *Config) *Orchestrator {
	return &Orchestrator{
		cfg: *cfg,
	}
}

func (orch *Orchestrator) Run(ctx context.Context) int {
	time.Sleep(3 * time.Second)
	return 0
}
