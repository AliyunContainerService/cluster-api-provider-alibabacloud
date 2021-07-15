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

package envtest

import (
	"bufio"
	"bytes"
<<<<<<< HEAD
<<<<<<< HEAD
	"context"
=======
>>>>>>> 79bfea2d (update vendor)
=======
	"context"
>>>>>>> 737a8f1c (add more test case)
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

<<<<<<< HEAD
<<<<<<< HEAD
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
=======
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
>>>>>>> 79bfea2d (update vendor)
=======
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
>>>>>>> 737a8f1c (add more test case)
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/rest"
<<<<<<< HEAD
<<<<<<< HEAD
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
=======
>>>>>>> 79bfea2d (update vendor)
=======
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
>>>>>>> 737a8f1c (add more test case)
	"sigs.k8s.io/yaml"
)

// CRDInstallOptions are the options for installing CRDs
type CRDInstallOptions struct {
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 737a8f1c (add more test case)
	// Paths is a list of paths to the directories or files containing CRDs
	Paths []string

	// CRDs is a list of CRDs to install
	CRDs []client.Object
<<<<<<< HEAD
=======
	// Paths is the path to the directory containing CRDs
	Paths []string

	// CRDs is a list of CRDs to install
	CRDs []*apiextensionsv1beta1.CustomResourceDefinition
>>>>>>> 79bfea2d (update vendor)
=======
>>>>>>> 737a8f1c (add more test case)

	// ErrorIfPathMissing will cause an error if a Path does not exist
	ErrorIfPathMissing bool

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 737a8f1c (add more test case)
	// MaxTime is the max time to wait
	MaxTime time.Duration

	// PollInterval is the interval to check
	PollInterval time.Duration

	// CleanUpAfterUse will cause the CRDs listed for installation to be
	// uninstalled when terminating the test environment.
	// Defaults to false.
	CleanUpAfterUse bool
<<<<<<< HEAD
=======
	// maxTime is the max time to wait
	maxTime time.Duration

	// pollInterval is the interval to check
	pollInterval time.Duration
>>>>>>> 79bfea2d (update vendor)
=======
>>>>>>> 737a8f1c (add more test case)
}

const defaultPollInterval = 100 * time.Millisecond
const defaultMaxWait = 10 * time.Second

