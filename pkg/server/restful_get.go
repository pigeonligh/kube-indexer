package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pigeonligh/kube-indexer/pkg/dataprocessor"
)

func (s *restfulServer) getKinds(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, s.s.template.KindList())
}

func (s *restfulServer) getKeys(ctx *gin.Context) {
	kind := ctx.Param("kind")

	data := s.s.data
	if data == nil {
		s.responseError(ctx, http.StatusBadRequest, errorNotInit)
		return
	}
	ks := data.Kind(kind)
	if ks == nil {
		s.responseError(ctx, http.StatusBadRequest, errorKindNotFound)
	} else {
		ctx.JSON(http.StatusOK, ks.Keys())
	}
}

func (s *restfulServer) getObject(ctx *gin.Context) {
	kind := ctx.Param("kind")
	key := ctx.Param("key")
	key = strings.TrimPrefix(key, "/")
	raw := ctx.Query("raw")

	data := s.s.data
	if data == nil {
		s.responseError(ctx, http.StatusBadRequest, errorNotInit)
		return
	}
	ks := data.Kind(kind)
	if ks == nil {
		s.responseError(ctx, http.StatusBadRequest, errorKindNotFound)
	} else {
		obj := ks.Get(key)
		if obj == nil {
			s.responseError(ctx, http.StatusBadRequest, errorObjectNotFound)
		} else {
			if raw == "true" {
				obj = obj.Get("_raw")
				obj = dataprocessor.UnrefObject(data, obj)
			}
			ctx.JSON(http.StatusOK, obj.Value())
		}
	}
}
