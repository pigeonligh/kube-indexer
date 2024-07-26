package server

import (
	"context"

	"github.com/pigeonligh/kube-indexer/pkg/cache"
	"github.com/pigeonligh/kube-indexer/pkg/dataprocessor"
)

type Server interface {
	Run(ctx context.Context) error
}

type server struct {
	c         cache.Cache
	data      dataprocessor.Source
	processor dataprocessor.Processor
}

func New(c cache.Cache, template *dataprocessor.Template) Server {
	return &server{
		c:         c,
		data:      dataprocessor.NewSource(),
		processor: dataprocessor.NewProcessor(template, rawPrefix),
	}
}

func (s *server) Run(ctx context.Context) error {
	go func(ctx context.Context) {
		restful := &restfulServer{s}

		if err := restful.Run(ctx); err != nil {
			panic(err)
		}
	}(ctx)

	go func(ctx context.Context) {
		cache := &cacheServer{s}

		if err := cache.Run(ctx); err != nil {
			panic(err)
		}
	}(ctx)

	<-ctx.Done()
	return nil
}
