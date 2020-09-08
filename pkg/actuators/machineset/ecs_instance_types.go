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

// InstanceType is sepc of ECS instance
type InstanceType struct {
	InstanceType string
	VCPU         int64
	MemoryMb     int64
	GPU          int64
}

// InstanceTypes is a map of ec2 resources
var InstanceTypes = map[string]*InstanceType{
	"ecs.g6.large": {
		InstanceType: "ecs.g6.large",
		VCPU:         2,
		MemoryMb:     8192,
		GPU:          0,
	},
	"ecs.g6.xlarge": {
		InstanceType: "ecs.g6.xlarge",
		VCPU:         4,
		MemoryMb:     16384,
		GPU:          0,
	},
	"ecs.g6.2xlarge": {
		InstanceType: "ecs.g6.2xlarge",
		VCPU:         8,
		MemoryMb:     32768,
		GPU:          0,
	},
	"ecs.g6.3xlarge": {
		InstanceType: "ecs.g6.3xlarge",
		VCPU:         12,
		MemoryMb:     49152,
		GPU:          0,
	},
	"ecs.g6.4xlarge": {
		InstanceType: "ecs.g6.4xlarge",
		VCPU:         16,
		MemoryMb:     65536,
		GPU:          0,
	},
	"ecs.g6.6xlarge": {
		InstanceType: "ecs.g6.6xlarge",
		VCPU:         24,
		MemoryMb:     98304,
		GPU:          0,
	},
	"ecs.g6.8xlarge": {
		InstanceType: "ecs.g6.8xlarge",
		VCPU:         32,
		MemoryMb:     130944,
		GPU:          0,
	},
	"ecs.g6.13xlarge": {
		InstanceType: "ecs.g6.13xlarge",
		VCPU:         52,
		MemoryMb:     196608,
		GPU:          0,
	},
	"ecs.g6.26xlarge": {
		InstanceType: "ecs.g6.26xlarge",
		VCPU:         104,
		MemoryMb:     393216,
		GPU:          0,
	},
	"ecs.g6e.large": {
		InstanceType: "ecs.g6e.large",
		VCPU:         2,
		MemoryMb:     8192,
		GPU:          0,
	},
	"ecs.g6e.xlarge": {
		InstanceType: "ecs.g6e.xlarge",
		VCPU:         4,
		MemoryMb:     16384,
		GPU:          0,
	},
	"ecs.g6e.2xlarge": {
		InstanceType: "ecs.g6e.2xlarge",
		VCPU:         8,
		MemoryMb:     32768,
		GPU:          0,
	},
	"ecs.g6e.4xlarge": {
		InstanceType: "ecs.g6e.4xlarge",
		VCPU:         16,
		MemoryMb:     65536,
		GPU:          0,
	},
	"ecs.g6e.8xlarge": {
		InstanceType: "ecs.g6e.8xlarge",
		VCPU:         32,
		MemoryMb:     130944,
		GPU:          0,
	},
	"ecs.g6e.13xlarge": {
		InstanceType: "ecs.g6e.13xlarge",
		VCPU:         52,
		MemoryMb:     196608,
		GPU:          0,
	},
	"ecs.g6e.26xlarge": {
		InstanceType: "ecs.g6e.26xlarge",
		VCPU:         104,
		MemoryMb:     393216,
		GPU:          0,
	},
	"ecs.g5.large": {
		InstanceType: "ecs.g5.large",
		VCPU:         2,
		MemoryMb:     8192,
		GPU:          0,
	},
	"ecs.g5.xlarge": {
		InstanceType: "ecs.g5.xlarge",
		VCPU:         4,
		MemoryMb:     16384,
		GPU:          0,
	},
	"ecs.g5.2xlarge": {
		InstanceType: "ecs.g5.2xlarge",
		VCPU:         8,
		MemoryMb:     32768,
		GPU:          0,
	},
	"ecs.g5.3xlarge": {
		InstanceType: "ecs.g5.3xlarge",
		VCPU:         12,
		MemoryMb:     49152,
		GPU:          0,
	},
	"ecs.g5.4xlarge": {
		InstanceType: "ecs.g5.4xlarge",
		VCPU:         16,
		MemoryMb:     65536,
		GPU:          0,
	},
	"ecs.g5.6xlarge": {
		InstanceType: "ecs.g5.6xlarge",
		VCPU:         24,
		MemoryMb:     98304,
		GPU:          0,
	},
	"ecs.g5.8xlarge": {
		InstanceType: "ecs.g5.8xlarge",
		VCPU:         32,
		MemoryMb:     130944,
		GPU:          0,
	},
	"ecs.g5.16xlarge": {
		InstanceType: "ecs.g5.16xlarge",
		VCPU:         64,
		MemoryMb:     262144,
		GPU:          0,
	},
	"ecs.c6.large": {
		InstanceType: "ecs.c6.large",
		VCPU:         2,
		MemoryMb:     4096,
		GPU:          0,
	},
	"ecs.c6.xlarge": {
		InstanceType: "ecs.c6.xlarge",
		VCPU:         4,
		MemoryMb:     8192,
		GPU:          0,
	},
	"ecs.c6.2xlarge": {
		InstanceType: "ecs.c6.2xlarge",
		VCPU:         8,
		MemoryMb:     16384,
		GPU:          0,
	},
	"ecs.c6.3xlarge": {
		InstanceType: "ecs.c6.3xlarge",
		VCPU:         12,
		MemoryMb:     24576,
		GPU:          0,
	},
	"ecs.c6.4xlarge": {
		InstanceType: "ecs.c6.4xlarge",
		VCPU:         16,
		MemoryMb:     32768,
		GPU:          0,
	},
	"ecs.c6.6xlarge": {
		InstanceType: "ecs.c6.6xlarge",
		VCPU:         24,
		MemoryMb:     49152,
		GPU:          0,
	},
	"ecs.c6.8xlarge": {
		InstanceType: "ecs.c6.8xlarge",
		VCPU:         32,
		MemoryMb:     65535,
		GPU:          0,
	},
	"ecs.c6.13xlarge": {
		InstanceType: "ecs.c6.13xlarge",
		VCPU:         52,
		MemoryMb:     98304,
		GPU:          0,
	},
	"ecs.c6.26xlarge": {
		InstanceType: "ecs.c6.26xlarge",
		VCPU:         104,
		MemoryMb:     196608,
		GPU:          0,
	},
	"ecs.c6a.large": {
		InstanceType: "ecs.c6a.large",
		VCPU:         2,
		MemoryMb:     4096,
		GPU:          0,
	},
	"ecs.c6a.xlarge": {
		InstanceType: "ecs.c6a.xlarge",
		VCPU:         4,
		MemoryMb:     8192,
		GPU:          0,
	},
	"ecs.c6a.2xlarge": {
		InstanceType: "ecs.c6a.2xlarge",
		VCPU:         8,
		MemoryMb:     16384,
		GPU:          0,
	},
	"ecs.c6a.4xlarge": {
		InstanceType: "ecs.c6a.4xlarge",
		VCPU:         16,
		MemoryMb:     32768,
		GPU:          0,
	},
	"ecs.c6a.8xlarge": {
		InstanceType: "ecs.c6a.8xlarge",
		VCPU:         32,
		MemoryMb:     65535,
		GPU:          0,
	},
	"ecs.c6a.16xlarge": {
		InstanceType: "ecs.c6a.16xlarge",
		VCPU:         64,
		MemoryMb:     131072,
		GPU:          0,
	},
	"ecs.c6a.32xlarge": {
		InstanceType: "ecs.c6a.32xlarge",
		VCPU:         128,
		MemoryMb:     262144,
		GPU:          0,
	},
	"ecs.c6a.64xlarge": {
		InstanceType: "ecs.c6a.64xlarge",
		VCPU:         256,
		MemoryMb:     524288,
		GPU:          0,
	},
	"ecs.g6a.large": {
		InstanceType: "ecs.g6a.large",
		VCPU:         2,
		MemoryMb:     8192,
		GPU:          0,
	},
	"ecs.g6a.xlarge": {
		InstanceType: "ecs.g6a.xlarge",
		VCPU:         4,
		MemoryMb:     16384,
		GPU:          0,
	},
	"ecs.g6a.2xlarge": {
		InstanceType: "ecs.g6a.2xlarge",
		VCPU:         8,
		MemoryMb:     32768,
		GPU:          0,
	},
	"ecs.g6a.4xlarge": {
		InstanceType: "ecs.g6a.4xlarge",
		VCPU:         16,
		MemoryMb:     65536,
		GPU:          0,
	},
	"ecs.g6a.8xlarge": {
		InstanceType: "ecs.g6a.8xlarge",
		VCPU:         32,
		MemoryMb:     131072,
		GPU:          0,
	},
	"ecs.g6a.16xlarge": {
		InstanceType: "ecs.g6a.16xlarge",
		VCPU:         64,
		MemoryMb:     262144,
		GPU:          0,
	},
	"ecs.g6a.32xlarge": {
		InstanceType: "ecs.g6a.32xlarge",
		VCPU:         128,
		MemoryMb:     524288,
		GPU:          0,
	},
	"ecs.g6a.64xlarge": {
		InstanceType: "ecs.g6a.64xlarge",
		VCPU:         256,
		MemoryMb:     1048576,
		GPU:          0,
	},
}
