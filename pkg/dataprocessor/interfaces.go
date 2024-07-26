package dataprocessor

type KindSource interface {
	Kind() string

	Set(k string, v Object)
	Keys() []string
	Get(key string) Object

	SetProperties(p ...string)
	Properties() []string
	HasProperties(string) bool
}

type Source interface {
	Set(s KindSource)

	Kind(k string) KindSource
}

type Processor interface {
	Process(Source) (Source, error)
}
