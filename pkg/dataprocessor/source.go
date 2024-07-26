package dataprocessor

import "sync"

type source struct {
	sync.Mutex
	ks map[string]KindSource
}

func NewSource() Source {
	return &source{
		ks: make(map[string]KindSource),
	}
}

func (s *source) Set(ks KindSource) {
	s.Lock()
	defer s.Unlock()

	if ks != nil {
		s.ks[ks.Kind()] = ks
	}
}

func (s *source) Kind(kind string) KindSource {
	s.Lock()
	defer s.Unlock()

	return s.ks[kind]
}
