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

package builder

import (
	"fmt"
	"strings"

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"sigs.k8s.io/controller-runtime/pkg/client"
<<<<<<< HEAD
=======
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
>>>>>>> 79bfea2d (update vendor)
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Supporting mocking out functions for testing
var newController = controller.New
var getGvk = apiutil.GVKForObject

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
// project represents other forms that the we can use to
// send/receive a given resource (metadata-only, unstructured, etc)
type objectProjection int

const (
	// projectAsNormal doesn't change the object from the form given
	projectAsNormal objectProjection = iota
	// projectAsMetadata turns this into an metadata-only watch
	projectAsMetadata
)

<<<<<<< HEAD
// Builder builds a Controller.
type Builder struct {
	forInput         ForInput
	ownsInput        []OwnsInput
	watchesInput     []WatchesInput
	mgr              manager.Manager
	globalPredicates []predicate.Predicate
	ctrl             controller.Controller
	ctrlOptions      controller.Options
	name             string
=======
// Builder builds a Controller.
type Builder struct {
	apiType        runtime.Object
	mgr            manager.Manager
	predicates     []predicate.Predicate
	managedObjects []runtime.Object
	watchRequest   []watchRequest
	config         *rest.Config
	ctrl           controller.Controller
	ctrlOptions    controller.Options
	name           string
>>>>>>> 79bfea2d (update vendor)
=======
// Builder builds a Controller.
type Builder struct {
	forInput         ForInput
	ownsInput        []OwnsInput
	watchesInput     []WatchesInput
	mgr              manager.Manager
	globalPredicates []predicate.Predicate
	ctrl             controller.Controller
	ctrlOptions      controller.Options
	name             string
>>>>>>> e879a141 (alibabacloud machine-api provider)
}

// ControllerManagedBy returns a new controller builder that will be started by the provided Manager
func ControllerManagedBy(m manager.Manager) *Builder {
	return &Builder{mgr: m}
}

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
// ForInput represents the information set by For method.
type ForInput struct {
	object           client.Object
	predicates       []predicate.Predicate
	objectProjection objectProjection
	err              error
<<<<<<< HEAD
=======
// ForType defines the type of Object being *reconciled*, and configures the ControllerManagedBy to respond to create / delete /
// update events by *reconciling the object*.
// This is the equivalent of calling
// Watches(&source.Kind{Type: apiType}, &handler.EnqueueRequestForObject{})
//
// Deprecated: Use For
func (blder *Builder) ForType(apiType runtime.Object) *Builder {
	return blder.For(apiType)
>>>>>>> 79bfea2d (update vendor)
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
}

// For defines the type of Object being *reconciled*, and configures the ControllerManagedBy to respond to create / delete /
// update events by *reconciling the object*.
// This is the equivalent of calling
// Watches(&source.Kind{Type: apiType}, &handler.EnqueueRequestForObject{})
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
func (blder *Builder) For(object client.Object, opts ...ForOption) *Builder {
	if blder.forInput.object != nil {
		blder.forInput.err = fmt.Errorf("For(...) should only be called once, could not assign multiple objects for reconciliation")
		return blder
	}
	input := ForInput{object: object}
	for _, opt := range opts {
		opt.ApplyToFor(&input)
	}

	blder.forInput = input
<<<<<<< HEAD
	return blder
}

// OwnsInput represents the information set by Owns method.
type OwnsInput struct {
	object           client.Object
	predicates       []predicate.Predicate
	objectProjection objectProjection
}

// Owns defines types of Objects being *generated* by the ControllerManagedBy, and configures the ControllerManagedBy to respond to
// create / delete / update events by *reconciling the owner object*.  This is the equivalent of calling
// Watches(&source.Kind{Type: <ForType-forInput>}, &handler.EnqueueRequestForOwner{OwnerType: apiType, IsController: true})
func (blder *Builder) Owns(object client.Object, opts ...OwnsOption) *Builder {
	input := OwnsInput{object: object}
	for _, opt := range opts {
		opt.ApplyToOwns(&input)
	}

	blder.ownsInput = append(blder.ownsInput, input)
	return blder
}

// WatchesInput represents the information set by Watches method.
type WatchesInput struct {
	src              source.Source
	eventhandler     handler.EventHandler
	predicates       []predicate.Predicate
	objectProjection objectProjection
=======
func (blder *Builder) For(apiType runtime.Object) *Builder {
	blder.apiType = apiType
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
	return blder
}

// OwnsInput represents the information set by Owns method.
type OwnsInput struct {
	object           client.Object
	predicates       []predicate.Predicate
	objectProjection objectProjection
}

// Owns defines types of Objects being *generated* by the ControllerManagedBy, and configures the ControllerManagedBy to respond to
// create / delete / update events by *reconciling the owner object*.  This is the equivalent of calling
// Watches(&source.Kind{Type: <ForType-forInput>}, &handler.EnqueueRequestForOwner{OwnerType: apiType, IsController: true})
func (blder *Builder) Owns(object client.Object, opts ...OwnsOption) *Builder {
	input := OwnsInput{object: object}
	for _, opt := range opts {
		opt.ApplyToOwns(&input)
	}

	blder.ownsInput = append(blder.ownsInput, input)
	return blder
}

<<<<<<< HEAD
type watchRequest struct {
	src          source.Source
	eventhandler handler.EventHandler
>>>>>>> 79bfea2d (update vendor)
=======
// WatchesInput represents the information set by Watches method.
type WatchesInput struct {
	src              source.Source
	eventhandler     handler.EventHandler
	predicates       []predicate.Predicate
	objectProjection objectProjection
>>>>>>> e879a141 (alibabacloud machine-api provider)
}

// Watches exposes the lower-level ControllerManagedBy Watches functions through the builder.  Consider using
// Owns or For instead of Watches directly.
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
// Specified predicates are registered only for given source.
func (blder *Builder) Watches(src source.Source, eventhandler handler.EventHandler, opts ...WatchesOption) *Builder {
	input := WatchesInput{src: src, eventhandler: eventhandler}
	for _, opt := range opts {
		opt.ApplyToWatches(&input)
	}
<<<<<<< HEAD

	blder.watchesInput = append(blder.watchesInput, input)
=======
func (blder *Builder) Watches(src source.Source, eventhandler handler.EventHandler) *Builder {
	blder.watchRequest = append(blder.watchRequest, watchRequest{src: src, eventhandler: eventhandler})
	return blder
}

// WithConfig sets the Config to use for configuring clients.  Defaults to the in-cluster config or to ~/.kube/config.
//
// Deprecated: Use ControllerManagedBy(Manager) and this isn't needed.
func (blder *Builder) WithConfig(config *rest.Config) *Builder {
	blder.config = config
>>>>>>> 79bfea2d (update vendor)
=======

	blder.watchesInput = append(blder.watchesInput, input)
>>>>>>> e879a141 (alibabacloud machine-api provider)
	return blder
}

// WithEventFilter sets the event filters, to filter which create/update/delete/generic events eventually
// trigger reconciliations.  For example, filtering on whether the resource version has changed.
<<<<<<< HEAD
<<<<<<< HEAD
// Given predicate is added for all watched objects.
// Defaults to the empty list.
func (blder *Builder) WithEventFilter(p predicate.Predicate) *Builder {
	blder.globalPredicates = append(blder.globalPredicates, p)
=======
// Defaults to the empty list.
func (blder *Builder) WithEventFilter(p predicate.Predicate) *Builder {
	blder.predicates = append(blder.predicates, p)
>>>>>>> 79bfea2d (update vendor)
=======
// Given predicate is added for all watched objects.
// Defaults to the empty list.
func (blder *Builder) WithEventFilter(p predicate.Predicate) *Builder {
	blder.globalPredicates = append(blder.globalPredicates, p)
>>>>>>> e879a141 (alibabacloud machine-api provider)
	return blder
}

// WithOptions overrides the controller options use in doController. Defaults to empty.
func (blder *Builder) WithOptions(options controller.Options) *Builder {
	blder.ctrlOptions = options
	return blder
}

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
// WithLogger overrides the controller options's logger used.
func (blder *Builder) WithLogger(log logr.Logger) *Builder {
	blder.ctrlOptions.Log = log
	return blder
}

<<<<<<< HEAD
=======
>>>>>>> 79bfea2d (update vendor)
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
// Named sets the name of the controller to the given name.  The name shows up
// in metrics, among other things, and thus should be a prometheus compatible name
// (underscores and alphanumeric characters only).
//
// By default, controllers are named using the lowercase version of their kind.
func (blder *Builder) Named(name string) *Builder {
	blder.name = name
	return blder
}

<<<<<<< HEAD
<<<<<<< HEAD
// Complete builds the Application Controller.
=======
// Complete builds the Application ControllerManagedBy.
>>>>>>> 79bfea2d (update vendor)
=======
// Complete builds the Application Controller.
>>>>>>> e879a141 (alibabacloud machine-api provider)
func (blder *Builder) Complete(r reconcile.Reconciler) error {
	_, err := blder.Build(r)
	return err
}

<<<<<<< HEAD
<<<<<<< HEAD
// Build builds the Application Controller and returns the Controller it created.
=======
// Build builds the Application ControllerManagedBy and returns the Controller it created.
>>>>>>> 79bfea2d (update vendor)
=======
// Build builds the Application Controller and returns the Controller it created.
>>>>>>> e879a141 (alibabacloud machine-api provider)
func (blder *Builder) Build(r reconcile.Reconciler) (controller.Controller, error) {
	if r == nil {
		return nil, fmt.Errorf("must provide a non-nil Reconciler")
	}
	if blder.mgr == nil {
		return nil, fmt.Errorf("must provide a non-nil Manager")
	}
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
	if blder.forInput.err != nil {
		return nil, blder.forInput.err
	}
	// Checking the reconcile type exist or not
	if blder.forInput.object == nil {
		return nil, fmt.Errorf("must provide an object for reconciliation")
	}
<<<<<<< HEAD
=======

	// Set the Config
	blder.loadRestConfig()
>>>>>>> 79bfea2d (update vendor)
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)

	// Set the ControllerManagedBy
	if err := blder.doController(r); err != nil {
		return nil, err
	}

	// Set the Watch
	if err := blder.doWatch(); err != nil {
		return nil, err
	}

	return blder.ctrl, nil
}

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
func (blder *Builder) project(obj client.Object, proj objectProjection) (client.Object, error) {
	switch proj {
	case projectAsNormal:
		return obj, nil
	case projectAsMetadata:
		metaObj := &metav1.PartialObjectMetadata{}
		gvk, err := getGvk(obj, blder.mgr.GetScheme())
		if err != nil {
			return nil, fmt.Errorf("unable to determine GVK of %T for a metadata-only watch: %w", obj, err)
		}
		metaObj.SetGroupVersionKind(gvk)
		return metaObj, nil
	default:
		panic(fmt.Sprintf("unexpected projection type %v on type %T, should not be possible since this is an internal field", proj, obj))
	}
}

<<<<<<< HEAD
func (blder *Builder) doWatch() error {
	// Reconcile type
	typeForSrc, err := blder.project(blder.forInput.object, blder.forInput.objectProjection)
	if err != nil {
		return err
	}
	src := &source.Kind{Type: typeForSrc}
	hdler := &handler.EnqueueRequestForObject{}
	allPredicates := append(blder.globalPredicates, blder.forInput.predicates...)
	if err := blder.ctrl.Watch(src, hdler, allPredicates...); err != nil {
=======
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
func (blder *Builder) doWatch() error {
	// Reconcile type
	typeForSrc, err := blder.project(blder.forInput.object, blder.forInput.objectProjection)
	if err != nil {
>>>>>>> 79bfea2d (update vendor)
		return err
	}
	src := &source.Kind{Type: typeForSrc}
	hdler := &handler.EnqueueRequestForObject{}
	allPredicates := append(blder.globalPredicates, blder.forInput.predicates...)
	if err := blder.ctrl.Watch(src, hdler, allPredicates...); err != nil {
		return err
	}

	// Watches the managed types
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
	for _, own := range blder.ownsInput {
		typeForSrc, err := blder.project(own.object, own.objectProjection)
		if err != nil {
			return err
		}
		src := &source.Kind{Type: typeForSrc}
<<<<<<< HEAD
		hdler := &handler.EnqueueRequestForOwner{
			OwnerType:    blder.forInput.object,
			IsController: true,
		}
		allPredicates := append([]predicate.Predicate(nil), blder.globalPredicates...)
		allPredicates = append(allPredicates, own.predicates...)
		if err := blder.ctrl.Watch(src, hdler, allPredicates...); err != nil {
=======
	for _, obj := range blder.managedObjects {
		src := &source.Kind{Type: obj}
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
		hdler := &handler.EnqueueRequestForOwner{
			OwnerType:    blder.forInput.object,
			IsController: true,
		}
<<<<<<< HEAD
		if err := blder.ctrl.Watch(src, hdler, blder.predicates...); err != nil {
>>>>>>> 79bfea2d (update vendor)
=======
		allPredicates := append([]predicate.Predicate(nil), blder.globalPredicates...)
		allPredicates = append(allPredicates, own.predicates...)
		if err := blder.ctrl.Watch(src, hdler, allPredicates...); err != nil {
>>>>>>> e879a141 (alibabacloud machine-api provider)
			return err
		}
	}

	// Do the watch requests
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
	for _, w := range blder.watchesInput {
		allPredicates := append([]predicate.Predicate(nil), blder.globalPredicates...)
		allPredicates = append(allPredicates, w.predicates...)

		// If the source of this watch is of type *source.Kind, project it.
		if srckind, ok := w.src.(*source.Kind); ok {
			typeForSrc, err := blder.project(srckind.Type, w.objectProjection)
			if err != nil {
				return err
			}
			srckind.Type = typeForSrc
<<<<<<< HEAD
		}

		if err := blder.ctrl.Watch(w.src, w.eventhandler, allPredicates...); err != nil {
			return err
		}
=======
	for _, w := range blder.watchRequest {
		if err := blder.ctrl.Watch(w.src, w.eventhandler, blder.predicates...); err != nil {
			return err
		}

>>>>>>> 79bfea2d (update vendor)
=======
		}

		if err := blder.ctrl.Watch(w.src, w.eventhandler, allPredicates...); err != nil {
			return err
		}
>>>>>>> e879a141 (alibabacloud machine-api provider)
	}
	return nil
}

<<<<<<< HEAD
<<<<<<< HEAD
func (blder *Builder) getControllerName(gvk schema.GroupVersionKind) string {
	if blder.name != "" {
		return blder.name
	}
	return strings.ToLower(gvk.Kind)
}

func (blder *Builder) doController(r reconcile.Reconciler) error {
	globalOpts := blder.mgr.GetControllerOptions()

	ctrlOptions := blder.ctrlOptions
	if ctrlOptions.Reconciler == nil {
		ctrlOptions.Reconciler = r
	}

	// Retrieve the GVK from the object we're reconciling
	// to prepopulate logger information, and to optionally generate a default name.
	gvk, err := getGvk(blder.forInput.object, blder.mgr.GetScheme())
	if err != nil {
		return err
	}

	// Setup concurrency.
	if ctrlOptions.MaxConcurrentReconciles == 0 {
		groupKind := gvk.GroupKind().String()

		if concurrency, ok := globalOpts.GroupKindConcurrency[groupKind]; ok && concurrency > 0 {
			ctrlOptions.MaxConcurrentReconciles = concurrency
		}
	}

	// Setup cache sync timeout.
	if ctrlOptions.CacheSyncTimeout == 0 && globalOpts.CacheSyncTimeout != nil {
		ctrlOptions.CacheSyncTimeout = *globalOpts.CacheSyncTimeout
	}

	// Setup the logger.
	if ctrlOptions.Log == nil {
		ctrlOptions.Log = blder.mgr.GetLogger()
	}
	ctrlOptions.Log = ctrlOptions.Log.WithValues("reconciler group", gvk.Group, "reconciler kind", gvk.Kind)

	// Build the controller and return.
	blder.ctrl, err = newController(blder.getControllerName(gvk), blder.mgr, ctrlOptions)
=======
func (blder *Builder) loadRestConfig() {
	if blder.config == nil {
		blder.config = blder.mgr.GetConfig()
	}
}

func (blder *Builder) getControllerName() (string, error) {
=======
func (blder *Builder) getControllerName(gvk schema.GroupVersionKind) string {
>>>>>>> e879a141 (alibabacloud machine-api provider)
	if blder.name != "" {
		return blder.name
	}
	return strings.ToLower(gvk.Kind)
}

func (blder *Builder) doController(r reconcile.Reconciler) error {
	globalOpts := blder.mgr.GetControllerOptions()

	ctrlOptions := blder.ctrlOptions
	if ctrlOptions.Reconciler == nil {
		ctrlOptions.Reconciler = r
	}

	// Retrieve the GVK from the object we're reconciling
	// to prepopulate logger information, and to optionally generate a default name.
	gvk, err := getGvk(blder.forInput.object, blder.mgr.GetScheme())
	if err != nil {
		return err
	}
<<<<<<< HEAD
	ctrlOptions := blder.ctrlOptions
	ctrlOptions.Reconciler = r
	blder.ctrl, err = newController(name, blder.mgr, ctrlOptions)
>>>>>>> 79bfea2d (update vendor)
=======

	// Setup concurrency.
	if ctrlOptions.MaxConcurrentReconciles == 0 {
		groupKind := gvk.GroupKind().String()

		if concurrency, ok := globalOpts.GroupKindConcurrency[groupKind]; ok && concurrency > 0 {
			ctrlOptions.MaxConcurrentReconciles = concurrency
		}
	}

	// Setup cache sync timeout.
	if ctrlOptions.CacheSyncTimeout == 0 && globalOpts.CacheSyncTimeout != nil {
		ctrlOptions.CacheSyncTimeout = *globalOpts.CacheSyncTimeout
	}

	// Setup the logger.
	if ctrlOptions.Log == nil {
		ctrlOptions.Log = blder.mgr.GetLogger()
	}
	ctrlOptions.Log = ctrlOptions.Log.WithValues("reconciler group", gvk.Group, "reconciler kind", gvk.Kind)

	// Build the controller and return.
	blder.ctrl, err = newController(blder.getControllerName(gvk), blder.mgr, ctrlOptions)
>>>>>>> e879a141 (alibabacloud machine-api provider)
	return err
}
