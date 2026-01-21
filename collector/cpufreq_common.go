// Copyright 2023 The Prometheus Authors
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

var (
	cpuFreqHertzDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "frequency_hertz"),
		"Current CPU thread frequency in hertz.",
		[]string{"cpu"}, nil,
	)
	cpuFreqMinDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "frequency_min_hertz"),
		"Minimum CPU thread frequency in hertz.",
		[]string{"cpu"}, nil,
	)
	cpuFreqMaxDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "frequency_max_hertz"),
		"Maximum CPU thread frequency in hertz.",
		[]string{"cpu"}, nil,
	)
	cpuFreqScalingFreqDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "scaling_frequency_hertz"),
		"Current scaled CPU thread frequency in hertz.",
		[]string{"cpu"}, nil,
	)
	cpuFreqScalingFreqMinDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "scaling_frequency_min_hertz"),
		"Minimum scaled CPU thread frequency in hertz.",
		[]string{"cpu"}, nil,
	)
	cpuFreqScalingFreqMaxDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "scaling_frequency_max_hertz"),
		"Maximum scaled CPU thread frequency in hertz.",
		[]string{"cpu"}, nil,
	)
	cpuFreqScalingGovernorDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "scaling_governor"),
		"Current enabled CPU frequency governor.",
		[]string{"cpu", "governor"}, nil,
	)

	// All-core aggregate metrics
	cpuFreqCurrentHzMinDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "frequency_current_hertz_min"),
		"Minimum current CPU frequency in hertz across all cores.",
		nil, nil,
	)
	cpuFreqCurrentHzMeanDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "frequency_current_hertz_mean"),
		"Mean current CPU frequency in hertz across all cores.",
		nil, nil,
	)
	cpuFreqCurrentHzMaxDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "frequency_current_hertz_max"),
		"Maximum current CPU frequency in hertz across all cores.",
		nil, nil,
	)
	cpuFreqCapabilityMinHzDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "frequency_capability_min_hertz"),
		"Minimum CPU frequency capability (hardware limit) in hertz.",
		nil, nil,
	)
	cpuFreqCapabilityMaxHzDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "frequency_capability_max_hertz"),
		"Maximum CPU frequency capability (hardware limit) in hertz.",
		nil, nil,
	)
	cpuFreqScalingCurrentHzMinDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "scaling_frequency_current_hertz_min"),
		"Minimum current scaled CPU frequency in hertz across all cores.",
		nil, nil,
	)
	cpuFreqScalingCurrentHzMeanDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "scaling_frequency_current_hertz_mean"),
		"Mean current scaled CPU frequency in hertz across all cores.",
		nil, nil,
	)
	cpuFreqScalingCurrentHzMaxDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "scaling_frequency_current_hertz_max"),
		"Maximum current scaled CPU frequency in hertz across all cores.",
		nil, nil,
	)
	cpuFreqScalingCapabilityMinHzDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "scaling_frequency_capability_min_hertz"),
		"Minimum scaled CPU frequency capability (hardware limit) in hertz.",
		nil, nil,
	)
	cpuFreqScalingCapabilityMaxHzDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "scaling_frequency_capability_max_hertz"),
		"Maximum scaled CPU frequency capability (hardware limit) in hertz.",
		nil, nil,
	)
	cpuFreqScalingGovernorModalDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "scaling_governor_modal"),
		"Modal (most common) CPU frequency scaling governor across all cores.",
		[]string{"governor"}, nil,
	)
)
