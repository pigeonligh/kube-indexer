package dataprocessor

import "reflect"

type processor struct {
	tmpl      *Template
	rawPrefix string

	rvCache map[string]any  // kindkey => resource_version
	cached  map[string]bool // kindkey => updated in process

	bindCache      map[bindCacheKey]bool
	newBindCache   map[bindCacheKey]bool
	bindCountCache map[bindCacheKey]int

	lastResult Source
}

func NewProcessor(tmpl *Template, rawPrefix string) Processor {
	return &processor{
		tmpl:      tmpl,
		rawPrefix: rawPrefix,

		rvCache:        make(map[string]any),
		bindCache:      make(map[bindCacheKey]bool),
		bindCountCache: make(map[bindCacheKey]int),
	}
}

func (p *processor) Process(src Source) (Source, error) {
	p.cached = make(map[string]bool)
	p.newBindCache = make(map[bindCacheKey]bool)

	lastRvCache := p.rvCache
	p.rvCache = make(map[string]any)

	s := NewSource()
	for _, kind := range p.tmpl.Kinds {
		ks := src.Kind(p.rawPrefix + kind.For)
		if ks != nil {
			s.Set(ks)

			newKs := NewKindSource(kind.Name)
			for _, key := range ks.Keys() {
				obj := ks.Get(key)
				if obj == nil {
					continue
				}
				rv := obj.Get("metadata").Get("resourceVersion").Value()

				newKs.Set(key, NewObject(map[string]any{
					"data": NewRef(Ref{
						Kind: "raw." + kind.For,
						Key:  key,
					}),
					kindKey: kind.Name,
					keyKey:  key,

					resourceVersionKey: rv,
				}))

				kindkey := getKindKey(kind.Name, key)
				lastRv := lastRvCache[kindkey]
				if lastRv != nil && lastRv == rv {
					p.cached[kindkey] = true
				}
				p.rvCache[kindkey] = rv
			}
			s.Set(newKs)
		}
	}

	if reflect.DeepEqual(lastRvCache, p.rvCache) {
		p.newBindCache = p.bindCache
		return p.lastResult, nil
	}

	for _, act := range p.tmpl.Actions {
		if act.Attr != nil {
			err := p.processAttr(s, act.Attr)
			if err != nil {
				return nil, err
			}
		}
		if act.Bind != nil {
			err := p.processBind(s, act.Bind)
			if err != nil {
				return nil, err
			}
		}
	}
	p.bindCache = p.newBindCache
	p.lastResult = s
	return s, nil
}
