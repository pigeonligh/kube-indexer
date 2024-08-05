package dataprocessor

import (
	"github.com/mitchellh/hashstructure/v2"
)

func EvalValue(src Source, obj Object, value any, valueFrom *ValueFrom) Object {
	if value != nil {
		return NewObject(value)
	}
	if valueFrom != nil {
		if valueFrom.Expr != nil {
			result, err := EvalExpr(src, *valueFrom.Expr, map[string]any{
				"cur":    obj,
				"raw":    obj.Get(rawKey),
				"global": src,
			})
			if err != nil {
				return NewObject(nil)
			}
			return NewObject(result)
		}
	}
	return NewObject(nil)
}

func UnrefObject(src Source, obj Object) Object {
	for obj != nil && obj.Ref() != nil {
		obj = obj.RefValue(src)
	}
	if arrlength := obj.Range(nil, nil).Len(); arrlength > 0 {
		ret := make(Objects, 0)
		for i := 0; i < arrlength; i++ {
			ret = append(ret, UnrefObject(src, obj.GetIndex(i)))
		}
		obj = NewObject(ret)
	}
	return obj
}

func HashValue(val any) uint64 {
	r, _ := hashstructure.Hash(val, hashstructure.FormatV2, nil)
	return r
}
