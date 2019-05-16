/*
Copyright 2019 Rancher Labs.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by main. DO NOT EDIT.

package v1

import (
	"context"

	v1 "github.com/rancher/rio/pkg/apis/rio.cattle.io/v1"
	clientset "github.com/rancher/rio/pkg/generated/clientset/versioned/typed/rio.cattle.io/v1"
	informers "github.com/rancher/rio/pkg/generated/informers/externalversions/rio.cattle.io/v1"
	listers "github.com/rancher/rio/pkg/generated/listers/rio.cattle.io/v1"
	"github.com/rancher/wrangler/pkg/generic"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type AppHandler func(string, *v1.App) (*v1.App, error)

type AppController interface {
	AppClient

	OnChange(ctx context.Context, name string, sync AppHandler)
	OnRemove(ctx context.Context, name string, sync AppHandler)
	Enqueue(namespace, name string)

	Cache() AppCache

	Informer() cache.SharedIndexInformer
	GroupVersionKind() schema.GroupVersionKind

	AddGenericHandler(ctx context.Context, name string, handler generic.Handler)
	AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler)
	Updater() generic.Updater
}

type AppClient interface {
	Create(*v1.App) (*v1.App, error)
	Update(*v1.App) (*v1.App, error)
	UpdateStatus(*v1.App) (*v1.App, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1.App, error)
	List(namespace string, opts metav1.ListOptions) (*v1.AppList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.App, err error)
}

type AppCache interface {
	Get(namespace, name string) (*v1.App, error)
	List(namespace string, selector labels.Selector) ([]*v1.App, error)

	AddIndexer(indexName string, indexer AppIndexer)
	GetByIndex(indexName, key string) ([]*v1.App, error)
}

type AppIndexer func(obj *v1.App) ([]string, error)

type appController struct {
	controllerManager *generic.ControllerManager
	clientGetter      clientset.AppsGetter
	informer          informers.AppInformer
	gvk               schema.GroupVersionKind
}

func NewAppController(gvk schema.GroupVersionKind, controllerManager *generic.ControllerManager, clientGetter clientset.AppsGetter, informer informers.AppInformer) AppController {
	return &appController{
		controllerManager: controllerManager,
		clientGetter:      clientGetter,
		informer:          informer,
		gvk:               gvk,
	}
}

func FromAppHandlerToHandler(sync AppHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1.App
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1.App))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *appController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1.App))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateAppOnChange(updater generic.Updater, handler AppHandler) AppHandler {
	return func(key string, obj *v1.App) (*v1.App, error) {
		if obj == nil {
			return handler(key, nil)
		}

		copyObj := obj.DeepCopy()
		newObj, err := handler(key, copyObj)
		if newObj != nil {
			copyObj = newObj
		}
		if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
			newObj, err := updater(copyObj)
			if newObj != nil && err == nil {
				copyObj = newObj.(*v1.App)
			}
		}

		return copyObj, err
	}
}

func (c *appController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, handler)
}

func (c *appController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), handler)
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, removeHandler)
}

func (c *appController) OnChange(ctx context.Context, name string, sync AppHandler) {
	c.AddGenericHandler(ctx, name, FromAppHandlerToHandler(sync))
}

func (c *appController) OnRemove(ctx context.Context, name string, sync AppHandler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), FromAppHandlerToHandler(sync))
	c.AddGenericHandler(ctx, name, removeHandler)
}

func (c *appController) Enqueue(namespace, name string) {
	c.controllerManager.Enqueue(c.gvk, namespace, name)
}

func (c *appController) Informer() cache.SharedIndexInformer {
	return c.informer.Informer()
}

func (c *appController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *appController) Cache() AppCache {
	return &appCache{
		lister:  c.informer.Lister(),
		indexer: c.informer.Informer().GetIndexer(),
	}
}

func (c *appController) Create(obj *v1.App) (*v1.App, error) {
	return c.clientGetter.Apps(obj.Namespace).Create(obj)
}

func (c *appController) Update(obj *v1.App) (*v1.App, error) {
	return c.clientGetter.Apps(obj.Namespace).Update(obj)
}

func (c *appController) UpdateStatus(obj *v1.App) (*v1.App, error) {
	return c.clientGetter.Apps(obj.Namespace).UpdateStatus(obj)
}

func (c *appController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return c.clientGetter.Apps(namespace).Delete(name, options)
}

func (c *appController) Get(namespace, name string, options metav1.GetOptions) (*v1.App, error) {
	return c.clientGetter.Apps(namespace).Get(name, options)
}

func (c *appController) List(namespace string, opts metav1.ListOptions) (*v1.AppList, error) {
	return c.clientGetter.Apps(namespace).List(opts)
}

func (c *appController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientGetter.Apps(namespace).Watch(opts)
}

func (c *appController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.App, err error) {
	return c.clientGetter.Apps(namespace).Patch(name, pt, data, subresources...)
}

type appCache struct {
	lister  listers.AppLister
	indexer cache.Indexer
}

func (c *appCache) Get(namespace, name string) (*v1.App, error) {
	return c.lister.Apps(namespace).Get(name)
}

func (c *appCache) List(namespace string, selector labels.Selector) ([]*v1.App, error) {
	return c.lister.Apps(namespace).List(selector)
}

func (c *appCache) AddIndexer(indexName string, indexer AppIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1.App))
		},
	}))
}

func (c *appCache) GetByIndex(indexName, key string) (result []*v1.App, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		result = append(result, obj.(*v1.App))
	}
	return result, nil
}