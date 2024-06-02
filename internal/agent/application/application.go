package application

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/Vojan-Najov/daec/internal/agent/config"
	"github.com/Vojan-Najov/daec/internal/http/client"
	"github.com/Vojan-Najov/daec/internal/result"
	"github.com/Vojan-Najov/daec/internal/task"
	"github.com/Vojan-Najov/daec/pkg/counter"
)

type Application struct {
	cfg     config.Config
	client  *client.Client
	tasks   chan task.Task
	results chan result.Result
	counter *counter.Counter
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
		client:  &client.Client{Hostname: cfg.Hostname, Port: cfg.Port},
		tasks:   make(chan task.Task),
		results: make(chan result.Result),
		counter: counter.NewCounter(cfg.ComputingPower),
	}
}

func (app *Application) Run(ctx context.Context) int {
	defer close(app.results)
	defer close(app.tasks)

	for i := 0; i < app.cfg.ComputingPower; i++ {
		go runWorker(app.tasks, app.results, app.counter)
	}

	for {
		select {
		case <-ctx.Done():
			return 0
		case res := <-app.results:
			app.client.SendResult(res)
		default:
			if app.counter.Value() > 0 {
				task := app.client.GetTask()
				if task != nil {
					fmt.Println("Task getted")
					app.tasks <- *task
				}
			}
		}
	}
}

func runWorker(
	tasks <-chan task.Task,
	results chan<- result.Result,
	counter *counter.Counter,
) {
	for {
		task, ok := <-tasks
		if !ok {
			return
		}

		counter.Decrement()
		time.Sleep(task.OperationTime)

		arg1, err1 := strconv.ParseFloat(task.Arg1, 64)
		arg2, err2 := strconv.ParseFloat(task.Arg2, 64)
		if err1 != nil || err2 != nil {
			results <- result.Result{
				ID:    task.ID,
				Value: fmt.Sprintf("%f", math.NaN),
			}
		} else {
			value := ops[task.Operation](arg1, arg2)
			results <- result.Result{
				ID:    task.ID,
				Value: fmt.Sprintf("%f", value),
			}
		}

		counter.Increment()
	}
}
