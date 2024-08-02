package dataprocessor

import (
	"fmt"
	"reflect"
)

const (
	kindKey            = "_kind"
	keyKey             = "_key"
	resourceVersionKey = "_resource_version"
)

func (p *processor) checkBind(src Source, bind *BindDef) (KindSource, KindSource, error) {
	var k1, k2 string
	switch len(bind.Kinds) {
	case 1:
		k1 = bind.Kinds[0]
		k2 = bind.Kinds[0]

	case 2:
		k1 = bind.Kinds[0]
		k2 = bind.Kinds[1]

	default:
		return nil, nil, fmt.Errorf("invalid bind kinds %v", bind.Kinds)
	}

	ks1 := src.Kind(k1)
	ks2 := src.Kind(k2)
	if ks1 == nil {
		return nil, nil, fmt.Errorf("unknown kind %v for bind %v", k1, bind.Name)
	}
	if ks2 == nil {
		return nil, nil, fmt.Errorf("unknown kind %v for bind %v", k1, bind.Name)
	}

	if ks1.HasProperties(bind.Name) {
		return nil, nil, fmt.Errorf("duplicate definition for %v.%v", k1, bind.Name)
	}
	if ks2.HasProperties(bind.Name) {
		return nil, nil, fmt.Errorf("duplicate definition for %v.%v", k2, bind.Name)
	}

	return ks1, ks2, nil
}

func (p *processor) processBind(src Source, bind *BindDef) error {
	ks1, ks2, err := p.checkBind(src, bind)
	if err != nil {
		return err
	}

	if ks1.Kind() == ks2.Kind() {
		p.processBindSameKind(src, bind, ks1)
	} else {
		p.processBindDiffKind(src, bind, ks1, ks2)
	}

	return nil
}

func (p *processor) processBindSameKind(
	src Source, bind *BindDef, ks KindSource,
) {
	setPropertiesForKindSource(ks, bind.Name)

	if bind.Condition != nil {
		// TODO: cannot bind same kind by condition value
		return
	}
	if bind.ConditionFrom != nil {
		if len(bind.ConditionFrom.Matches) > 0 {
			matchBindSameKinds(src, bind.Name, ks, bind.ConditionFrom.Matches)
		}
	}
}

func matchBindSameKinds(src Source, name string, ks KindSource, matches []BindMatch) {
	keys := ks.Keys()
	kind := ks.Kind()

	m := make(map[uint64][]int)
	data := make([]any, len(keys))

	for i, key := range keys {
		object := ks.Get(key)
		if object == nil {
			continue
		}
		firstValues := make([]any, 0)
		secondValues := make([]any, 0)
		for _, match := range matches {
			firstValue := EvalValue(src, object, match.FirstValue, match.FirstValueFrom)
			firstValues = append(firstValues, firstValue)
			secondValue := EvalValue(src, object, match.SecondValue, match.SecondValueFrom)
			secondValues = append(secondValues, secondValue)
		}
		hash1 := HashValue(firstValues)
		hash2 := HashValue(secondValues)

		for _, candidateIndex := range m[hash2] {
			if reflect.DeepEqual(data[candidateIndex], secondValues) {
				bindObjects(name,
					kind, keys[candidateIndex], ks.Get(keys[candidateIndex]),
					kind, key, object,
				)
			}
		}

		m[hash1] = append(m[hash1], i)
		data[i] = firstValues
	}
}

func (p *processor) processBindDiffKind(
	src Source, bind *BindDef, ks1, ks2 KindSource,
) {
	setPropertiesForKindSource(ks1, bind.Name)
	setPropertiesForKindSource(ks2, bind.Name)

	if bind.Condition != nil {
		if *bind.Condition {
			matchBindDiffKinds(src, bind.Name, ks1, ks2, nil)
		}
	} else if bind.ConditionFrom != nil {
		if len(bind.ConditionFrom.Matches) > 0 {
			matchBindDiffKinds(src, bind.Name, ks1, ks2, bind.ConditionFrom.Matches)
		}
	}
}

func matchBindDiffKinds(src Source, name string, ks1, ks2 KindSource, matches []BindMatch) {
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
				bindObjects(name,
					kind1, keys1[candidateIndex], ks1.Get(keys1[candidateIndex]),
					kind2, key, object,
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

func bindObjects(name string, kind1, key1 string, obj1 Object, kind2, key2 string, obj2 Object) {
	bindList1 := obj1.Get(name)
	bindList1.Push(NewRef(Ref{Kind: kind2, Key: key2}))
	obj1.Set(name, bindList1)

	bindList2 := obj2.Get(name)
	bindList2.Push(NewRef(Ref{Kind: kind1, Key: key1}))
	obj2.Set(name, bindList2)
}
