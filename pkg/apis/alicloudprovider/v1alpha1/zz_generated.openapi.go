// +build !

/*
Copyright 2019 The Kubernetes Authors.
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

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1alpha1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/apis/alibabacloudproviderconfig/v1alpha1.AlibabaCloudMachineProviderConfig": schema_pkg_apis_alibabacloudproviderconfig_v1alpha1_AlibabaCloudMachineProviderConfig(ref),
		"github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/apis/alibabacloudproviderconfig/v1alpha1.AlibabaCloudMachineProviderSpec":   schema_pkg_apis_alibabacloudproviderconfig_v1alpha1_AlibabaCloudMachineProviderConfigSpec(ref),
		"github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/apis/alibabacloudproviderconfig/v1alpha1.AlibabaCloudMachineProviderStatus": schema_pkg_apis_alibabacloudproviderconfig_v1alpha1_AlibabaCloudMachineProviderConfigStatus(ref),
	}
}

func schema_pkg_apis_alibabacloudproviderconfig_v1alpha1_AlibabaCloudMachineProviderConfig(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "AlibabaCloudMachineProviderConfig is the Schema for the alicloudmachineproviderconfigs API",
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/apis/alicloudproviderconfig/v1alpha1.AlicloudMachineProviderSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/apis/alicloudproviderconfig/v1alpha1.AlibabaCloudMachineProviderStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/apis/alicloudproviderconfig/v1alpha1.AlicloudMachineProviderSpec", "github.com/AliyunContainerService/cluster-api-provider-alibabacloud/pkg/apis/alicloudproviderconfig/v1alpha1.AlibabaCloudMachineProviderStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_alibabacloudproviderconfig_v1alpha1_AlibabaCloudMachineProviderConfigSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "AlicloudMachineProviderSpec defines the desired state of AlibabaCloudMachineProviderConfig",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_alibabacloudproviderconfig_v1alpha1_AlibabaCloudMachineProviderConfigStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "AlibabaCloudMachineProviderStatus defines the observed state of AlibabaCloudMachineProviderConfig",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}
