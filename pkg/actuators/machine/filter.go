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

func alicloudTagFilter(name string) string {
	return fmt.Sprint("tag:", name)
}

func clusterFilterKey(name string) string {
	return fmt.Sprint(clusterFilterKeyPrefix, name)
}

func clusterFilter(name string) *ecs.DescribeInstancesTag {
	return &ecs.DescribeInstancesTag{
		Key:   alicloudTagFilter(clusterFilterKey(name)),
		Value: clusterFilterValue,
	}
}
