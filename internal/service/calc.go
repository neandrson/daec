package service

import (
	"fmt"
	"slices"

	"github.com/Vojan-Najov/daec/pkg/rpn"
)

const (
	ErrorStatus     = "Error"
	DoneStatus      = "Done"
	InProcessStatus = "In process"
)

type Expression struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Result string `json:"result"`
	rpn    *rpn.RPN
}

type CalcService struct {
	table map[string]Expression
}

// Структура для ответа по запросу на endpoint expressions/{id}
type ExpressionUnit struct {
	Expr Expression `json:"expression`
}

// Структура для ответа по запросу на endpoint expressions
type ExpressionList struct {
	Exprs []Expression `json:"expressions"`
}

type Task struct {
	Id        int64   `json:"id"`
	Arg1      float64 `json:"arg1"`
	Arg2      float64 `json:"arg1"`
	Operation string  `json:"operation"`
}

func NewCalcService() *CalcService {
	return &CalcService{
		table: make(map[string]Expression),
	}
}

func (cs *CalcService) AddExpression(id, expr string) error {
	if len(id) == 0 {
		return fmt.Errorf("empty ID")
	}
	if len(expr) == 0 {
		return fmt.Errorf("empty expression")
	}
	if _, found := cs.table[id]; found {
		return fmt.Errorf("not a unique ID: '%s'", id)
	}

	rpn, err := rpn.NewRPN(expr)
	if err != nil {
		return err
	}

	cs.table[id] = Expression{
		ID:     id,
		Status: InProcessStatus,
		Result: "",
		rpn:    rpn,
	}

	return nil
}

func (cs *CalcService) ListAll() ExpressionList {
	lst := ExpressionList{}
	for _, expr := range cs.table {
		lst.Exprs = append(lst.Exprs, expr)
	}

	slices.SortFunc(lst.Exprs, func(a, b Expression) int {
		if a.ID > b.ID {
			return 1
		} else if a.ID < b.ID {
			return -1
		}
		return 0
	})

	return lst
}

func (cs *CalcService) FindById(id string) (*ExpressionUnit, error) {
	expr, found := cs.table[id]
	if !found {
		return nil, fmt.Errorf("id %s no found", id)
	}
	return &ExpressionUnit{Expr: expr}, nil
}

/*
func (cs *CalcService) GetTask() (*Task) {
	for k, v := range cs.table {
		rpn := &v.rpn
		for e := l.Front(); e != nil; e = e.Next() {
			string(e)
		}
	}
}
*/
