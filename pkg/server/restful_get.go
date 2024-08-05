package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pigeonligh/kube-indexer/pkg/dataprocessor"
)

func (rr *restfulRegisterer) getKinds(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, rr.s.template.KindList())
}

func (rr *restfulRegisterer) getKeys(ctx *gin.Context) {
	kind := ctx.Param("kind")

	data := rr.s.data
	if data == nil {
		rr.responseError(ctx, http.StatusBadRequest, errorNotInit)
		return
	}
	ks := data.Kind(kind)
	if ks == nil {
		rr.responseError(ctx, http.StatusBadRequest, errorKindNotFound)
	} else {
		ctx.JSON(http.StatusOK, ks.Keys())
	}
}

func (rr *restfulRegisterer) getObject(ctx *gin.Context) {
	kind := ctx.Param("kind")
	key := ctx.Param("key")
	key = strings.TrimPrefix(key, "/")
	raw := ctx.Query("raw")

	data := rr.s.data
	if data == nil {
		rr.responseError(ctx, http.StatusBadRequest, errorNotInit)
		return
	}
	ks := data.Kind(kind)
	if ks == nil {
		rr.responseError(ctx, http.StatusBadRequest, errorKindNotFound)
	} else {
		obj := ks.Get(key)
		if obj == nil {
			rr.responseError(ctx, http.StatusBadRequest, errorObjectNotFound)
		} else {
			if raw == "true" {
				obj = obj.Get("_raw")
				obj = dataprocessor.UnrefObject(data, obj)
			}
			ctx.JSON(http.StatusOK, obj.Value())
		}
	}
}
