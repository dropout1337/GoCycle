package GoCycle

import (
	"bufio"
	"os"
	"sync"
	"time"
)

type Cycle struct {
	Mutex  *sync.Mutex
	Locked []string
	List   []string
	I      int

	WaitTime time.Duration
}

func New(List []string) *Cycle {
	return &Cycle{
		WaitTime: 50 * time.Millisecond,
		Mutex:    &sync.Mutex{},
		Locked:   []string{},
		List:     List,
		I:        0,
	}
}

func NewFromFile(Path string) (*Cycle, error) {
	file, err := os.Open(Path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return New(lines), nil
}

// fuck duplicate code i care dont bully me
func (c *Cycle) IsInList(Element string) bool {
	for _, v := range c.List {
		if Element == v {
			return true
		}
	}
	return false
}

func (c *Cycle) IsLocked(Element string) bool {
	for _, v := range c.Locked {
		if Element == v {
			return true
		}
	}
	return false
}

func isInList(List []string, Element string) bool {
	for _, v := range List {
		if Element == v {
			return true
		}
	}
	return false
}

func (c *Cycle) Next() string {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	for {
		c.I++
		if c.I >= len(c.List) {
			c.I = 0
		}

		if !c.IsLocked(c.List[c.I]) {
			return c.List[c.I]
		}

		time.Sleep(c.WaitTime)
	}
}

func (c *Cycle) Lock(Element string) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if c.IsInList(Element) {
		c.Locked = append(c.Locked, Element)
	}
}

func (c *Cycle) Unlock(Element string) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	for i, v := range c.Locked {
		if Element == v {
			c.Locked = append(c.Locked[:i], c.Locked[i+1:]...)
		}
	}
}

func (c *Cycle) ClearDuplicates() {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	var list []string
	for _, v := range c.List {
		if !isInList(list, v) {
			list = append(list, v)
		}
	}
	c.List = list
}

func (c *Cycle) Remove(Element string) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	for i, v := range c.List {
		if Element == v {
			c.List = append(c.List[:i], c.List[i+1:]...)
		}
	}

	for i, v := range c.Locked {
		if Element == v {
			c.Locked = append(c.Locked[:i], c.Locked[i+1:]...)
		}
	}
}

func (c *Cycle) LockByTimeout(Element string, Timeout time.Duration) {
	defer c.Unlock(Element)

	c.Lock(Element)
	time.Sleep(Timeout)
}