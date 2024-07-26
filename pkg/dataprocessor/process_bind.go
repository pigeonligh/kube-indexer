package dataprocessor

import "fmt"

const (
	kindKey            = "_kind"
	keyKey             = "_key"
	resourceVersionKey = "_resource_version"
)

type processorCondFn func(kind1, key1 string, obj1 Object, kind2, key2 string, obj2 Object) bool

type bindCacheKey struct {
	kindKey1 string
	kindKey2 string
	name     string
}

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
		p.checkBindCount(bind.Name, ks1)
	} else {
		p.processBindDiffKind(src, bind, ks1, ks2)
		p.checkBindCount(bind.Name, ks1)
		p.checkBindCount(bind.Name, ks2)
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
		if bind.ConditionFrom.Expr != nil {
			expr := *bind.ConditionFrom.Expr
			p.condBindSameKinds(bind.Name, ks, p.wrapCondFnByCache(getCondFnBySourceAndExpr(src, expr), bind.Name))
		}
	}
}

func (p *processor) condBindSameKinds(name string, ks KindSource, cond processorCondFn) {
	keys := ks.Keys()
	kind := ks.Kind()
	for j, key2 := range keys {
		object2 := ks.Get(key2)
		if object2 == nil {
			continue
		}

		for i := 0; i < j; i++ {
			key1 := keys[i]
			object1 := ks.Get(key1)
			if object1 == nil {
				continue
			}

			if cond(kind, key1, object1, kind, key2, object2) {
				bindObjects(name, kind, key1, object1, kind, key2, object2)
			}
		}
	}
}

func (p *processor) processBindDiffKind(
	src Source, bind *BindDef, ks1, ks2 KindSource,
) {
	setPropertiesForKindSource(ks1, bind.Name)
	setPropertiesForKindSource(ks2, bind.Name)

	if bind.Condition != nil {
		if *bind.Condition {
			p.condBindDiffKinds(bind.Name, ks1, ks2,
				func(kind1, key1 string, obj1 Object, kind2, key2 string, obj2 Object) bool {
					return true
				},
			)
		}
	} else if bind.ConditionFrom != nil {
		if bind.ConditionFrom.Expr != nil {
			expr := *bind.ConditionFrom.Expr
			p.condBindDiffKinds(bind.Name, ks1, ks2, p.wrapCondFnByCache(getCondFnBySourceAndExpr(src, expr), bind.Name))
		}
	}
}

func (p *processor) condBindDiffKinds(name string, ks1, ks2 KindSource, cond processorCondFn) {
	keys1 := ks1.Keys()
	keys2 := ks2.Keys()
	kind1 := ks1.Kind()
	kind2 := ks2.Kind()

	for _, key1 := range keys1 {
		object1 := ks1.Get(key1)
		if object1 == nil {
			continue
		}

		for _, key2 := range keys2 {
			object2 := ks2.Get(key2)
			if object2 == nil {
				continue
			}

			if cond(kind1, key1, object1, kind2, key2, object2) {
				bindObjects(name, kind1, key1, object1, kind2, key2, object2)
			}
		}
	}
}

func (p *processor) checkBindCount(name string, ks KindSource) {
	for _, key := range ks.Keys() {
		object := ks.Get(key)
		if object == nil {
			continue
		}

		kindkey := getKindKey(ks.Kind(), key)
		bck := bindCacheKey{
			kindKey1: kindkey,
			name:     name,
		}
		newCount := object.Get(name).Len()
		oldCount := p.bindCountCache[bck]
		if newCount != oldCount {
			p.bindCountCache[bck] = newCount
			p.cached[kindkey] = false
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

func getCondFnBySourceAndExpr(src Source, expr string) processorCondFn {
	return func(kind1, key1 string, obj1 Object, kind2, key2 string, obj2 Object) bool {
		result, err := EvalExpr(src, expr, map[string]Object{
			"value1": obj1,
			"value2": obj2,
		})
		if err != nil {
			return false
		}
		b, _ := result.Value().(bool)
		return b
	}
}

func (p *processor) wrapCondFnByCache(cond processorCondFn, name string) processorCondFn {
	return func(kind1, key1 string, obj1 Object, kind2, key2 string, obj2 Object) bool {
		ret := false

		kindKey1 := getKindKey(kind1, key1)
		kindKey2 := getKindKey(kind2, key2)
		bck := bindCacheKey{
			kindKey1: kindKey1,
			kindKey2: kindKey2,
			name:     name,
		}

		oldResult, found := p.bindCache[bck]
		if found && p.cached[kindKey1] && p.cached[kindKey2] {
			ret = oldResult
		} else {
			ret = cond(kind1, key1, obj1, kind2, key2, obj2)
			if ret != oldResult {
				p.cached[kindKey1] = false
				p.cached[kindKey2] = false
			}
		}

		p.newBindCache[bck] = ret
		return ret
	}
}
