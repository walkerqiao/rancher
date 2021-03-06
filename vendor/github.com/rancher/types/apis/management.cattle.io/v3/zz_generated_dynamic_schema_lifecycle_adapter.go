package v3

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type DynamicSchemaLifecycle interface {
	Create(obj *DynamicSchema) (*DynamicSchema, error)
	Remove(obj *DynamicSchema) (*DynamicSchema, error)
	Updated(obj *DynamicSchema) (*DynamicSchema, error)
}

type dynamicSchemaLifecycleAdapter struct {
	lifecycle DynamicSchemaLifecycle
}

func (w *dynamicSchemaLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*DynamicSchema))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *dynamicSchemaLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*DynamicSchema))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *dynamicSchemaLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*DynamicSchema))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewDynamicSchemaLifecycleAdapter(name string, clusterScoped bool, client DynamicSchemaInterface, l DynamicSchemaLifecycle) DynamicSchemaHandlerFunc {
	adapter := &dynamicSchemaLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *DynamicSchema) (*DynamicSchema, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(*DynamicSchema); ok {
			return o, err
		}
		return nil, err
	}
}
