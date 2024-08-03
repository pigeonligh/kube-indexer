package dataprocessor

import (
	"fmt"
	"reflect"
)

const (
	rawKey             = "_raw"
	kindKey            = "_kind"
	keyKey             = "_key"
	resourceVersionKey = "_resource_version"
)

const (
	suffixFirstKind  = "X"
	suffixSecondKind = "Y"
)

func (p *processor) checkBind(src Source, bind *BindDef) (KindSource, KindSource, error) {
	if len(bind.Kinds) == 0 || len(bind.Kinds) > 2 {
		return nil, nil, fmt.Errorf("invalid bind kinds %v", bind.Kinds)
	}

	k1 := bind.Kinds[0]
	ks1 := src.Kind(k1)
	if ks1 == nil {
		return nil, nil, fmt.Errorf("unknown kind %v for bind %v", k1, bind.Name)
	}

	if len(bind.Kinds) == 1 {
		if ks1.HasProperties(bind.Name) {
			return nil, nil, fmt.Errorf("duplicate definition for %v.%v", k1, bind.Name)
		}
		return ks1, nil, nil
	}

	if ks1.HasProperties(bind.Name + suffixFirstKind) {
		return nil, nil, fmt.Errorf("duplicate definition for %v.%v", k1, bind.Name+suffixFirstKind)
	}

	k2 := bind.Kinds[1]
	ks2 := src.Kind(k2)
	if ks2 == nil {
		return nil, nil, fmt.Errorf("unknown kind %v for bind %v", k1, bind.Name)
	}
	if ks2.HasProperties(bind.Name + suffixSecondKind) {
		return nil, nil, fmt.Errorf("duplicate definition for %v.%v", k2, bind.Name+suffixSecondKind)
	}

	return ks1, ks2, nil
}

func (p *processor) processBind(src Source, bind *BindDef) error {
	ks1, ks2, err := p.checkBind(src, bind)
	if err != nil {
		return err
	}

	if ks2 == nil {
		p.processBindOneKind(src, bind, ks1)
	} else {
		p.processBindTwoKind(src, bind, ks1, ks2)
	}

	return nil
}

func (p *processor) processBindOneKind(
	src Source, bind *BindDef, ks KindSource,
) {
	setPropertiesForKindSource(ks, bind.Name)

	if bind.Condition != nil {
		// TODO: cannot bind same kind by condition value
		return
	}
	if bind.ConditionFrom != nil {
		if len(bind.ConditionFrom.Matches) > 0 {
			matchBindOneKinds(src, bind.Name, ks, bind.ConditionFrom.Matches)
		}
	}
}

func matchBindOneKinds(src Source, name string, ks KindSource, matches []BindMatch) {
	keys := ks.Keys()
	kind := ks.Kind()

	m := make(map[uint64][]int)
	data := make([]any, len(keys))

	for i, key := range keys {
		object := ks.Get(key)
		if object == nil {
			continue
		}
		values := make([]any, 0)
		hasNull := false
		for _, match := range matches {
			value := EvalValue(src, object, match.FirstValue, match.FirstValueFrom)
			if UnrefObject(src, value).Value() == nil && !match.AllowNull {
				hasNull = true
				break
			}
			values = append(values, value)
		}
		if hasNull {
			continue
		}
		hash := HashValue(values)

		for _, candidateIndex := range m[hash] {
			if reflect.DeepEqual(data[candidateIndex], values) {
				bindObjects(
					name, kind, keys[candidateIndex], ks.Get(keys[candidateIndex]),
					name, kind, key, object,
				)
			}
		}

		m[hash] = append(m[hash], i)
		data[i] = values
	}
}

func (p *processor) processBindTwoKind(
	src Source, bind *BindDef, ks1, ks2 KindSource,
) {
	setPropertiesForKindSource(ks1, bind.Name+suffixFirstKind)
	setPropertiesForKindSource(ks2, bind.Name+suffixSecondKind)

	if bind.Condition != nil {
		if *bind.Condition {
			matchBindTwoKinds(src, bind.Name, ks1, ks2, nil)
		}
	} else if bind.ConditionFrom != nil {
		if len(bind.ConditionFrom.Matches) > 0 {
			matchBindTwoKinds(src, bind.Name, ks1, ks2, bind.ConditionFrom.Matches)
		}
	}
}

func matchBindTwoKinds(src Source, name string, ks1, ks2 KindSource, matches []BindMatch) {
	keys1 := ks1.Keys()
	kind1 := ks1.Kind()
	keys2 := ks2.Keys()
	kind2 := ks2.Kind()

	m := make(map[uint64][]int)
	data := make([]any, len(keys1))

	for i, key := range keys1 {
		object := ks1.Get(key)
		if object == nil {
			continue
		}
		firstValues := make([]any, 0)
		for _, match := range matches {
			firstValue := EvalValue(src, object, match.FirstValue, match.FirstValueFrom)
			firstValues = append(firstValues, firstValue)
		}
		hash1 := HashValue(firstValues)
		m[hash1] = append(m[hash1], i)
		data[i] = firstValues
	}
	for _, key := range keys2 {
		object := ks2.Get(key)
		if object == nil {
			continue
		}
		secondValues := make([]any, 0)
		for _, match := range matches {
			secondValue := EvalValue(src, object, match.SecondValue, match.SecondValueFrom)
			secondValues = append(secondValues, secondValue)
		}
		hash2 := HashValue(secondValues)

		for _, candidateIndex := range m[hash2] {
			if reflect.DeepEqual(data[candidateIndex], secondValues) {
				bindObjects(
					name+suffixFirstKind, kind1, keys1[candidateIndex], ks1.Get(keys1[candidateIndex]),
					name+suffixSecondKind, kind2, key, object,
				)
			}
		}
	}
}

func setPropertiesForKindSource(ks KindSource, name string) {
	for _, key := range ks.Keys() {
		obj := ks.Get(key)
		if obj != nil {
			obj.Set(name, NewObject(Objects{}))
		}
	}
	ks.SetProperties(name)
}

func bindObjects(name1, kind1, key1 string, obj1 Object, name2, kind2, key2 string, obj2 Object) {
	bindList1 := obj1.Get(name1)
	bindList1.Push(NewRef(Ref{Kind: kind2, Key: key2}))
	obj1.Set(name1, bindList1)

	bindList2 := obj2.Get(name2)
	bindList2.Push(NewRef(Ref{Kind: kind1, Key: key1}))
	obj2.Set(name2, bindList2)
}
