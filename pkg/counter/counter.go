package counter

import (
	"sync"
)

func NewCounter(value int) *Counter {
	return &Counter{
		value: value,
	}
}

type Counter struct {
	value  int
	locker sync.RWMutex
}

func (c *Counter) Increment() {
	c.locker.Lock()
	defer c.locker.Unlock()
	c.value++
}

func (c *Counter) Decrement() {
	c.locker.Lock()
	defer c.locker.Unlock()
	c.value--
}

func (c *Counter) Value() int {
	c.locker.RLock()
	defer c.locker.RUnlock()
	return c.value
}
