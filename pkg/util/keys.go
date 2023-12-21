package util

import (
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
)

// Key struct used to add items to the workqueue. The key is used to identify
// the group/version/Kind of an object, used to index the listers for all
// objects, and the namespace/name key for the object. For deleted objects,
// since they are no longer in the cache, the key stores a shallow copy of the
// deleted object.
type Key struct {
	GvkKey           string
	NamespaceNameKey string
	DeletedObject    *runtime.Object
}

// Given an object that implements runtime.Object, create a key of type Key
// that contains the groupVersionKind key and name/namespace
func KeyForGroupVersionKindNamespaceName(obj any) (Key, error) {
	rObj := obj.(runtime.Object)
	ok := rObj.GetObjectKind()
	gvk := ok.GroupVersionKind()
	namespaceName, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		return Key{}, err
	}
	key := Key{
		GvkKey:           fmt.Sprintf("%s/%s", gvk.GroupVersion().String(), gvk.Kind),
		NamespaceNameKey: namespaceName,
	}

	return key, nil
}

// Create a string key in the form group/version/Kind or version/Kind if the group is empty
func KeyForGroupVersionKind(group, version, kind string) string {
	if group == "" {
		return fmt.Sprintf("%s/%s", version, kind)
	}
	return fmt.Sprintf("%s/%s/%s", group, version, kind)
}

// Used for generating a single string unique representation of the object for logging info
func GenerateObjectInfoString(obj runtime.Object) string {
	group := obj.GetObjectKind().GroupVersionKind().Group
	kind := strings.ToLower(obj.GetObjectKind().GroupVersionKind().Kind)
	mObj := obj.(metav1.Object)

	prefix := kind
	if group != "" {
		prefix = fmt.Sprintf("%s.%s", kind, group)

	}

	return fmt.Sprintf("[%s] %s/%s", mObj.GetNamespace(), prefix, mObj.GetName())
}
