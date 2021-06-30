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
<<<<<<< HEAD
<<<<<<< HEAD
=======

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

func Test_DescribeImages(t *testing.T) {
	ecsClient, err := ecs.NewClientWithAccessKey("cn-beijing", "LTAI5tQ9JSqZeJShVqVghmdE", "eitOmHrNw53JdkWkjUPhCtV8ca93Vj")
	assert.Nil(t, err)
	assert.NotEmpty(t, ecsClient)

	request := ecs.CreateDescribeInstancesRequest()
	request.RegionId = "cn-beijing"
	instancesIds, _ := json.Marshal([]string{"i-2ze8nunk5kt23u21lkjt"})
	request.InstanceIds = string(instancesIds)
	request.Scheme = "https"

	response, err := ecsClient.DescribeInstances(request)
	assert.Nil(t, err)
	t.Logf("%v", response)
}
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
>>>>>>> 7e2c5241 (remove test case)
