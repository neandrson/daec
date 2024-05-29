package service

import (
	"container/list"
	"fmt"
	"slices"
)

type CalcService struct {
	table     map[string]*Expression
	tasks     []*Task
	taskID    int64
	taskTable map[int64]*list.Element
}

func NewCalcService() *CalcService {
	return &CalcService{
		table:     make(map[string]*Expression),
		taskTable: make(map[int64]*list.Element),
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

	expression, err := NewExpression(id, expr)
	if err != nil {
		return err
	}

	cs.table[id] = expression
	return nil
}

func (cs *CalcService) ListAll() ExpressionList {
	lst := ExpressionList{}
	for _, expr := range cs.table {
		lst.Exprs = append(lst.Exprs, *expr)
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
	return &ExpressionUnit{Expr: *expr}, nil
}

func (cs *CalcService) GetTask() *Task {
	size := cs.updateTasks()
	if size == 0 {
		return nil
	}

	task := cs.tasks[0]
	cs.tasks = cs.tasks[1:]
	return task
}

func (cs *CalcService) updateTasks() int {
	fmt.Println("size of tasks = ", len(cs.tasks))
	if len(cs.tasks) > 0 {
		return len(cs.tasks)
	}
	for i, expr := range cs.table {
		fmt.Println("iteration #", i)
		el := expr.Front()
		for el != nil {
			el1 := el
			fmt.Println(el1.Value.(Token).Type())
			if el1.Value.(Token).Type() != TokenTypeNumber {
				el = el.Next()
				continue
			}
			el2 := el1.Next()
			fmt.Println(el2.Value.(Token).Type())
			if el2 == nil || el2.Value.(Token).Type() != TokenTypeNumber {
				el = el.Next()
				continue
			}
			op := el2.Next()
			if op == nil || op.Value.(Token).Type() != TokenTypeOperation {
				el = el.Next()
				continue
			}
			fmt.Println("tut")
			task := new(Task)
			task.ID = cs.taskID
			cs.taskID++
			taskToken := TaskToken{ID: task.ID}
			cs.taskTable[task.ID] = expr.InsertBefore(&taskToken, el)
			task.Arg1 = el1.Value.(NumToken).Value
			task.Arg2 = el2.Value.(NumToken).Value
			task.Operation = op.Value.(OpToken).Value
			cs.tasks = append(cs.tasks, task)
			el = op.Next()
			expr.Remove(el1)
			expr.Remove(el2)
			expr.Remove(op)
		}
	}
	return len(cs.tasks)
}
