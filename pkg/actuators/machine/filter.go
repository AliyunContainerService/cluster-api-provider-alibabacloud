package machine

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

const (
	clusterFilterKeyPrefix = "kubernetes.io/cluster/"
	clusterFilterValue     = "owned"
)

func clusterTagFilter(clusterId, machineName string) []ecs.DescribeInstancesTag {
	tagsList := make([]ecs.DescribeInstancesTag, 0)
	tagsList = append(tagsList, ecs.DescribeInstancesTag{
		Key:   fmt.Sprintf("%s%s", clusterFilterKeyPrefix, clusterId),
		Value: clusterFilterValue,
	})
	tagsList = append(tagsList, ecs.DescribeInstancesTag{
		Key:   "Name",
		Value: machineName,
	})

	return tagsList
}
