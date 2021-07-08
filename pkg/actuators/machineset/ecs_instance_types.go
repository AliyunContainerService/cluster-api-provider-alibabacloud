/*
Copyright The Kubernetes Authors.
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

package machineset

type instanceType struct {
	InstanceType string
	VCPU         int64
	MemoryMb     int64
	GPU          int64
}

// InstanceTypes is a map of ecs resources
<<<<<<< HEAD
<<<<<<< HEAD
// TODO next version will be supported
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
// TODO next version will be supported
>>>>>>> 24c35849 (fix stop ecs instance func)
var InstanceTypes = map[string]*instanceType{
	"ecs.c6.2xlarge": {
		InstanceType: "ecs.c6.2xlarge",
		VCPU:         8,
		MemoryMb:     16384,
		GPU:          0,
	},
}
