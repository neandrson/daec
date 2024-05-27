// В этом пакетк содержится код обработчиков http запросов.
package handler

import (
	"context"
	"net/http"

	"github.com/Vojan-Najov/daec/internal/service"
)

// тип Decorator служит для добавления middleware к обработчикам
type Decorator func(http.Handler) http.Handler

// объект для обработки запросов
type calcStates struct {
	CalcService *service.CalcService
}

func NewHandler(
	ctx context.Context,
	calcService *service.CalcService,
) (http.Handler, error) {
	serveMux := http.NewServeMux()

	calcState := calcStates{
		CalcService: calcService,
	}

	serveMux.HandleFunc("/api/v1/calculate", calcState.calculate)
	serveMux.HandleFunc("/api/v1/expressions", calcState.listAll)
	serveMux.HandleFunc("/api/v1/expressions/{id}", calcState.listByID)
	serveMux.HandleFunc("/internal/task", calcState.sendTask)

	return serveMux, nil
}

// функция добавления middleware
func Decorate(next http.Handler, ds ...Decorator) http.Handler {
	decorated := next
	for d := len(ds) - 1; d >= 0; d-- {
		decorated = ds[d](decorated)
	}

	return decorated
}

// Добавление вычисления арифметического выражения
func (ls *calcStates) calculate(w http.ResponseWriter, r *http.Request) {
	//worldState := ls.LifeService.NewState()

	//err := json.NewEncoder(w).Encode(worldState.Cells)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//}
}

func (ls *calcStates) listAll(w http.ResponseWriter, r *http.Request) {
}

func (ls *calcStates) listByID(w http.ResponseWriter, r *http.Request) {
}

func (ls *calcStates) sendTask(w http.ResponseWriter, r *http.Request) {
}
