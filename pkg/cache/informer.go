package cache

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/tools/cache"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func NewInformer(factory cmdutil.Factory, kind string) (cache.SharedInformer, error) {
	result := factory.NewBuilder().
		Unstructured().
		ResourceTypeOrNameArgs(true, kind).
		Do()
	if err := result.Err(); err != nil {
		return nil, err
	}
	infos, err := result.Infos()
	if err != nil {
		return nil, err
	}
	if len(infos) != 1 {
		return nil, fmt.Errorf("invalid gvk")
	}
	info := infos[0]

	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(info.Mapping.GroupVersionKind)

	informer := cache.NewSharedInformer(&cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return resource.NewHelper(info.Client, info.Mapping).List("", obj.GetAPIVersion(), &options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return resource.NewHelper(info.Client, info.Mapping).Watch("", obj.GetAPIVersion(), &options)
		},
	}, obj, 0)
	return informer, nil
}
