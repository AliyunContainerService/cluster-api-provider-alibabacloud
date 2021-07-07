/*
Copyright 2018 The Kubernetes Authors.

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

package client

import (
	"context"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

var _ Reader = &unstructuredClient{}
var _ Writer = &unstructuredClient{}
var _ StatusWriter = &unstructuredClient{}

// client is a client.Client that reads and writes directly from/to an API server.  It lazily initializes
// new clients at the time they are used, and caches the client.
type unstructuredClient struct {
	cache      *clientCache
	paramCodec runtime.ParameterCodec
}

// Create implements client.Client
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
func (uc *unstructuredClient) Create(ctx context.Context, obj Object, opts ...CreateOption) error {
=======
func (uc *unstructuredClient) Create(_ context.Context, obj runtime.Object, opts ...CreateOption) error {
>>>>>>> 79bfea2d (update vendor)
=======
func (uc *unstructuredClient) Create(ctx context.Context, obj Object, opts ...CreateOption) error {
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
func (uc *unstructuredClient) Create(ctx context.Context, obj Object, opts ...CreateOption) error {
>>>>>>> 03397665 (update api)
	u, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return fmt.Errorf("unstructured client did not understand object: %T", obj)
	}

	gvk := u.GroupVersionKind()

	o, err := uc.cache.getObjMeta(obj)
	if err != nil {
		return err
	}

	createOpts := &CreateOptions{}
	createOpts.ApplyOptions(opts)
	result := o.Post().
		NamespaceIfScoped(o.GetNamespace(), o.isNamespaced()).
		Resource(o.resource()).
		Body(obj).
		VersionedParams(createOpts.AsCreateOptions(), uc.paramCodec).
		Do(ctx).
		Into(obj)

	u.SetGroupVersionKind(gvk)
	return result
}

// Update implements client.Client
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
func (uc *unstructuredClient) Update(ctx context.Context, obj Object, opts ...UpdateOption) error {
=======
func (uc *unstructuredClient) Update(_ context.Context, obj runtime.Object, opts ...UpdateOption) error {
>>>>>>> 79bfea2d (update vendor)
=======
func (uc *unstructuredClient) Update(ctx context.Context, obj Object, opts ...UpdateOption) error {
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
func (uc *unstructuredClient) Update(ctx context.Context, obj Object, opts ...UpdateOption) error {
>>>>>>> 03397665 (update api)
	u, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return fmt.Errorf("unstructured client did not understand object: %T", obj)
	}

	gvk := u.GroupVersionKind()

	o, err := uc.cache.getObjMeta(obj)
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 03397665 (update api)
	if err != nil {
		return err
	}

	updateOpts := UpdateOptions{}
	updateOpts.ApplyOptions(opts)
	result := o.Put().
		NamespaceIfScoped(o.GetNamespace(), o.isNamespaced()).
		Resource(o.resource()).
		Name(o.GetName()).
		Body(obj).
		VersionedParams(updateOpts.AsUpdateOptions(), uc.paramCodec).
		Do(ctx).
		Into(obj)

	u.SetGroupVersionKind(gvk)
	return result
}

// Delete implements client.Client
func (uc *unstructuredClient) Delete(ctx context.Context, obj Object, opts ...DeleteOption) error {
	_, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return fmt.Errorf("unstructured client did not understand object: %T", obj)
	}

	o, err := uc.cache.getObjMeta(obj)
<<<<<<< HEAD
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
>>>>>>> 03397665 (update api)
	if err != nil {
		return err
	}

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 03397665 (update api)
	deleteOpts := DeleteOptions{}
	deleteOpts.ApplyOptions(opts)
	return o.Delete().
		NamespaceIfScoped(o.GetNamespace(), o.isNamespaced()).
		Resource(o.resource()).
		Name(o.GetName()).
		Body(deleteOpts.AsDeleteOptions()).
		Do(ctx).
		Error()
<<<<<<< HEAD
=======
	updateOpts := UpdateOptions{}
	updateOpts.ApplyOptions(opts)
	result := o.Put().
		NamespaceIfScoped(o.GetNamespace(), o.isNamespaced()).
		Resource(o.resource()).
		Name(o.GetName()).
		Body(obj).
		VersionedParams(updateOpts.AsUpdateOptions(), uc.paramCodec).
		Do(ctx).
		Into(obj)

	u.SetGroupVersionKind(gvk)
	return result
>>>>>>> e879a141 (alibabacloud machine-api provider)
}

<<<<<<< HEAD
// DeleteAllOf implements client.Client
func (uc *unstructuredClient) DeleteAllOf(ctx context.Context, obj Object, opts ...DeleteAllOfOption) error {
	_, ok := obj.(*unstructured.Unstructured)
=======
// Delete implements client.Client
<<<<<<< HEAD
func (uc *unstructuredClient) Delete(_ context.Context, obj runtime.Object, opts ...DeleteOption) error {
	u, ok := obj.(*unstructured.Unstructured)
>>>>>>> 79bfea2d (update vendor)
=======
func (uc *unstructuredClient) Delete(ctx context.Context, obj Object, opts ...DeleteOption) error {
	_, ok := obj.(*unstructured.Unstructured)
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
}

// DeleteAllOf implements client.Client
func (uc *unstructuredClient) DeleteAllOf(ctx context.Context, obj Object, opts ...DeleteAllOfOption) error {
	_, ok := obj.(*unstructured.Unstructured)
>>>>>>> 03397665 (update api)
	if !ok {
		return fmt.Errorf("unstructured client did not understand object: %T", obj)
	}

	o, err := uc.cache.getObjMeta(obj)
	if err != nil {
		return err
	}
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 03397665 (update api)

	deleteAllOfOpts := DeleteAllOfOptions{}
	deleteAllOfOpts.ApplyOptions(opts)
	return o.Delete().
		NamespaceIfScoped(deleteAllOfOpts.ListOptions.Namespace, o.isNamespaced()).
		Resource(o.resource()).
		VersionedParams(deleteAllOfOpts.AsListOptions(), uc.paramCodec).
		Body(deleteAllOfOpts.AsDeleteOptions()).
		Do(ctx).
		Error()
}

// Patch implements client.Client
func (uc *unstructuredClient) Patch(ctx context.Context, obj Object, patch Patch, opts ...PatchOption) error {
	_, ok := obj.(*unstructured.Unstructured)
<<<<<<< HEAD
=======
=======

>>>>>>> e879a141 (alibabacloud machine-api provider)
	deleteOpts := DeleteOptions{}
	deleteOpts.ApplyOptions(opts)
	return o.Delete().
		NamespaceIfScoped(o.GetNamespace(), o.isNamespaced()).
		Resource(o.resource()).
		Name(o.GetName()).
		Body(deleteOpts.AsDeleteOptions()).
		Do(ctx).
		Error()
}

// DeleteAllOf implements client.Client
func (uc *unstructuredClient) DeleteAllOf(ctx context.Context, obj Object, opts ...DeleteAllOfOption) error {
	_, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return fmt.Errorf("unstructured client did not understand object: %T", obj)
	}

	o, err := uc.cache.getObjMeta(obj)
	if err != nil {
		return err
	}

	deleteAllOfOpts := DeleteAllOfOptions{}
	deleteAllOfOpts.ApplyOptions(opts)
	return o.Delete().
		NamespaceIfScoped(deleteAllOfOpts.ListOptions.Namespace, o.isNamespaced()).
		Resource(o.resource()).
		VersionedParams(deleteAllOfOpts.AsListOptions(), uc.paramCodec).
		Body(deleteAllOfOpts.AsDeleteOptions()).
		Do(ctx).
		Error()
}

// Patch implements client.Client
<<<<<<< HEAD
func (uc *unstructuredClient) Patch(_ context.Context, obj runtime.Object, patch Patch, opts ...PatchOption) error {
	u, ok := obj.(*unstructured.Unstructured)
>>>>>>> 79bfea2d (update vendor)
=======
func (uc *unstructuredClient) Patch(ctx context.Context, obj Object, patch Patch, opts ...PatchOption) error {
	_, ok := obj.(*unstructured.Unstructured)
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
>>>>>>> 03397665 (update api)
	if !ok {
		return fmt.Errorf("unstructured client did not understand object: %T", obj)
	}

	o, err := uc.cache.getObjMeta(obj)
	if err != nil {
		return err
	}

	data, err := patch.Data(obj)
	if err != nil {
		return err
	}

	patchOpts := &PatchOptions{}
	return o.Patch(patch.Type()).
		NamespaceIfScoped(o.GetNamespace(), o.isNamespaced()).
		Resource(o.resource()).
		Name(o.GetName()).
		VersionedParams(patchOpts.ApplyOptions(opts).AsPatchOptions(), uc.paramCodec).
		Body(data).
		Do(ctx).
		Into(obj)
}

// Get implements client.Client
func (uc *unstructuredClient) Get(ctx context.Context, key ObjectKey, obj Object) error {
	u, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return fmt.Errorf("unstructured client did not understand object: %T", obj)
	}

	gvk := u.GroupVersionKind()

	r, err := uc.cache.getResource(obj)
	if err != nil {
		return err
	}

	result := r.Get().
		NamespaceIfScoped(key.Namespace, r.isNamespaced()).
		Resource(r.resource()).
		Name(key.Name).
		Do(ctx).
		Into(obj)

	u.SetGroupVersionKind(gvk)

	return result
}

// List implements client.Client
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
func (uc *unstructuredClient) List(ctx context.Context, obj ObjectList, opts ...ListOption) error {
=======
func (uc *unstructuredClient) List(_ context.Context, obj runtime.Object, opts ...ListOption) error {
>>>>>>> 79bfea2d (update vendor)
=======
func (uc *unstructuredClient) List(ctx context.Context, obj ObjectList, opts ...ListOption) error {
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
func (uc *unstructuredClient) List(ctx context.Context, obj ObjectList, opts ...ListOption) error {
>>>>>>> 03397665 (update api)
	u, ok := obj.(*unstructured.UnstructuredList)
	if !ok {
		return fmt.Errorf("unstructured client did not understand object: %T", obj)
	}

	gvk := u.GroupVersionKind()
	if strings.HasSuffix(gvk.Kind, "List") {
		gvk.Kind = gvk.Kind[:len(gvk.Kind)-4]
	}

	listOpts := ListOptions{}
	listOpts.ApplyOptions(opts)

	r, err := uc.cache.getResource(obj)
	if err != nil {
		return err
	}

	return r.Get().
		NamespaceIfScoped(listOpts.Namespace, r.isNamespaced()).
		Resource(r.resource()).
		VersionedParams(listOpts.AsListOptions(), uc.paramCodec).
		Do(ctx).
		Into(obj)
}

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
func (uc *unstructuredClient) UpdateStatus(ctx context.Context, obj Object, opts ...UpdateOption) error {
	_, ok := obj.(*unstructured.Unstructured)
=======
func (uc *unstructuredClient) UpdateStatus(_ context.Context, obj runtime.Object, opts ...UpdateOption) error {
	u, ok := obj.(*unstructured.Unstructured)
>>>>>>> 79bfea2d (update vendor)
=======
func (uc *unstructuredClient) UpdateStatus(ctx context.Context, obj Object, opts ...UpdateOption) error {
	_, ok := obj.(*unstructured.Unstructured)
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
func (uc *unstructuredClient) UpdateStatus(ctx context.Context, obj Object, opts ...UpdateOption) error {
	_, ok := obj.(*unstructured.Unstructured)
>>>>>>> 03397665 (update api)
	if !ok {
		return fmt.Errorf("unstructured client did not understand object: %T", obj)
	}

	o, err := uc.cache.getObjMeta(obj)
	if err != nil {
		return err
	}

	return o.Put().
		NamespaceIfScoped(o.GetNamespace(), o.isNamespaced()).
		Resource(o.resource()).
		Name(o.GetName()).
		SubResource("status").
		Body(obj).
		VersionedParams((&UpdateOptions{}).ApplyOptions(opts).AsUpdateOptions(), uc.paramCodec).
		Do(ctx).
		Into(obj)
}

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
func (uc *unstructuredClient) PatchStatus(ctx context.Context, obj Object, patch Patch, opts ...PatchOption) error {
=======
func (uc *unstructuredClient) PatchStatus(_ context.Context, obj runtime.Object, patch Patch, opts ...PatchOption) error {
>>>>>>> 79bfea2d (update vendor)
=======
func (uc *unstructuredClient) PatchStatus(ctx context.Context, obj Object, patch Patch, opts ...PatchOption) error {
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
func (uc *unstructuredClient) PatchStatus(ctx context.Context, obj Object, patch Patch, opts ...PatchOption) error {
>>>>>>> 03397665 (update api)
	u, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return fmt.Errorf("unstructured client did not understand object: %T", obj)
	}

	gvk := u.GroupVersionKind()

	o, err := uc.cache.getObjMeta(obj)
	if err != nil {
		return err
	}

	data, err := patch.Data(obj)
	if err != nil {
		return err
	}

	patchOpts := &PatchOptions{}
	result := o.Patch(patch.Type()).
		NamespaceIfScoped(o.GetNamespace(), o.isNamespaced()).
		Resource(o.resource()).
		Name(o.GetName()).
		SubResource("status").
		Body(data).
		VersionedParams(patchOpts.ApplyOptions(opts).AsPatchOptions(), uc.paramCodec).
		Do(ctx).
		Into(u)

	u.SetGroupVersionKind(gvk)
	return result
}
