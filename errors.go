package xtpl

import (
	"fmt"
	"strings"
	"sync"
)

type errors struct {
	s    sync.Mutex
	list []string
}

func (e *errors) Add(err error) string {
	e.s.Lock()
	number := serialID.Next()
	e.list = append(e.list, fmt.Sprintf("Runtime error #%d: %s\n", number, err.Error()))
	e.s.Unlock()
	return fmt.Sprintf("[Runtime error: See #%d]", number)
}

func (e *errors) Error() error {
	if len(e.list) == 0 {
		return nil
	}
	return fmt.Errorf(strings.Join(e.list, "\n"))
}
