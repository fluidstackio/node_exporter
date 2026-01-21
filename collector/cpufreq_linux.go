// Copyright 2019 The Prometheus Authors
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
	"fmt"
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
)

type cpuFreqCollector struct {
	fs     sysfs.FS
	logger *slog.Logger
}

func init() {
	registerCollector("cpufreq", defaultEnabled, NewCPUFreqCollector)
}

// NewCPUFreqCollector returns a new Collector exposing kernel/system statistics.
func NewCPUFreqCollector(logger *slog.Logger) (Collector, error) {
	fs, err := sysfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}

	return &cpuFreqCollector{
		fs:     fs,
		logger: logger,
	}, nil
}

// Update implements Collector and exposes cpu related metrics from /proc/stat and /sys/.../cpu/.
func (c *cpuFreqCollector) Update(ch chan<- prometheus.Metric) error {
	cpuFreqs, err := c.fs.SystemCpufreq()
	if err != nil {
		return err
	}

	if len(cpuFreqs) == 0 {
		return nil
	}

	// sysfs cpufreq values are kHz, thus multiply by 1000 to export base units (hz).
	// See https://www.kernel.org/doc/Documentation/cpu-freq/user-guide.txt

	// Aggregate values across all cores (per-core metrics removed for cardinality reduction)
	var (
		currentFreqs     []float64
		minFreqs         []float64
		maxFreqs         []float64
		scalingFreqs     []float64
		scalingMinFreqs  []float64
		scalingMaxFreqs  []float64
		governorCounts   = make(map[string]int)
	)

	for _, stats := range cpuFreqs {
		if stats.CpuinfoCurrentFrequency != nil {
			currentFreqs = append(currentFreqs, float64(*stats.CpuinfoCurrentFrequency)*1000.0)
		}
		if stats.CpuinfoMinimumFrequency != nil {
			minFreqs = append(minFreqs, float64(*stats.CpuinfoMinimumFrequency)*1000.0)
		}
		if stats.CpuinfoMaximumFrequency != nil {
			maxFreqs = append(maxFreqs, float64(*stats.CpuinfoMaximumFrequency)*1000.0)
		}
		if stats.ScalingCurrentFrequency != nil {
			scalingFreqs = append(scalingFreqs, float64(*stats.ScalingCurrentFrequency)*1000.0)
		}
		if stats.ScalingMinimumFrequency != nil {
			scalingMinFreqs = append(scalingMinFreqs, float64(*stats.ScalingMinimumFrequency)*1000.0)
		}
		if stats.ScalingMaximumFrequency != nil {
			scalingMaxFreqs = append(scalingMaxFreqs, float64(*stats.ScalingMaximumFrequency)*1000.0)
		}
		if stats.Governor != "" {
			governorCounts[stats.Governor]++
		}
	}

	// Emit current frequency aggregates (min/mean/max)
	if len(currentFreqs) > 0 {
		min, mean, max := calculateMinMeanMax(currentFreqs)
		ch <- prometheus.MustNewConstMetric(cpuFreqHertzAllCoreMinDesc, prometheus.GaugeValue, min)
		ch <- prometheus.MustNewConstMetric(cpuFreqHertzAllCoreMeanDesc, prometheus.GaugeValue, mean)
		ch <- prometheus.MustNewConstMetric(cpuFreqHertzAllCoreMaxDesc, prometheus.GaugeValue, max)
	}

	// Emit min frequency aggregate (absolute minimum)
	if len(minFreqs) > 0 {
		min := minFreqs[0]
		for _, v := range minFreqs {
			if v < min {
				min = v
			}
		}
		ch <- prometheus.MustNewConstMetric(cpuFreqMinAllCoreMinDesc, prometheus.GaugeValue, min)
	}

	// Emit max frequency aggregate (absolute maximum)
	if len(maxFreqs) > 0 {
		max := maxFreqs[0]
		for _, v := range maxFreqs {
			if v > max {
				max = v
			}
		}
		ch <- prometheus.MustNewConstMetric(cpuFreqMaxAllCoreMaxDesc, prometheus.GaugeValue, max)
	}

	// Emit scaling current frequency aggregates (min/mean/max)
	if len(scalingFreqs) > 0 {
		min, mean, max := calculateMinMeanMax(scalingFreqs)
		ch <- prometheus.MustNewConstMetric(cpuFreqScalingFreqAllCoreMinDesc, prometheus.GaugeValue, min)
		ch <- prometheus.MustNewConstMetric(cpuFreqScalingFreqAllCoreMeanDesc, prometheus.GaugeValue, mean)
		ch <- prometheus.MustNewConstMetric(cpuFreqScalingFreqAllCoreMaxDesc, prometheus.GaugeValue, max)
	}

	// Emit scaling min frequency aggregate (absolute minimum)
	if len(scalingMinFreqs) > 0 {
		min := scalingMinFreqs[0]
		for _, v := range scalingMinFreqs {
			if v < min {
				min = v
			}
		}
		ch <- prometheus.MustNewConstMetric(cpuFreqScalingFreqMinAllCoreMinDesc, prometheus.GaugeValue, min)
	}

	// Emit scaling max frequency aggregate (absolute maximum)
	if len(scalingMaxFreqs) > 0 {
		max := scalingMaxFreqs[0]
		for _, v := range scalingMaxFreqs {
			if v > max {
				max = v
			}
		}
		ch <- prometheus.MustNewConstMetric(cpuFreqScalingFreqMaxAllCoreMaxDesc, prometheus.GaugeValue, max)
	}

	// Emit governor aggregate (most common governor)
	if len(governorCounts) > 0 {
		maxCount := 0
		modalGovernor := ""
		for gov, count := range governorCounts {
			if count > maxCount {
				maxCount = count
				modalGovernor = gov
			}
		}
		ch <- prometheus.MustNewConstMetric(cpuFreqScalingGovernorAllCoreAggregateDesc, prometheus.GaugeValue, 1, modalGovernor)
	}

	return nil
}

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
