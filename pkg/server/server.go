package server

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/pigeonligh/kube-indexer/pkg/cache"
	"github.com/pigeonligh/kube-indexer/pkg/dataprocessor"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type Server interface {
	Init(ctx context.Context) error
	Run(ctx context.Context) error
}

type server struct {
	config   *genericclioptions.ConfigFlags
	router   gin.IRouter
	template *dataprocessor.Template

	c cache.Cache

	data      dataprocessor.Source
	processor dataprocessor.Processor
}

func New(config *genericclioptions.ConfigFlags, router gin.IRouter, template *dataprocessor.Template) Server {
	return &server{
		config:   config,
		router:   router,
		template: template,

		data:      dataprocessor.NewSource(),
		processor: dataprocessor.NewProcessor(template, rawPrefix),
	}
}

func (s *server) Init(ctx context.Context) error {
	c, err := cache.New(s.config, s.template.ForList()...)
	if err != nil {
		return err
	}
	c.Init()
	s.c = c

	restful := &restfulRegisterer{s}
	restful.Init()

	return nil
}

func (s *server) Run(ctx context.Context) error {
	go func() {
		_ = s.c.Run(ctx)
	}()
	s.c.WaitForCacheSync(ctx)

	go func(ctx context.Context) {
		cache := &cacheServer{s}

		if err := cache.Run(ctx); err != nil {
			panic(err)
		}
	}(ctx)

	<-ctx.Done()
	return nil
}
