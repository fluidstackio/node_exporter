// Copyright 2016 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build !nocpu
// +build !nocpu

package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	cpuCollectorSubsystem = "cpu"
)

var (
	// nodeCPUSecondsDesc is used by non-Linux platforms (Darwin, FreeBSD, etc.)
	// Linux uses aggregate metrics instead to reduce cardinality
	nodeCPUSecondsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "seconds_total"),
		"Seconds the CPUs spent in each mode.",
		[]string{"cpu", "mode"}, nil,
	)
	nodeCPUSecondsAllCoreMeanDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "seconds_all_core_mean"),
		"Mean seconds the CPUs spent in each mode across all cores.",
		[]string{"mode"}, nil,
	)
	nodeCPUGuestSecondsAllCoreMeanDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "guest_seconds_all_core_mean"),
		"Mean seconds the CPUs spent in guest mode across all cores.",
		[]string{"mode"}, nil,
	)
	nodeCPUInfoAllCoreAggregateDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "info_all_core_aggregate"),
		"Aggregate CPU information across all cores (modal values).",
		[]string{"vendor", "family", "model", "model_name", "microcode", "stepping", "cachesize"}, nil,
	)
	nodeCPUFrequencyHzAllCoreMinDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "frequency_hertz_all_core_min"),
		"Minimum CPU frequency in hertz across all cores.",
		nil, nil,
	)
	nodeCPUFrequencyHzAllCoreMeanDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "frequency_hertz_all_core_mean"),
		"Mean CPU frequency in hertz across all cores.",
		nil, nil,
	)
	nodeCPUFrequencyHzAllCoreMaxDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "frequency_hertz_all_core_max"),
		"Maximum CPU frequency in hertz across all cores.",
		nil, nil,
	)
	nodeCPUThrottlesAllCoreTotalDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "throttles_all_core_total"),
		"Total number of CPU core throttle events across all cores.",
		nil, nil,
	)
	nodeCPUThrottlesAllPackageTotalDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "throttles_all_package_total"),
		"Total number of CPU package throttle events across all packages.",
		nil, nil,
	)
)

// calculateMinMeanMax calculates min, mean, and max of a slice of floats
func calculateMinMeanMax(values []float64) (min, mean, max float64) {
	if len(values) == 0 {
		return 0, 0, 0
	}

	min = values[0]
	max = values[0]
	sum := 0.0

	for _, v := range values {
		sum += v
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	mean = sum / float64(len(values))
	return min, mean, max
}

// findModalValue finds the most common value in a count map.
// Uses lexicographic ordering for deterministic tie-breaking.
func findModalValue(counts map[string]int) string {
	maxCount := 0
	modal := ""
	for val, count := range counts {
		// Use lexicographic ordering for deterministic tie-breaking
		if count > maxCount || (count == maxCount && (modal == "" || val < modal)) {
			maxCount = count
			modal = val
		}
	}
	return modal
}

// findMin finds the minimum value in a slice of floats.
// Returns 0 if the slice is empty.
func findMin(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	min := values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
	}
	return min
}

// findMax finds the maximum value in a slice of floats.
// Returns 0 if the slice is empty.
func findMax(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	max := values[0]
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	return max
}
