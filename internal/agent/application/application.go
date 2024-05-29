package application

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Vojan-Najov/daec/internal/agent/config"
	"github.com/Vojan-Najov/daec/internal/task"
)

type Application struct {
	cfg    config.Config
	client *http.Client
	tasks  chan task.Task
	result chan float64
}

var operations map[string]func(float64, float64) float64

func init() {
	operations = make(map[string]func(float64, float64) float64)
	operations["+"] = addition
	operations["-"] = subtraction
	operations["*"] = multiplication
	operations["/"] = division
}

func addition(a, b float64) float64       { return a + b }
func subtraction(a, b float64) float64    { return a - b }
func multiplication(a, b float64) float64 { return a * b }
func division(a, b float64) float64       { return a / b }

func NewApplication(cfg *config.Config) *Application {
	return &Application{
		cfg:    *cfg,
		client: &http.Client{},
		tasks:  make(chan task.Task),
		result: make(chan float64),
	}
}

func (app *Application) Run(ctx context.Context) int {
	defer close(app.result)
	defer close(app.tasks)
	for i := 0; i < app.cfg.ComputingPower; i++ {
		go func(tasks <-chan task.Task, result chan<- float64) {
			for {
				select {
				case task, ok := <-tasks:
					if !ok {
						return
					}
					// time.Sleep(task.OperationTime)
					result <- operations[task.Operation](
						task.Arg1,
						task.Arg2,
					)
					if !ok {
						return
					}
				}
			}
		}(app.tasks, app.result)
	}
	for {
		select {
		case <-ctx.Done():
			return 0
		case res := <-app.result:
			fmt.Println(res)
			//sendPostRequest(res)
		default:
			app.sendGetRequest()
		}
	}
	return 0
}

func (app *Application) sendGetRequest() {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8081/internal/task", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := app.client.Do(req.WithContext(ctx))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	answer := struct {
		Task task.Task `json:"task"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&answer)
	if err != nil {
		fmt.Println(err)
		return
	}

	app.tasks <- answer.Task
}
