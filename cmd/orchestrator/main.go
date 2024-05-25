package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Vojan-Najov/daec/internal/orchestrator"
)

func main() {
	cfg, err := orchestrator.NewConfigFromEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	app := orchestrator.New(cfg)
	exitCode := app.Run(ctx)
	os.Exit(exitCode)
}
