package xtpl

import "sync"

type serial struct {
	sync.Mutex
	n int64
}

func (s *serial) Next() (n int64) {
	s.Lock()
	s.n++
	n = s.n
	s.Unlock()
	return
}
