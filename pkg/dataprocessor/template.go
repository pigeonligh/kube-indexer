package dataprocessor

type KindDef struct {
	Name string `yaml:"name"`
	For  string `yaml:"for"`
}

type AttrValueFrom struct {
	Expr *string `yaml:"expr"`
}

type AttrDef struct {
	Kind  string   `yaml:"kind"`
	Kinds []string `yaml:"kinds"`
	Name  string   `yaml:"name"`

	Value     any            `yaml:"value"`
	ValueFrom *AttrValueFrom `yaml:"valueFrom"`
}

type BindConditionFrom struct {
	Expr *string `yaml:"expr"`
}

type BindDef struct {
	Kinds []string `yaml:"kinds"`
	Name  string   `yaml:"name"`

	Condition     *bool              `yaml:"condition"`
	ConditionFrom *BindConditionFrom `yaml:"conditionFrom"`
}

type Action struct {
	Attr *AttrDef `yaml:"attr"`
	Bind *BindDef `yaml:"bind"`
}

type Template struct {
	Kinds   []KindDef `yaml:"kinds"`
	Actions []Action  `yaml:"actions"`
}

func (t Template) CacheFor() []string {
	ret := make([]string, 0)
	for _, kind := range t.Kinds {
		ret = append(ret, kind.For)
	}
	return ret
}
