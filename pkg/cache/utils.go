package cache

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

func ConvertObject[T metav1.Object](obj interface{}) (T, bool) {
	tObj, ok := obj.(T)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			return tObj, false
		}
		tObj, ok = tombstone.Obj.(T)
		if !ok {
			return tObj, false
		}
	}
	return tObj, true
}
