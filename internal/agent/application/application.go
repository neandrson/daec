package application

import (
	"context"
	"time"

	"github.com/Vojan-Najov/daec/internal/agent/config"
	"github.com/Vojan-Najov/daec/internal/http/client"
	"github.com/Vojan-Najov/daec/internal/result"
	"github.com/Vojan-Najov/daec/internal/task"
)

type Application struct {
	cfg     config.Config
	client  *client.Client
	tasks   chan task.Task
	results chan result.Result
}

var ops map[string]func(float64, float64) float64

func init() {
	ops = make(map[string]func(float64, float64) float64)
	ops["+"] = addition
	ops["-"] = subtraction
	ops["*"] = multiplication
	ops["/"] = division
}

func addition(a, b float64) float64       { return a + b }
func subtraction(a, b float64) float64    { return a - b }
func multiplication(a, b float64) float64 { return a * b }
func division(a, b float64) float64       { return a / b }

func NewApplication(cfg *config.Config) *Application {
	return &Application{
		cfg:     *cfg,
		client:  &client.Client{},
		tasks:   make(chan task.Task),
		results: make(chan result.Result),
	}
}

func (app *Application) Run(ctx context.Context) int {
	defer close(app.results)
	defer close(app.tasks)
	for i := 0; i < app.cfg.ComputingPower; i++ {
		go func(tasks <-chan task.Task, results chan<- result.Result) {
			for {
				select {
				case task, ok := <-tasks:
					if !ok {
						return
					}
					time.Sleep(task.OperationTime)
					value := ops[task.Operation](task.Arg1, task.Arg2)
					results <- result.Result{
						ID:    task.ID,
						Value: value,
					}
				}
			}
		}(app.tasks, app.results)
	}
	for {
		select {
		case <-ctx.Done():
			return 0
		case res := <-app.results:
			app.client.SendResult(res)
		default:
			task := app.client.GetTask()
			if task != nil {
				app.tasks <- *task
			}
		}
	}
	return 0
}
