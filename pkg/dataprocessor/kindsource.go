package dataprocessor

import (
	"sort"
	"sync"
)

type kindSource struct {
	sync.Mutex
	k          string
	properties map[string]struct{}
	data       map[string]Object
}

func NewKindSource(k string) KindSource {
	return &kindSource{
		k:          k,
		properties: make(map[string]struct{}),
		data:       make(map[string]Object),
	}
}

func CloneKindSource(k string, s KindSource) KindSource {
	ks := NewKindSource(k)
	for _, key := range s.Keys() {
		if v := s.Get(key); v != nil {
			ks.Set(key, v)
		}
	}
	ks.SetProperties(s.Properties()...)
	return ks
}

func (s *kindSource) Set(k string, v Object) {
	s.Lock()
	defer s.Unlock()

	s.data[k] = v
}

func (s *kindSource) Kind() string {
	return s.k
}

func (s *kindSource) Keys() []string {
	s.Lock()
	defer s.Unlock()

	keys := make([]string, 0)
	for k := range s.data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (s *kindSource) Get(key string) Object {
	s.Lock()
	defer s.Unlock()

	return s.data[key]
}

func (s *kindSource) SetProperties(properties ...string) {
	s.Lock()
	defer s.Unlock()

	for _, p := range properties {
		s.properties[p] = struct{}{}
	}
}

func (s *kindSource) Properties() []string {
	s.Lock()
	defer s.Unlock()

	properties := make([]string, 0)
	for p := range s.properties {
		properties = append(properties, p)
	}
	return properties
}

func (s *kindSource) HasProperties(p string) bool {
	_, found := s.properties[p]
	return found
}
