package dataprocessor

type processor struct {
	tmpl      *Template
	rawPrefix string
}

func NewProcessor(tmpl *Template, rawPrefix string) Processor {
	return &processor{
		tmpl:      tmpl,
		rawPrefix: rawPrefix,
	}
}

func (p *processor) Process(src Source) (Source, error) {
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
					rawKey: NewRef(Ref{
						Kind: "raw." + kind.For,
						Key:  key,
					}),
					kindKey: kind.Name,
					keyKey:  key,

					resourceVersionKey: rv,
				}))
			}
			s.Set(newKs)
		}
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
	return s, nil
}
