package server

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

type restfulServer struct {
	s *server
}

func (s *restfulServer) Run(ctx context.Context) error {
	r := gin.Default()

	r.GET("/api/kinds", s.getKinds)
	r.GET("/api/resource/:kind", s.getKeys)
	r.GET("/api/resource/:kind/*key", s.getObject)

	r.POST("/api/eval", s.eval)
	r.POST("/api/resource/:kind", s.listObjects)

	return r.Run(fmt.Sprintf(":%v", s.s.port))
}
