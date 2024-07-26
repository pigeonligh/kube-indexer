package server

import (
	"context"
	"fmt"
	"time"
)

const (
	rawPrefix = "raw."
)

type cacheServer struct {
	s *server
}

func (s *cacheServer) sync() {
	start := time.Now()

	src := s.s.c.Snapshot(rawPrefix)

	src, err := s.s.processor.Process(src)
	if err != nil {
		fmt.Println(err)
		return
	}
	s.s.data = src

	t := time.Since(start)
	_ = t // sync time
}

func (s *cacheServer) Run(ctx context.Context) error {
	s.sync()
	fmt.Println("first synced")

	t := time.NewTicker(time.Second * 5)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil

		case <-t.C:
			s.sync()
		}
	}
}
