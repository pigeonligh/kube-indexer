package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type restfulServer struct {
	s *server
}

func (s *restfulServer) getKinds(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, s.s.template.KindList())
}

func (s *restfulServer) getKeys(ctx *gin.Context) {
	kind := ctx.Param("kind")

	if s.s.data == nil {
		ctx.String(http.StatusBadRequest, "not init")
		return
	}
	ks := s.s.data.Kind(kind)
	if ks == nil {
		ctx.String(http.StatusBadRequest, "kind not found")
	} else {
		ctx.JSON(http.StatusOK, ks.Keys())
	}
}

func (s *restfulServer) getObject(ctx *gin.Context) {
	kind := ctx.Param("kind")
	key := ctx.Param("key")
	key = strings.TrimPrefix(key, "/")

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

	r.GET("/kinds", s.getKinds)
	r.GET("/:kind", s.getKeys)
	r.GET("/:kind/*key", s.getObject)

	return r.Run(fmt.Sprintf(":%v", s.s.port))
}
