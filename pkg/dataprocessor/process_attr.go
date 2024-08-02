package dataprocessor

import (
	"fmt"
)

func (p *processor) checkAttr(src Source, attr *AttrDef) ([]KindSource, error) {
	kinds := attr.Kinds
	if len(kinds) == 0 && attr.Kind != "" {
		kinds = []string{attr.Kind}
	}

	kindSources := make([]KindSource, 0, len(kinds))
	for _, kind := range kinds {
		ks := src.Kind(kind)
		if ks == nil {
			return nil, fmt.Errorf("unknown kind %v for attr %v", kind, attr.Name)
		}
		if ks.HasProperties(attr.Name) {
			return nil, fmt.Errorf("duplicate definition for %v.%v", kind, attr.Name)
		}
		kindSources = append(kindSources, ks)
	}
	return kindSources, nil
}

func (p *processor) processAttr(src Source, attr *AttrDef) error {
	kindSources, err := p.checkAttr(src, attr)
	if err != nil {
		return err
	}

	for _, ks := range kindSources {
		for _, key := range ks.Keys() {
			object := ks.Get(key)
			if object == nil {
				continue
			}
			DefaultProcessAttrForObject(src, attr, object)
			ks.Set(key, object)
		}

		ks.SetProperties(attr.Name)
		src.Set(ks)
	}
	return nil
}

func DefaultProcessAttrForObject(src Source, attr *AttrDef, obj Object) {
	result := EvalValue(src, obj, attr.Value, attr.ValueFrom)
	obj.Set(attr.Name, result)
}