// InstallCRDs installs a collection of CRDs into a cluster by reading the crd yaml files from a directory
<<<<<<< HEAD
<<<<<<< HEAD
func InstallCRDs(config *rest.Config, options CRDInstallOptions) ([]client.Object, error) {
=======
func InstallCRDs(config *rest.Config, options CRDInstallOptions) ([]*apiextensionsv1beta1.CustomResourceDefinition, error) {
>>>>>>> 79bfea2d (update vendor)
=======
func InstallCRDs(config *rest.Config, options CRDInstallOptions) ([]client.Object, error) {
>>>>>>> 737a8f1c (add more test case)
	defaultCRDOptions(&options)

	// Read the CRD yamls into options.CRDs
	if err := readCRDFiles(&options); err != nil {
		return nil, err
	}

	// Create the CRDs in the apiserver
	if err := CreateCRDs(config, options.CRDs); err != nil {
		return options.CRDs, err
	}

	// Wait for the CRDs to appear as Resources in the apiserver
	if err := WaitForCRDs(config, options.CRDs, options); err != nil {
		return options.CRDs, err
	}

	return options.CRDs, nil
}

// readCRDFiles reads the directories of CRDs in options.Paths and adds the CRD structs to options.CRDs
func readCRDFiles(options *CRDInstallOptions) error {
	if len(options.Paths) > 0 {
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 737a8f1c (add more test case)
		crdList, err := renderCRDs(options)
		if err != nil {
			return err
		}

		options.CRDs = append(options.CRDs, crdList...)
<<<<<<< HEAD
=======
		for _, path := range options.Paths {
			if _, err := os.Stat(path); !options.ErrorIfPathMissing && os.IsNotExist(err) {
				continue
			}
			new, err := readCRDs(path)
			if err != nil {
				return err
			}
			options.CRDs = append(options.CRDs, new...)
		}
>>>>>>> 79bfea2d (update vendor)
=======
>>>>>>> 737a8f1c (add more test case)
	}
	return nil
}

// defaultCRDOptions sets the default values for CRDs
func defaultCRDOptions(o *CRDInstallOptions) {
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 737a8f1c (add more test case)
	if o.MaxTime == 0 {
		o.MaxTime = defaultMaxWait
	}
	if o.PollInterval == 0 {
		o.PollInterval = defaultPollInterval
<<<<<<< HEAD
=======
	if o.maxTime == 0 {
		o.maxTime = defaultMaxWait
	}
	if o.pollInterval == 0 {
		o.pollInterval = defaultPollInterval
>>>>>>> 79bfea2d (update vendor)
=======
>>>>>>> 737a8f1c (add more test case)
	}
}

// WaitForCRDs waits for the CRDs to appear in discovery
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 737a8f1c (add more test case)
func WaitForCRDs(config *rest.Config, crds []client.Object, options CRDInstallOptions) error {
	// Add each CRD to a map of GroupVersion to Resource
	waitingFor := map[schema.GroupVersion]*sets.String{}
	for _, crd := range runtimeCRDListToUnstructured(crds) {
		gvs := []schema.GroupVersion{}
		crdGroup, _, err := unstructured.NestedString(crd.Object, "spec", "group")
		if err != nil {
			return err
		}
		crdPlural, _, err := unstructured.NestedString(crd.Object, "spec", "names", "plural")
		if err != nil {
			return err
		}
		crdVersion, _, err := unstructured.NestedString(crd.Object, "spec", "version")
		if err != nil {
			return err
		}
		versions, found, err := unstructured.NestedSlice(crd.Object, "spec", "versions")
		if err != nil {
			return err
		}

		// gvs should be added here only if single version is found. If multiple version is found we will add those version
		// based on the version is served or not.
		if crdVersion != "" && !found {
			gvs = append(gvs, schema.GroupVersion{Group: crdGroup, Version: crdVersion})
		}

		for _, version := range versions {
			versionMap, ok := version.(map[string]interface{})
			if !ok {
				continue
			}
			served, _, err := unstructured.NestedBool(versionMap, "served")
			if err != nil {
				return err
			}
			if served {
				versionName, _, err := unstructured.NestedString(versionMap, "name")
				if err != nil {
					return err
				}
				gvs = append(gvs, schema.GroupVersion{Group: crdGroup, Version: versionName})
			}
		}

<<<<<<< HEAD
=======
func WaitForCRDs(config *rest.Config, crds []*apiextensionsv1beta1.CustomResourceDefinition, options CRDInstallOptions) error {
	// Add each CRD to a map of GroupVersion to Resource
	waitingFor := map[schema.GroupVersion]*sets.String{}
	for _, crd := range crds {
		gvs := []schema.GroupVersion{}
		if crd.Spec.Version != "" {
			gvs = append(gvs, schema.GroupVersion{Group: crd.Spec.Group, Version: crd.Spec.Version})
		}
		for _, ver := range crd.Spec.Versions {
			if ver.Served {
				gvs = append(gvs, schema.GroupVersion{Group: crd.Spec.Group, Version: ver.Name})
			}
		}
>>>>>>> 79bfea2d (update vendor)
=======
>>>>>>> 737a8f1c (add more test case)
		for _, gv := range gvs {
			log.V(1).Info("adding API in waitlist", "GV", gv)
			if _, found := waitingFor[gv]; !found {
				// Initialize the set
				waitingFor[gv] = &sets.String{}
			}
			// Add the Resource
<<<<<<< HEAD
<<<<<<< HEAD
			waitingFor[gv].Insert(crdPlural)
=======
			waitingFor[gv].Insert(crd.Spec.Names.Plural)
>>>>>>> 79bfea2d (update vendor)
=======
			waitingFor[gv].Insert(crdPlural)
>>>>>>> 737a8f1c (add more test case)
		}
	}

	// Poll until all resources are found in discovery
	p := &poller{config: config, waitingFor: waitingFor}
<<<<<<< HEAD
<<<<<<< HEAD
	return wait.PollImmediate(options.PollInterval, options.MaxTime, p.poll)
=======
	return wait.PollImmediate(options.pollInterval, options.maxTime, p.poll)
>>>>>>> 79bfea2d (update vendor)
=======
	return wait.PollImmediate(options.PollInterval, options.MaxTime, p.poll)
>>>>>>> 737a8f1c (add more test case)
}

// poller checks if all the resources have been found in discovery, and returns false if not
type poller struct {
	// config is used to get discovery
	config *rest.Config

	// waitingFor is the map of resources keyed by group version that have not yet been found in discovery
	waitingFor map[schema.GroupVersion]*sets.String
}

// poll checks if all the resources have been found in discovery, and returns false if not
func (p *poller) poll() (done bool, err error) {
	// Create a new clientset to avoid any client caching of discovery
	cs, err := clientset.NewForConfig(p.config)
	if err != nil {
		return false, err
	}

	allFound := true
	for gv, resources := range p.waitingFor {
		// All resources found, do nothing
		if resources.Len() == 0 {
			delete(p.waitingFor, gv)
			continue
		}

		// Get the Resources for this GroupVersion
		// TODO: Maybe the controller-runtime client should be able to do this...
		resourceList, err := cs.Discovery().ServerResourcesForGroupVersion(gv.Group + "/" + gv.Version)
		if err != nil {
			return false, nil
		}

		// Remove each found resource from the resources set that we are waiting for
		for _, resource := range resourceList.APIResources {
			resources.Delete(resource.Name)
		}

		// Still waiting on some resources in this group version
		if resources.Len() != 0 {
			allFound = false
		}
	}
	return allFound, nil
}

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 737a8f1c (add more test case)
// UninstallCRDs uninstalls a collection of CRDs by reading the crd yaml files from a directory
func UninstallCRDs(config *rest.Config, options CRDInstallOptions) error {

	// Read the CRD yamls into options.CRDs
	if err := readCRDFiles(&options); err != nil {
		return err
	}

	// Delete the CRDs from the apiserver
	cs, err := client.New(config, client.Options{})
	if err != nil {
		return err
	}

	// Uninstall each CRD
	for _, crd := range runtimeCRDListToUnstructured(options.CRDs) {
		log.V(1).Info("uninstalling CRD", "crd", crd.GetName())
		if err := cs.Delete(context.TODO(), crd); err != nil {
			// If CRD is not found, we can consider success
			if !apierrors.IsNotFound(err) {
				return err
			}
		}
	}

	return nil
}

// CreateCRDs creates the CRDs
func CreateCRDs(config *rest.Config, crds []client.Object) error {
	cs, err := client.New(config, client.Options{})
<<<<<<< HEAD
=======
// CreateCRDs creates the CRDs
func CreateCRDs(config *rest.Config, crds []*apiextensionsv1beta1.CustomResourceDefinition) error {
	cs, err := clientset.NewForConfig(config)
>>>>>>> 79bfea2d (update vendor)
=======
>>>>>>> 737a8f1c (add more test case)
	if err != nil {
		return err
	}

	// Create each CRD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 737a8f1c (add more test case)
	for _, crd := range runtimeCRDListToUnstructured(crds) {
		log.V(1).Info("installing CRD", "crd", crd.GetName())
		existingCrd := crd.DeepCopy()
		err := cs.Get(context.TODO(), client.ObjectKey{Name: crd.GetName()}, existingCrd)
		switch {
		case apierrors.IsNotFound(err):
			if err := cs.Create(context.TODO(), crd); err != nil {
				return err
			}
		case err != nil:
			return err
		default:
			log.V(1).Info("CRD already exists, updating", "crd", crd.GetName())
			if err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
				if err := cs.Get(context.TODO(), client.ObjectKey{Name: crd.GetName()}, existingCrd); err != nil {
					return err
				}
				crd.SetResourceVersion(existingCrd.GetResourceVersion())
				return cs.Update(context.TODO(), crd)
			}); err != nil {
				return err
			}
<<<<<<< HEAD
=======
	for _, crd := range crds {
		log.V(1).Info("installing CRD", "crd", crd.Name)
		if _, err := cs.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd); err != nil {
			return err
>>>>>>> 79bfea2d (update vendor)
=======
>>>>>>> 737a8f1c (add more test case)
		}
	}
	return nil
}

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 737a8f1c (add more test case)
// renderCRDs iterate through options.Paths and extract all CRD files.
func renderCRDs(options *CRDInstallOptions) ([]client.Object, error) {
	var (
		err   error
		info  os.FileInfo
		files []os.FileInfo
	)

	type GVKN struct {
		GVK  schema.GroupVersionKind
		Name string
	}

	crds := map[GVKN]*unstructured.Unstructured{}

	for _, path := range options.Paths {
		var filePath = path

		// Return the error if ErrorIfPathMissing exists
		if info, err = os.Stat(path); os.IsNotExist(err) {
			if options.ErrorIfPathMissing {
				return nil, err
			}
			continue
		}

		if !info.IsDir() {
			filePath, files = filepath.Dir(path), []os.FileInfo{info}
		} else {
			if files, err = ioutil.ReadDir(path); err != nil {
				return nil, err
			}
		}

		log.V(1).Info("reading CRDs from path", "path", path)
		crdList, err := readCRDs(filePath, files)
		if err != nil {
			return nil, err
		}

		for i, crd := range crdList {
			gvkn := GVKN{GVK: crd.GroupVersionKind(), Name: crd.GetName()}
			if _, found := crds[gvkn]; found {
				// Currently, we only print a log when there are duplicates. We may want to error out if that makes more sense.
				log.Info("there are more than one CRD definitions with the same <Group, Version, Kind, Name>", "GVKN", gvkn)
			}
			// We always use the CRD definition that we found last.
			crds[gvkn] = crdList[i]
		}
	}

	// Converting map to a list to return
	var res []client.Object
	for _, obj := range crds {
		res = append(res, obj)
	}
	return res, nil
}

// readCRDs reads the CRDs from files and Unmarshals them into structs
func readCRDs(basePath string, files []os.FileInfo) ([]*unstructured.Unstructured, error) {
	var crds []*unstructured.Unstructured
<<<<<<< HEAD
=======
// readCRDs reads the CRDs from files and Unmarshals them into structs
func readCRDs(path string) ([]*apiextensionsv1beta1.CustomResourceDefinition, error) {
	// Get the CRD files
	var files []os.FileInfo
	var err error
	log.V(1).Info("reading CRDs from path", "path", path)
	if files, err = ioutil.ReadDir(path); err != nil {
		return nil, err
	}
>>>>>>> 79bfea2d (update vendor)
=======
>>>>>>> 737a8f1c (add more test case)

	// White list the file extensions that may contain CRDs
	crdExts := sets.NewString(".json", ".yaml", ".yml")

<<<<<<< HEAD
<<<<<<< HEAD
	for _, file := range files {
		// Only parse allowlisted file types
=======
	var crds []*apiextensionsv1beta1.CustomResourceDefinition
	for _, file := range files {
		// Only parse whitelisted file types
>>>>>>> 79bfea2d (update vendor)
=======
	for _, file := range files {
		// Only parse allowlisted file types
>>>>>>> 737a8f1c (add more test case)
		if !crdExts.Has(filepath.Ext(file.Name())) {
			continue
		}

		// Unmarshal CRDs from file into structs
<<<<<<< HEAD
<<<<<<< HEAD
		docs, err := readDocuments(filepath.Join(basePath, file.Name()))
=======
		docs, err := readDocuments(filepath.Join(path, file.Name()))
>>>>>>> 79bfea2d (update vendor)
=======
		docs, err := readDocuments(filepath.Join(basePath, file.Name()))
>>>>>>> 737a8f1c (add more test case)
		if err != nil {
			return nil, err
		}

		for _, doc := range docs {
<<<<<<< HEAD
<<<<<<< HEAD
			crd := &unstructured.Unstructured{}
=======
			crd := &apiextensionsv1beta1.CustomResourceDefinition{}
>>>>>>> 79bfea2d (update vendor)
=======
			crd := &unstructured.Unstructured{}
>>>>>>> 737a8f1c (add more test case)
			if err = yaml.Unmarshal(doc, crd); err != nil {
				return nil, err
			}

			// Check that it is actually a CRD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 737a8f1c (add more test case)
			crdKind, _, err := unstructured.NestedString(crd.Object, "spec", "names", "kind")
			if err != nil {
				return nil, err
			}
			crdGroup, _, err := unstructured.NestedString(crd.Object, "spec", "group")
			if err != nil {
				return nil, err
			}

			if crd.GetKind() != "CustomResourceDefinition" || crdKind == "" || crdGroup == "" {
<<<<<<< HEAD
=======
			if crd.Spec.Names.Kind == "" || crd.Spec.Group == "" {
>>>>>>> 79bfea2d (update vendor)
=======
>>>>>>> 737a8f1c (add more test case)
				continue
			}
			crds = append(crds, crd)
		}

		log.V(1).Info("read CRDs from file", "file", file.Name())
	}
	return crds, nil
}

// readDocuments reads documents from file
func readDocuments(fp string) ([][]byte, error) {
	b, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, err
	}

	docs := [][]byte{}
	reader := k8syaml.NewYAMLReader(bufio.NewReader(bytes.NewReader(b)))
	for {
		// Read document
		doc, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		docs = append(docs, doc)
	}

	return docs, nil
}
