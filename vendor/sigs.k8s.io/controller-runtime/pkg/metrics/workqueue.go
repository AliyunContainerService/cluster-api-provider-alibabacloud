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

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/client-go/util/workqueue"
<<<<<<< HEAD
<<<<<<< HEAD
)

=======
	logf "sigs.k8s.io/controller-runtime/pkg/internal/log"
)

var log = logf.RuntimeLog.WithName("metrics")

>>>>>>> 79bfea2d (update vendor)
=======
)

>>>>>>> e879a141 (alibabacloud machine-api provider)
// This file is copied and adapted from k8s.io/kubernetes/pkg/util/workqueue/prometheus
// which registers metrics to the default prometheus Registry. We require very
// similar functionality, but must register metrics to a different Registry.

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
// Metrics subsystem and all keys used by the workqueue.
const (
	WorkQueueSubsystem         = "workqueue"
	DepthKey                   = "depth"
	AddsKey                    = "adds_total"
	QueueLatencyKey            = "queue_duration_seconds"
	WorkDurationKey            = "work_duration_seconds"
	UnfinishedWorkKey          = "unfinished_work_seconds"
	LongestRunningProcessorKey = "longest_running_processor_seconds"
	RetriesKey                 = "retries_total"
)
<<<<<<< HEAD

var (
	depth = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: WorkQueueSubsystem,
		Name:      DepthKey,
		Help:      "Current depth of workqueue",
	}, []string{"name"})

	adds = prometheus.NewCounterVec(prometheus.CounterOpts{
		Subsystem: WorkQueueSubsystem,
		Name:      AddsKey,
		Help:      "Total number of adds handled by workqueue",
	}, []string{"name"})

	latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem: WorkQueueSubsystem,
		Name:      QueueLatencyKey,
		Help:      "How long in seconds an item stays in workqueue before being requested",
		Buckets:   prometheus.ExponentialBuckets(10e-9, 10, 10),
	}, []string{"name"})

	workDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem: WorkQueueSubsystem,
		Name:      WorkDurationKey,
		Help:      "How long in seconds processing an item from workqueue takes.",
		Buckets:   prometheus.ExponentialBuckets(10e-9, 10, 10),
	}, []string{"name"})

	unfinished = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: WorkQueueSubsystem,
		Name:      UnfinishedWorkKey,
		Help: "How many seconds of work has been done that " +
			"is in progress and hasn't been observed by work_duration. Large " +
			"values indicate stuck threads. One can deduce the number of stuck " +
			"threads by observing the rate at which this increases.",
	}, []string{"name"})

	longestRunningProcessor = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: WorkQueueSubsystem,
		Name:      LongestRunningProcessorKey,
		Help: "How many seconds has the longest running " +
			"processor for workqueue been running.",
	}, []string{"name"})

	retries = prometheus.NewCounterVec(prometheus.CounterOpts{
		Subsystem: WorkQueueSubsystem,
		Name:      RetriesKey,
		Help:      "Total number of retries handled by workqueue",
	}, []string{"name"})
)

func init() {
	Registry.MustRegister(depth)
	Registry.MustRegister(adds)
	Registry.MustRegister(latency)
	Registry.MustRegister(workDuration)
	Registry.MustRegister(unfinished)
	Registry.MustRegister(longestRunningProcessor)
	Registry.MustRegister(retries)

	workqueue.SetProvider(workqueueMetricsProvider{})
}

type workqueueMetricsProvider struct{}

func (workqueueMetricsProvider) NewDepthMetric(name string) workqueue.GaugeMetric {
	return depth.WithLabelValues(name)
}

func (workqueueMetricsProvider) NewAddsMetric(name string) workqueue.CounterMetric {
	return adds.WithLabelValues(name)
}

func (workqueueMetricsProvider) NewLatencyMetric(name string) workqueue.HistogramMetric {
	return latency.WithLabelValues(name)
}

func (workqueueMetricsProvider) NewWorkDurationMetric(name string) workqueue.HistogramMetric {
	return workDuration.WithLabelValues(name)
}

func (workqueueMetricsProvider) NewUnfinishedWorkSecondsMetric(name string) workqueue.SettableGaugeMetric {
	return unfinished.WithLabelValues(name)
}

func (workqueueMetricsProvider) NewLongestRunningProcessorSecondsMetric(name string) workqueue.SettableGaugeMetric {
	return longestRunningProcessor.WithLabelValues(name)
}

func (workqueueMetricsProvider) NewRetriesMetric(name string) workqueue.CounterMetric {
	return retries.WithLabelValues(name)
}
=======
func init() {
	workqueue.SetProvider(workqueueMetricsProvider{})
}

func registerWorkqueueMetric(c prometheus.Collector, name, queue string) {
	if err := Registry.Register(c); err != nil {
		log.Error(err, "failed to register metric", "name", name, "queue", queue)
	}
}

type workqueueMetricsProvider struct{}

func (workqueueMetricsProvider) NewDepthMetric(queue string) workqueue.GaugeMetric {
	const name = "workqueue_depth"
	m := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        name,
		Help:        "Current depth of workqueue",
		ConstLabels: prometheus.Labels{"name": queue},
	})
	registerWorkqueueMetric(m, name, queue)
	return m
}

