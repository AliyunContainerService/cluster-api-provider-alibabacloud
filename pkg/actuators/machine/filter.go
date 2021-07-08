/*
Copyright 2021 The Kubernetes Authors.

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

package machine

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

const (
	clusterFilterKeyPrefix = "kubernetes.io/cluster/"
	clusterFilterValue     = "owned"
	clusterFilterName      = "Name"
	clusterOwnedKey        = "kubernetes-sigs/cluster-api"
	clusterOwnedValue      = "cluster-api-provider-alibabacloud"
)

func clusterTagFilter(clusterID, machineName string) []ecs.DescribeInstancesTag {
	tagsList := make([]ecs.DescribeInstancesTag, 0)
	tagsList = append(tagsList, ecs.DescribeInstancesTag{
		Key:   fmt.Sprintf("%s%s", clusterFilterKeyPrefix, clusterID),
		Value: clusterFilterValue,
	})
	tagsList = append(tagsList, ecs.DescribeInstancesTag{
		Key:   clusterFilterName,
		Value: machineName,
	})

	return tagsList
}

func tagResourceTags(clusterID, machineName string) *[]ecs.TagResourcesTag {
	tagsList := make([]ecs.TagResourcesTag, 0)

	tagsList = append(tagsList, ecs.TagResourcesTag{
		Key:   fmt.Sprintf("%s%s", clusterFilterKeyPrefix, clusterID),
		Value: clusterFilterValue,
	})
	tagsList = append(tagsList, ecs.TagResourcesTag{
		Key:   clusterFilterName,
		Value: machineName,
	})
	tagsList = append(tagsList, ecs.TagResourcesTag{
		Key:   clusterOwnedKey,
		Value: clusterOwnedValue,
	})

	return &tagsList
}
