package dataprocessor

import (
	"encoding/json"
	"fmt"
)

type ObjectType int

const ()

type Ref struct {
	Kind string `json:"kind"`
	Key  string `json:"key"`
}

type Object interface {
	Ref() *Ref
	RefValue(Source) Object

	Set(string, Object)
	Get(string) Object

	SetIndex(int, Object)
	GetIndex(int) Object
	Range(start, end *int) Object
	Concat() Object
	Push(Object)

	Len() int
	Value() any
}

type Objects []Object

type object struct {
	R   *Ref
	Val any
}

func NewRef(r Ref) Object {
	return &object{
		R: &r,
	}
}

func NewObject(v any) Object {
	if v != nil {
		if o, ok := v.(Object); ok {
			return o
		}
	}
	return &object{Val: v}
}

func (obj *object) Ref() *Ref {
	return obj.R
}

func (obj *object) RefValue(src Source) Object {
	if obj.R != nil {
		ret := src.Kind(obj.R.Kind).Get(obj.R.Key)
		return ret
	}
	return obj
}

func (obj *object) Set(k string, v Object) {
	if obj.R != nil {
		panic("cannot call Set for ref object")
	}

	switch val := obj.Val.(type) {
	case map[string]any:
		val[k] = v

	case Object:
		val.Set(k, v)

	case []any:
		for _, item := range val {
			NewObject(item).Set(k, v)
		}

	case Objects:
		for _, item := range val {
			item.Set(k, v)
		}
	}
}

func (obj *object) SetIndex(i int, v Object) {
	if obj.R != nil {
		panic("cannot call SetIndex for ref object")
	}

	switch val := obj.Val.(type) {
	case Object:
		val.SetIndex(i, v)

	case []any:
		if i >= 0 && i < len(val) {
			val[i] = v
		}

	case Objects:
		if i >= 0 && i < len(val) {
			val[i] = v
		}
	}
}

func (obj *object) Get(k string) Object {
	if obj.R != nil {
		panic("cannot call Get for ref object")
	}

	switch val := obj.Val.(type) {
	case map[string]any:
		return NewObject(val[k])

	case Object:
		return val.Get(k)

	case []any:
		ret := make(Objects, 0)
		for _, item := range val {
			ret = append(ret, NewObject(item).Get(k))
		}
		return NewObject(ret)

	case Objects:
		ret := make(Objects, 0)
		for _, item := range val {
			ret = append(ret, item.Get(k))
		}
		return NewObject(ret)
	}
	return NewObject(nil)
}

func (obj *object) GetIndex(i int) Object {
	if obj.R != nil {
		panic("cannot call GetIndex for ref object")
	}

	switch val := obj.Val.(type) {
	case Object:
		return val.GetIndex(i)

	case []any:
		if i >= 0 && i < len(val) {
			return NewObject(val[i])
		}

	case Objects:
		if i >= 0 && i < len(val) {
			return val[i]
		}
	}
	return NewObject(nil)
}

func (obj *object) Range(start, end *int) Object {
	if obj.R != nil {
		panic("cannot call Range for ref object")
	}

	switch val := obj.Val.(type) {
	case Object:
		return NewObject(val.Range(start, end))

	case []any:
		ret := make(Objects, 0)
		for _, v := range rangeT(val, sliceRange{start, end}) {
			ret = append(ret, NewObject(v))
		}
		return NewObject(ret)

	case Objects:
		ret := make(Objects, 0)
		for _, v := range rangeT(val, sliceRange{start, end}) {
			ret = append(ret, v)
		}
		return NewObject(ret)
	}
	return NewObject(nil)
}

func (obj *object) Concat() Object {
	if obj.R != nil {
		panic("cannot call Concat for ref object")
	}

	switch val := obj.Val.(type) {
	case Object:
		return val.Concat()

	case []any:
		ret := make(Objects, 0)
		for _, subVal := range val {
			subObj := NewObject(subVal)
			arrlength := subObj.Range(nil, nil).Len()
			for i := 0; i < arrlength; i++ {
				ret = append(ret, subObj.GetIndex(i))
			}
		}
		return NewObject(ret)

	case Objects:
		ret := make(Objects, 0)
		for _, subObj := range val {
			arrlength := subObj.Range(nil, nil).Len()
			for i := 0; i < arrlength; i++ {
				ret = append(ret, subObj.GetIndex(i))
			}
		}
		return NewObject(ret)
	}
	return NewObject(nil)
}

func (obj *object) Push(v Object) {
	if obj.R != nil {
		panic("cannot call Range for ref object")
	}

	switch val := obj.Val.(type) {
	case Object:
		val.Push(v)

	case []any:
		val = append(val, v.Value())
		obj.Val = val

	case Objects:
		val = append(val, v)
		obj.Val = val
	}
}

func (obj *object) Len() int {
	if obj.R != nil {
		panic("cannot call Len for ref object")
	}

	switch val := obj.Val.(type) {
	case map[string]any:
		return len(val)

	case Object:
		return val.Len()

	case []any:
		return len(val)

	case Objects:
		return len(val)
	}
	return 0
}

func (obj *object) Value() any {
	if obj.R != nil {
		panic("cannot call Value for ref object")
	}

	if o, ok := obj.Val.(Object); ok {
		return o.Value()
	}
	return obj.Val
}

func (obj *object) String() string {
	if obj.R != nil {
		return fmt.Sprintf("ref(%v,%v)", obj.R.Kind, obj.R.Key)
	}
	return fmt.Sprint(obj.Val)
}

func (obj *object) MarshalJSON() ([]byte, error) {
	if obj.R != nil {
		m := map[string]any{
			"ref": *obj.R,
		}
		return json.Marshal(m)
	}
	return json.Marshal(obj.Val)
}