func (workqueueMetricsProvider) NewAddsMetric(queue string) workqueue.CounterMetric {
	const name = "workqueue_adds_total"
	m := prometheus.NewCounter(prometheus.CounterOpts{
		Name:        name,
		Help:        "Total number of adds handled by workqueue",
		ConstLabels: prometheus.Labels{"name": queue},
	})
	registerWorkqueueMetric(m, name, queue)
	return m
}

func (workqueueMetricsProvider) NewLatencyMetric(queue string) workqueue.HistogramMetric {
	const name = "workqueue_queue_duration_seconds"
	m := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:        name,
		Help:        "How long in seconds an item stays in workqueue before being requested.",
		ConstLabels: prometheus.Labels{"name": queue},
		Buckets:     prometheus.ExponentialBuckets(10e-9, 10, 10),
	})
	registerWorkqueueMetric(m, name, queue)
	return m
}

func (workqueueMetricsProvider) NewWorkDurationMetric(queue string) workqueue.HistogramMetric {
	const name = "workqueue_work_duration_seconds"
	m := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:        name,
		Help:        "How long in seconds processing an item from workqueue takes.",
		ConstLabels: prometheus.Labels{"name": queue},
		Buckets:     prometheus.ExponentialBuckets(10e-9, 10, 10),
	})
	registerWorkqueueMetric(m, name, queue)
	return m
}
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)

var (
	depth = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: WorkQueueSubsystem,
		Name:      DepthKey,
		Help:      "Current depth of workqueue",
	}, []string{"name"})

	adds = prometheus.NewCounterVec(prometheus.CounterOpts{
		Subsystem: WorkQueueSubsystem,
		Name:      AddsKey,
		Help:      "Total number of adds handled by workqueue",
	}, []string{"name"})

	latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem: WorkQueueSubsystem,
		Name:      QueueLatencyKey,
		Help:      "How long in seconds an item stays in workqueue before being requested",
		Buckets:   prometheus.ExponentialBuckets(10e-9, 10, 10),
	}, []string{"name"})

	workDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem: WorkQueueSubsystem,
		Name:      WorkDurationKey,
		Help:      "How long in seconds processing an item from workqueue takes.",
		Buckets:   prometheus.ExponentialBuckets(10e-9, 10, 10),
	}, []string{"name"})

	unfinished = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: WorkQueueSubsystem,
		Name:      UnfinishedWorkKey,
		Help: "How many seconds of work has been done that " +
			"is in progress and hasn't been observed by work_duration. Large " +
			"values indicate stuck threads. One can deduce the number of stuck " +
			"threads by observing the rate at which this increases.",
	}, []string{"name"})

	longestRunningProcessor = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: WorkQueueSubsystem,
		Name:      LongestRunningProcessorKey,
		Help: "How many seconds has the longest running " +
			"processor for workqueue been running.",
	}, []string{"name"})

	retries = prometheus.NewCounterVec(prometheus.CounterOpts{
		Subsystem: WorkQueueSubsystem,
		Name:      RetriesKey,
		Help:      "Total number of retries handled by workqueue",
	}, []string{"name"})
)

func init() {
	Registry.MustRegister(depth)
	Registry.MustRegister(adds)
	Registry.MustRegister(latency)
	Registry.MustRegister(workDuration)
	Registry.MustRegister(unfinished)
	Registry.MustRegister(longestRunningProcessor)
	Registry.MustRegister(retries)

	workqueue.SetProvider(workqueueMetricsProvider{})
}

type workqueueMetricsProvider struct{}

func (workqueueMetricsProvider) NewDepthMetric(name string) workqueue.GaugeMetric {
	return depth.WithLabelValues(name)
}

func (workqueueMetricsProvider) NewAddsMetric(name string) workqueue.CounterMetric {
	return adds.WithLabelValues(name)
}

func (workqueueMetricsProvider) NewLatencyMetric(name string) workqueue.HistogramMetric {
	return latency.WithLabelValues(name)
}

func (workqueueMetricsProvider) NewWorkDurationMetric(name string) workqueue.HistogramMetric {
	return workDuration.WithLabelValues(name)
}

func (workqueueMetricsProvider) NewUnfinishedWorkSecondsMetric(name string) workqueue.SettableGaugeMetric {
	return unfinished.WithLabelValues(name)
}

func (workqueueMetricsProvider) NewLongestRunningProcessorSecondsMetric(name string) workqueue.SettableGaugeMetric {
	return longestRunningProcessor.WithLabelValues(name)
}

<<<<<<< HEAD
type noopMetric struct{}

func (noopMetric) Inc()            {}
func (noopMetric) Dec()            {}
func (noopMetric) Set(float64)     {}
func (noopMetric) Observe(float64) {}
>>>>>>> 79bfea2d (update vendor)
=======
func (workqueueMetricsProvider) NewRetriesMetric(name string) workqueue.CounterMetric {
	return retries.WithLabelValues(name)
}
>>>>>>> e879a141 (alibabacloud machine-api provider)
