package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type restfulServer struct {
	s *server
}

func (s *restfulServer) getObject(ctx *gin.Context) {
	kind := ctx.Param("kind")
	ns := ctx.Param("ns")
	name := ctx.Param("name")

	key := name
	if ns != "" {
		key = ns + "/" + name
	}

	if s.s.data == nil {
		ctx.String(http.StatusBadRequest, "not init")
		return
	}
	ks := s.s.data.Kind(kind)
	if ks == nil {
		ctx.String(http.StatusBadRequest, "kind not found")
	} else {
		obj := ks.Get(key)
		if obj == nil {
			ctx.String(http.StatusBadRequest, "object not found")
		} else {
			ctx.JSON(http.StatusOK, obj.Value())
		}
	}
}

func (s *restfulServer) Run(ctx context.Context) error {
	r := gin.Default()

	r.GET("/:kind/namespace/:ns/name/:name", s.getObject)
	r.GET("/:kind/name/:name", s.getObject)

	return r.Run(":8082")
}
