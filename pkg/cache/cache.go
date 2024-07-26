package cache

import (
	"context"
	"fmt"
	"sync"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/cache"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	"github.com/pigeonligh/kube-indexer/pkg/dataprocessor"
)

type Cache interface {
	Init()
	WaitForCacheSync(ctx context.Context)
	Run(ctx context.Context) error

	GetObject(kind, ns, name string) (*unstructured.Unstructured, error)
	Snapshot(kindPrefix string) dataprocessor.Source
}

type objectCacheData struct {
	sync.Mutex
	data map[string]*unstructured.Unstructured // ns/name => object
}

type defaultCache struct {
	sync.Mutex

	config    *genericclioptions.ConfigFlags
	informers map[string]cache.SharedInformer // kind => informer
	data      map[string]*objectCacheData     // kind => cache
}

func New(config *genericclioptions.ConfigFlags, kinds ...string) (Cache, error) {
	factory := cmdutil.NewFactory(config)

	informers := make(map[string]cache.SharedInformer)
	for _, kind := range kinds {
		informer, err := NewInformer(factory, kind)
		if err != nil {
			return nil, err
		}
		informers[kind] = informer
	}

	return &defaultCache{
		config:    config,
		informers: informers,
		data:      make(map[string]*objectCacheData),
	}, nil
}

func (c *defaultCache) Init() {
	c.Lock()
	defer c.Unlock()

	for kind, informer := range c.informers {
		oc := &objectCacheData{
			data: make(map[string]*unstructured.Unstructured),
		}
		informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				oc.Lock()
				defer oc.Unlock()

				if u, ok := ConvertObject[*unstructured.Unstructured](obj); ok {
					oc.data[cache.MetaObjectToName(u).String()] = u
				}
			},
			UpdateFunc: func(_, obj interface{}) {
				oc.Lock()
				defer oc.Unlock()

				if u, ok := ConvertObject[*unstructured.Unstructured](obj); ok {
					oc.data[cache.MetaObjectToName(u).String()] = u
				}
			},
			DeleteFunc: func(obj interface{}) {
				oc.Lock()
				defer oc.Unlock()

				if u, ok := ConvertObject[*unstructured.Unstructured](obj); ok {
					delete(oc.data, cache.MetaObjectToName(u).String())
				}
			},
		})
		c.data[kind] = oc
	}
}

func (c *defaultCache) WaitForCacheSync(ctx context.Context) {
	syncs := make([]cache.InformerSynced, 0)
	for _, informer := range c.informers {
		syncs = append(syncs, informer.HasSynced)
	}
	cache.WaitForCacheSync(ctx.Done(), syncs...)
}

func (c *defaultCache) Run(ctx context.Context) error {
	for _, informer := range c.informers {
		go func(informer cache.SharedInformer) {
			informer.Run(ctx.Done())
		}(informer)
	}
	<-ctx.Done()
	return nil
}

func (c *defaultCache) GetObject(kind, ns, name string) (*unstructured.Unstructured, error) {
	c.Lock()
	oc, found := c.data[kind]
	c.Unlock()

	if found {
		oc.Lock()
		defer oc.Unlock()

		key := name
		if ns != "" {
			key = ns + "/" + name
		}

		if obj, found := oc.data[key]; found {
			return obj, nil
		}
		return nil, fmt.Errorf("%v %v not found", kind, key)
	}
	return nil, fmt.Errorf("kind %v not found", kind)
}

func (c *defaultCache) Snapshot(kindPrefix string) dataprocessor.Source {
	c.Lock()
	defer c.Unlock()

	ret := dataprocessor.NewSource()
	for k, oc := range c.data {
		ks := func(k string, oc *objectCacheData) dataprocessor.KindSource {
			oc.Lock()
			defer oc.Unlock()

			ret := dataprocessor.NewKindSource(kindPrefix + k)
			for key, value := range oc.data {
				ret.Set(key, dataprocessor.NewObject(value.Object))
			}
			return ret
		}(k, oc)
		ret.Set(ks)
	}
	return ret
}
