package server

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
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

	data := s.s.data
	if data == nil {
		ctx.String(http.StatusBadRequest, "data not init")
		return
	}
	result := dataprocessor.EvalValue(data, dataprocessor.NewObject(nil), nil, &dataprocessor.ValueFrom{
		Expr: &expr,
	})
	result = dataprocessor.UnrefObject(data, result)
	ctx.JSON(http.StatusOK, result.Value())
}

func (s *restfulServer) getKeys(ctx *gin.Context) {
	kind := ctx.Param("kind")

	data := s.s.data
	if data == nil {
		ctx.String(http.StatusBadRequest, "data not init")
		return
	}
	ks := data.Kind(kind)
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

	data := s.s.data
	if data == nil {
		ctx.String(http.StatusBadRequest, "data not init")
		return
	}
	ks := data.Kind(kind)
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

func (s *restfulServer) listObjects(ctx *gin.Context) {
	kind := ctx.Param("kind")
	m := map[string]string{}
	err := ctx.BindJSON(&m)
	if err != nil {
		ctx.String(http.StatusBadRequest, "invalid body")
		fmt.Println(err)
		return
	}

	data := s.s.data
	if data == nil {
		ctx.String(http.StatusBadRequest, "data not init")
		return
	}
	ks := data.Kind(kind)
	if ks == nil {
		ctx.String(http.StatusBadRequest, "kind not found")
		return
	}

	filter := m["filter"]
	groupBy := m["group_by"]

	filterFunc := getFilterFunc(data, filter)
	groupFunc := getGroupFunc(data, groupBy)

	ret := make(map[string][]dataprocessor.Object)
	for _, key := range ks.Keys() {
		obj := ks.Get(key)
		if obj == nil {
			continue
		}

		if filterFunc(obj) {
			group := groupFunc(obj)
			ret[group] = append(ret[group], obj)
		}
	}
	ctx.JSON(http.StatusOK, ret)
}

func (s *restfulServer) Run(ctx context.Context) error {
	r := gin.Default()

	r.GET("/api/kinds", s.getKinds)
	r.GET("/api/eval", s.eval)
	r.GET("/api/resource/:kind", s.getKeys)
	r.GET("/api/resource/:kind/*key", s.getObject)
	r.POST("/api/resource/:kind", s.listObjects)

	return r.Run(fmt.Sprintf(":%v", s.s.port))
}

type filterFunc func(dataprocessor.Object) bool

func getFilterFunc(src dataprocessor.Source, filter string) filterFunc {
	return func(o dataprocessor.Object) bool {
		if filter == "" {
			return true
		}
		val := dataprocessor.EvalValue(src, o, nil, &dataprocessor.ValueFrom{
			Expr: &filter,
		})
		val = dataprocessor.UnrefObject(src, val)
		return val.Value() == true
	}
}

type groupFunc func(dataprocessor.Object) string

func getGroupFunc(src dataprocessor.Source, groupBy string) groupFunc {
	return func(o dataprocessor.Object) string {
		if groupBy == "" {
			return "-"
		}
		val := dataprocessor.EvalValue(src, o, nil, &dataprocessor.ValueFrom{
			Expr: &groupBy,
		})
		val = dataprocessor.UnrefObject(src, val)
		switch v := val.Value().(type) {
		case string:
			return v

		case int:
			return fmt.Sprint(v)

		case bool:
			return fmt.Sprint(v)

		case []any:
			return "array"

		case map[string]any:
			return "object"

		default:
			if v == nil {
				return "null"
			}
			return reflect.TypeOf(v).String()
		}
	}
}
