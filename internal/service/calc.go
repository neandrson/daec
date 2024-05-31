package service

import (
	"fmt"
	"slices"
	"time"
	"sync"

	"github.com/Vojan-Najov/daec/internal/orchestrator/config"
	"github.com/Vojan-Najov/daec/internal/task"
	"github.com/Vojan-Najov/daec/internal/result"
)

type CalcService struct {
	locker    sync.RWMutex
	exprTable map[string]*Expression
	taskID    int64
	tasks     []*task.Task
	//taskTable map[int64]*list.Element
	taskTable map[int64]ExprElement
	timeTable map[string]time.Duration
}

func NewCalcService(cfg config.Config) *CalcService {
	cs := CalcService{
		exprTable: make(map[string]*Expression),
		taskTable: make(map[int64]ExprElement),
		timeTable: make(map[string]time.Duration),
	}
	cs.timeTable["+"] = cfg.Addtime
	cs.timeTable["-"] = cfg.Subtime
	cs.timeTable["*"] = cfg.Multime
	cs.timeTable["/"] = cfg.Divtime

	return &cs
}

func (cs *CalcService) AddExpression(id, expr string) error {
	if len(id) == 0 {
		return fmt.Errorf("empty ID")
	}
	if len(expr) == 0 {
		return fmt.Errorf("empty expression")
	}

	cs.locker.Lock()
	defer cs.locker.Unlock()

	if _, found := cs.exprTable[id]; found {
		return fmt.Errorf("not a unique ID: %q", id)
	}

	expression, err := NewExpression(id, expr)
	if err != nil {
		return err
	}

	cs.exprTable[id] = expression
	cs.extractTasksFromExpression(expression)
	return nil
}

func (cs *CalcService) ListAll() ExpressionList {
	cs.locker.RLock()
	defer cs.locker.RUnlock()

	lst := ExpressionList{}
	for _, expr := range cs.exprTable {
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
	cs.locker.RLock()
	defer cs.locker.RUnlock()

	expr, found := cs.exprTable[id]
	if !found {
		return nil, fmt.Errorf("id %q not found", id)
	}
	return &ExpressionUnit{Expr: *expr}, nil
}

func (cs *CalcService) GetTask() *task.Task {
	cs.locker.Lock()
	cs.locker.Unlock()
	if len(cs.tasks) == 0 {
		return nil
	}

	task := cs.tasks[0]
	cs.tasks = cs.tasks[1:]
	return task
}

func (cs *CalcService) PutResult(res result.Result) error {
	cs.locker.Lock()
	defer cs.locker.Unlock()

	_, found := cs.taskTable[res.ID]
	if !found {
		return fmt.Errorf("Task id %d not found", res.ID)
	}

	fmt.Println(res)

	el := cs.taskTable[res.ID].Ptr
	exprID := cs.taskTable[res.ID].ID
	delete(cs.taskTable, res.ID)
	expr, found := cs.exprTable[exprID]
	if !found {
		return fmt.Errorf("Expression for task %d not found", res.ID)
	}
	
	fmt.Println("len = ", expr.Len())

	if expr.Len() == 1 {
		expr.Result = fmt.Sprintf("%g", res.Value)
		expr.Status = StatusDone
		expr.Remove(el)
	} else {
		numToken := NumToken{res.Value}
		expr.InsertBefore(numToken, el)
		expr.Remove(el)
		cs.extractTasksFromExpression(expr)		
	}

	return nil
}

/*
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
			task := new(task.Task)
			task.ID = cs.taskID
			cs.taskID++
			taskToken := TaskToken{ID: task.ID}
			cs.taskTable[task.ID] = expr.InsertBefore(&taskToken, el)
			task.Arg1 = el1.Value.(NumToken).Value
			task.Arg2 = el2.Value.(NumToken).Value
			task.Operation = op.Value.(OpToken).Value
			task.Time = cs.timeTable[task.Operation]
			cs.tasks = append(cs.tasks, task)
			el = op.Next()
			expr.Remove(el1)
			expr.Remove(el2)
			expr.Remove(op)
		}
	}
	return len(cs.tasks)
}
*/

func (cs *CalcService) extractTasksFromExpression(expr *Expression) int {
	var taskCount int
	el := expr.Front()
	for el != nil {
		fmt.Println(el.Value)
		el1 := el
		if el1.Value.(Token).Type() != TokenTypeNumber {
			el = el.Next()
			continue
		}

		el2 := el1.Next()
		if el2 == nil || el2.Value.(Token).Type() != TokenTypeNumber {
			el = el.Next()
			continue
		}

		op := el2.Next()
		if op == nil || op.Value.(Token).Type() != TokenTypeOperation {
			el = el.Next()
			continue
		}

		task := new(task.Task)
		task.ID = cs.taskID
		cs.taskID++
		taskToken := TaskToken{ID: task.ID}
		taskElement := expr.InsertBefore(&taskToken, el)
		cs.taskTable[task.ID] = ExprElement{expr.ID, taskElement}
		task.Arg1 = el1.Value.(NumToken).Value
		task.Arg2 = el2.Value.(NumToken).Value
		task.Operation = op.Value.(OpToken).Value
		task.OperationTime = cs.timeTable[task.Operation]

		taskCount++
		cs.tasks = append(cs.tasks, task)
		el = op.Next()
		expr.Remove(el1)
		expr.Remove(el2)
		expr.Remove(op)
	}

	fmt.Printf("Add %d new tasks\n", taskCount)
	return taskCount
}
