package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pigeonligh/kube-indexer/pkg/dataprocessor"
)

type restfulServer struct {
	s *server
}

func (s *restfulServer) getKinds(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, s.s.template.KindList())
}

func (s *restfulServer) eval(ctx *gin.Context) {
	expr := ctx.Query("expr")

	if s.s.data == nil {
		ctx.String(http.StatusBadRequest, "not init")
		return
	}
	result := dataprocessor.EvalValue(s.s.data, dataprocessor.NewObject(nil), nil, &dataprocessor.ValueFrom{
		Expr: &expr,
	})
	result = dataprocessor.UnrefObject(s.s.data, result)
	ctx.JSON(http.StatusOK, result.Value())
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

	r.GET("/api/kinds", s.getKinds)
	r.GET("/api/eval", s.eval)
	r.GET("/api/resource/:kind", s.getKeys)
	r.GET("/api/resource/:kind/*key", s.getObject)

	return r.Run(fmt.Sprintf(":%v", s.s.port))
}
