package server

import (
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/pigeonligh/kube-indexer/pkg/dataprocessor"
)

func (rr *restfulRegisterer) eval(ctx *gin.Context) {
	expr := ctx.Query("expr")

	data := rr.s.data
	if data == nil {
		rr.responseError(ctx, http.StatusBadRequest, errorNotInit)
		return
	}
	result := dataprocessor.EvalValue(data, dataprocessor.NewObject(nil), nil, &dataprocessor.ValueFrom{
		Expr: &expr,
	})
	result = dataprocessor.UnrefObject(data, result)
	ctx.JSON(http.StatusOK, result.Value())
}

func (rr *restfulRegisterer) listObjects(ctx *gin.Context) {
	kind := ctx.Param("kind")

	lparam := listParam{}
	err := ctx.BindJSON(&lparam)
	if err != nil {
		rr.responseError(ctx, http.StatusBadRequest, errorInvalidBody)
		fmt.Println(err)
		return
	}

	data := rr.s.data
	if data == nil {
		rr.responseError(ctx, http.StatusBadRequest, errorNotInit)
		return
	}
	ks := data.Kind(kind)
	if ks == nil {
		rr.responseError(ctx, http.StatusBadRequest, errorKindNotFound)
		return
	}

	filterFunc := getFilterFunc(data, lparam.Filter)
	groupFunc := getGroupFunc(data, lparam.GroupBy)

	groups := make(map[string][]listItem)
	for _, key := range getKeys(data, ks, lparam.From) {
		obj := ks.Get(key)
		if obj == nil {
			continue
		}

		if filterFunc(obj) {
			groupName := groupFunc(obj)
			groups[groupName] = append(groups[groupName], listItem{
				Key:   key,
				Value: obj.Value(),
			})
		}
	}

	result := listResult{
		Kind:       kind,
		Param:      lparam,
		Time:       time.Now().Format(time.RFC3339),
		GroupCount: len(groups),
	}
	for _, kdef := range rr.s.template.Kinds {
		if kdef.Name == kind {
			result.Headers = kdef.Headers
		}
	}
	groupNames := make([]string, 0)
	for groupName := range groups {
		groupNames = append(groupNames, groupName)
	}
	sort.Strings(groupNames)

	for _, groupName := range groupNames {
		result.ResultGroups = append(result.ResultGroups, listGroups{
			Name:  groupName,
			Count: len(groups[groupName]),
			Items: groups[groupName],
		})
	}

	ctx.JSON(http.StatusOK, result)
}
