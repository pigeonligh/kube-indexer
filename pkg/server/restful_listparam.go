package server

import (
	"fmt"
	"reflect"

	"github.com/pigeonligh/kube-indexer/pkg/dataprocessor"
)

type listParam struct {
	Filter  string `json:"filter"`
	GroupBy string `json:"group_by"`
}

type listItem struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

type listGroups struct {
	Name  string     `json:"name"`
	Count int        `json:"count"`
	Items []listItem `json:"items"`
}

type listResult struct {
	Kind    string    `json:"kind"`
	Headers []string  `json:"headers"`
	Param   listParam `json:"list_param"`
	Time    string    `json:"time"`

	GroupCount   int          `json:"group_count"`
	ResultGroups []listGroups `json:"result_groups"`
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
			return "{array}"

		case map[string]any:
			return "{object}"

		default:
			if v == nil {
				return "{null}"
			}
			return "{" + reflect.TypeOf(v).String() + "}"
		}
	}
}