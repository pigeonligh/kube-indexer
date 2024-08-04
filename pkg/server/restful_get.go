package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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
			ctx.JSON(http.StatusOK, obj.Value())
		}
	}
}
