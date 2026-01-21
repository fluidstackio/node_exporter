// Copyright 2021 The Prometheus Authors
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
	"io"
	"log/slog"
	"reflect"
	"testing"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

func copyStats(d, s map[int64]procfs.CPUStat) {
	for k := range s {
		v := s[k]
		d[k] = v
	}
}

func makeTestCPUCollector(s map[int64]procfs.CPUStat) *cpuCollector {
	dup := make(map[int64]procfs.CPUStat, len(s))
	copyStats(dup, s)
	return &cpuCollector{
		logger:   slog.New(slog.NewTextHandler(io.Discard, nil)),
		cpuStats: dup,
	}
}

// Helper function to extract metric value and labels from a prometheus.Metric
func readMetric(m prometheus.Metric) (float64, map[string]string) {
	pb := &dto.Metric{}
	if err := m.Write(pb); err != nil {
		panic(err)
	}

	labels := make(map[string]string)
	for _, lp := range pb.Label {
		labels[*lp.Name] = *lp.Value
	}

	var value float64
	if pb.Gauge != nil {
		value = pb.Gauge.GetValue()
	} else if pb.Counter != nil {
		value = pb.Counter.GetValue()
	} else if pb.Untyped != nil {
		value = pb.Untyped.GetValue()
	}

	return value, labels
}

func TestCPU(t *testing.T) {
	firstCPUStat := map[int64]procfs.CPUStat{
		0: {
			User:      100.0,
			Nice:      100.0,
			System:    100.0,
			Idle:      100.0,
			Iowait:    100.0,
			IRQ:       100.0,
			SoftIRQ:   100.0,
			Steal:     100.0,
			Guest:     100.0,
			GuestNice: 100.0,
		}}

	c := makeTestCPUCollector(firstCPUStat)
	want := map[int64]procfs.CPUStat{
		0: {
			User:      101.0,
			Nice:      101.0,
			System:    101.0,
			Idle:      101.0,
			Iowait:    101.0,
			IRQ:       101.0,
			SoftIRQ:   101.0,
			Steal:     101.0,
			Guest:     101.0,
			GuestNice: 101.0,
		}}
	c.updateCPUStats(want)
	got := c.cpuStats
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("should have %v CPU Stat: got %v", want, got)
	}

	c = makeTestCPUCollector(firstCPUStat)
	jumpBack := map[int64]procfs.CPUStat{
		0: {
			User:      99.9,
			Nice:      99.9,
			System:    99.9,
			Idle:      99.9,
			Iowait:    99.9,
			IRQ:       99.9,
			SoftIRQ:   99.9,
			Steal:     99.9,
			Guest:     99.9,
			GuestNice: 99.9,
		}}
	c.updateCPUStats(jumpBack)
	got = c.cpuStats
	if reflect.DeepEqual(jumpBack, got) {
		t.Fatalf("should have %v CPU Stat: got %v", firstCPUStat, got)
	}

	c = makeTestCPUCollector(firstCPUStat)
	resetIdle := map[int64]procfs.CPUStat{
		0: {
			User:      102.0,
			Nice:      102.0,
			System:    102.0,
			Idle:      1.0,
			Iowait:    102.0,
			IRQ:       102.0,
			SoftIRQ:   102.0,
			Steal:     102.0,
			Guest:     102.0,
			GuestNice: 102.0,
		}}
	c.updateCPUStats(resetIdle)
	got = c.cpuStats
	if !reflect.DeepEqual(resetIdle, got) {
		t.Fatalf("should have %v CPU Stat: got %v", resetIdle, got)
	}
}

func TestCPUOffline(t *testing.T) {
	// CPU 1 goes offline.
	firstCPUStat := map[int64]procfs.CPUStat{
		0: {
			User:      100.0,
			Nice:      100.0,
			System:    100.0,
			Idle:      100.0,
			Iowait:    100.0,
			IRQ:       100.0,
			SoftIRQ:   100.0,
			Steal:     100.0,
			Guest:     100.0,
			GuestNice: 100.0,
		},
		1: {
			User:      101.0,
			Nice:      101.0,
			System:    101.0,
			Idle:      101.0,
			Iowait:    101.0,
			IRQ:       101.0,
			SoftIRQ:   101.0,
			Steal:     101.0,
			Guest:     101.0,
			GuestNice: 101.0,
		},
	}

	c := makeTestCPUCollector(firstCPUStat)
	want := map[int64]procfs.CPUStat{
		0: {
			User:      100.0,
			Nice:      100.0,
			System:    100.0,
			Idle:      100.0,
			Iowait:    100.0,
			IRQ:       100.0,
			SoftIRQ:   100.0,
			Steal:     100.0,
			Guest:     100.0,
			GuestNice: 100.0,
		},
	}
	c.updateCPUStats(want)
	got := c.cpuStats
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("should have %v CPU Stat: got %v", want, got)
	}

	// CPU 1 comes back online.
	want = map[int64]procfs.CPUStat{
		0: {
			User:      100.0,
			Nice:      100.0,
			System:    100.0,
			Idle:      100.0,
			Iowait:    100.0,
			IRQ:       100.0,
			SoftIRQ:   100.0,
			Steal:     100.0,
			Guest:     100.0,
			GuestNice: 100.0,
		},
		1: {
			User:      101.0,
			Nice:      101.0,
			System:    101.0,
			Idle:      101.0,
			Iowait:    101.0,
			IRQ:       101.0,
			SoftIRQ:   101.0,
			Steal:     101.0,
			Guest:     101.0,
			GuestNice: 101.0,
		},
	}
	c.updateCPUStats(want)
	got = c.cpuStats
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("should have %v CPU Stat: got %v", want, got)
	}

}

func TestCPUSecondsAllCoreMean(t *testing.T) {
	// Test with multiple CPUs with different values.
	cpuStats := map[int64]procfs.CPUStat{
		0: {
			User:    100.0,
			Nice:    20.0,
			System:  30.0,
			Idle:    400.0,
			Iowait:  10.0,
			IRQ:     5.0,
			SoftIRQ: 5.0,
			Steal:   0.0,
		},
		1: {
			User:    200.0,
			Nice:    40.0,
			System:  60.0,
			Idle:    800.0,
			Iowait:  20.0,
			IRQ:     10.0,
			SoftIRQ: 10.0,
			Steal:   0.0,
		},
	}

	c := makeTestCPUCollector(cpuStats)
	c.cpuSecondsAllCoreMean = nodeCPUSecondsAllCoreMeanDesc

	ch := make(chan prometheus.Metric, 8)
	c.emitCPUSecondsAllCoreMean(ch)
	close(ch)

	// Expected means (sum / 2)
	expectedValues := map[string]float64{
		"user":    150.0, // (100 + 200) / 2
		"nice":    30.0,  // (20 + 40) / 2
		"system":  45.0,  // (30 + 60) / 2
		"idle":    600.0, // (400 + 800) / 2
		"iowait":  15.0,  // (10 + 20) / 2
		"irq":     7.5,   // (5 + 10) / 2
		"softirq": 7.5,   // (5 + 10) / 2
		"steal":   0.0,   // (0 + 0) / 2
	}

	count := 0
	for metric := range ch {
		value, labels := readMetric(metric)
		mode := labels["mode"]
		expectedValue, ok := expectedValues[mode]
		if !ok {
			t.Fatalf("unexpected mode: %s", mode)
		}
		if value != expectedValue {
			t.Errorf("mode %s: expected value %f, got %f", mode, expectedValue, value)
		}
		count++
	}

	if count != 8 {
		t.Fatalf("expected 8 metrics, got %d", count)
	}
}

func TestCPUSecondsAllCoreMeanSingleCPU(t *testing.T) {
	// Test with a single CPU - mean should equal the CPU value.
	cpuStats := map[int64]procfs.CPUStat{
		0: {
			User:    100.0,
			Nice:    20.0,
			System:  30.0,
			Idle:    400.0,
			Iowait:  10.0,
			IRQ:     5.0,
			SoftIRQ: 5.0,
			Steal:   0.0,
		},
	}

	c := makeTestCPUCollector(cpuStats)
	c.cpuSecondsAllCoreMean = nodeCPUSecondsAllCoreMeanDesc

	ch := make(chan prometheus.Metric, 8)
	c.emitCPUSecondsAllCoreMean(ch)
	close(ch)

	// With a single CPU, the mean should be the same as the CPU values.
	// Verify we got 8 metrics.
	count := 0
	for range ch {
		count++
	}
	if count != 8 {
		t.Fatalf("expected 8 metrics, got %d", count)
	}
}

func TestCPUSecondsAllCoreMeanNoCPUs(t *testing.T) {
	// Test with no CPUs - should not emit any metrics.
	cpuStats := map[int64]procfs.CPUStat{}

	c := makeTestCPUCollector(cpuStats)
	c.cpuSecondsAllCoreMean = nodeCPUSecondsAllCoreMeanDesc

	ch := make(chan prometheus.Metric, 8)
	c.emitCPUSecondsAllCoreMean(ch)
	close(ch)

	// Should not emit any metrics when there are no CPUs.
	count := 0
	for range ch {
		count++
	}
	if count != 0 {
		t.Fatalf("expected 0 metrics, got %d", count)
	}
}

func TestCPUGuestSecondsAllCoreMean(t *testing.T) {
	// Test guest seconds all-core mean with multiple CPUs.
	cpuStats := map[int64]procfs.CPUStat{
		0: {
			Guest:     50.0,
			GuestNice: 10.0,
		},
		1: {
			Guest:     100.0,
			GuestNice: 20.0,
		},
	}

	c := makeTestCPUCollector(cpuStats)
	c.cpuGuestSecondsAllCoreMean = nodeCPUGuestSecondsAllCoreMeanDesc

	// Enable CPU guest metrics
	enableGuestTrue := true
	oldEnableCPUGuest := enableCPUGuest
	enableCPUGuest = &enableGuestTrue
	defer func() { enableCPUGuest = oldEnableCPUGuest }()

	ch := make(chan prometheus.Metric, 2)
	c.emitCPUGuestSecondsAllCoreMean(ch)
	close(ch)

	// Expected: guest=(50+100)/2=75, guest_nice=(10+20)/2=15
	expectedValues := map[string]float64{
		"user": 75.0, // (50 + 100) / 2
		"nice": 15.0, // (10 + 20) / 2
	}

	count := 0
	for metric := range ch {
		value, labels := readMetric(metric)
		mode := labels["mode"]
		expectedValue, ok := expectedValues[mode]
		if !ok {
			t.Fatalf("unexpected mode: %s", mode)
		}
		if value != expectedValue {
			t.Errorf("mode %s: expected value %f, got %f", mode, expectedValue, value)
		}
		count++
	}

	if count != 2 {
		t.Fatalf("expected 2 metrics, got %d", count)
	}
}

func TestCPUGuestSecondsAllCoreMeanDisabled(t *testing.T) {
	// Test that guest metrics are not emitted when disabled.
	cpuStats := map[int64]procfs.CPUStat{
		0: {
			Guest:     50.0,
			GuestNice: 10.0,
		},
	}

	c := makeTestCPUCollector(cpuStats)
	c.cpuGuestSecondsAllCoreMean = nodeCPUGuestSecondsAllCoreMeanDesc

	// Disable CPU guest metrics
	enableGuestFalse := false
	oldEnableCPUGuest := enableCPUGuest
	enableCPUGuest = &enableGuestFalse
	defer func() { enableCPUGuest = oldEnableCPUGuest }()

	ch := make(chan prometheus.Metric, 2)
	c.emitCPUGuestSecondsAllCoreMean(ch)
	close(ch)

	count := 0
	for range ch {
		count++
	}
	if count != 0 {
		t.Fatalf("expected 0 metrics when guest is disabled, got %d", count)
	}
}

func TestCPUFrequencyAllCoreAggregates(t *testing.T) {
	// Test frequency min/mean/max aggregates.
	info := []procfs.CPUInfo{
		{CPUMHz: 2500},
		{CPUMHz: 3000},
		{CPUMHz: 3500},
	}

	c := &cpuCollector{
		logger:                    slog.New(slog.NewTextHandler(io.Discard, nil)),
		cpuFrequencyHzAllCoreMin:  nodeCPUFrequencyHzAllCoreMinDesc,
		cpuFrequencyHzAllCoreMean: nodeCPUFrequencyHzAllCoreMeanDesc,
		cpuFrequencyHzAllCoreMax:  nodeCPUFrequencyHzAllCoreMaxDesc,
	}

	ch := make(chan prometheus.Metric, 3)
	c.emitCPUFrequencyAllCoreAggregates(ch, info)
	close(ch)

	// Expected: min=2500MHz, mean=3000MHz, max=3500MHz (converted to Hz)
	metrics := make([]prometheus.Metric, 0, 3)
	for metric := range ch {
		metrics = append(metrics, metric)
	}

	if len(metrics) != 3 {
		t.Fatalf("expected 3 metrics, got %d", len(metrics))
	}

	// Collect values by metric name (simplified by checking descriptor match)
	var minValue, meanValue, maxValue float64
	for _, metric := range metrics {
		value, _ := readMetric(metric)
		desc := metric.Desc().String()
		if desc == nodeCPUFrequencyHzAllCoreMinDesc.String() {
			minValue = value
		} else if desc == nodeCPUFrequencyHzAllCoreMeanDesc.String() {
			meanValue = value
		} else if desc == nodeCPUFrequencyHzAllCoreMaxDesc.String() {
			maxValue = value
		}
	}

	expectedMin := 2500.0 * 1e6
	expectedMean := 3000.0 * 1e6
	expectedMax := 3500.0 * 1e6

	if minValue != expectedMin {
		t.Errorf("expected min frequency %f, got %f", expectedMin, minValue)
	}
	if meanValue != expectedMean {
		t.Errorf("expected mean frequency %f, got %f", expectedMean, meanValue)
	}
	if maxValue != expectedMax {
		t.Errorf("expected max frequency %f, got %f", expectedMax, maxValue)
	}
}

func TestCPUFrequencyAllCoreAggregatesEmpty(t *testing.T) {
	// Test with no CPU info - should not emit any metrics.
	c := &cpuCollector{
		logger:                    slog.New(slog.NewTextHandler(io.Discard, nil)),
		cpuFrequencyHzAllCoreMin:  nodeCPUFrequencyHzAllCoreMinDesc,
		cpuFrequencyHzAllCoreMean: nodeCPUFrequencyHzAllCoreMeanDesc,
		cpuFrequencyHzAllCoreMax:  nodeCPUFrequencyHzAllCoreMaxDesc,
	}

	ch := make(chan prometheus.Metric, 3)
	c.emitCPUFrequencyAllCoreAggregates(ch, []procfs.CPUInfo{})
	close(ch)

	count := 0
	for range ch {
		count++
	}
	if count != 0 {
		t.Fatalf("expected 0 metrics, got %d", count)
	}
}

func TestCPUInfoAllCoreAggregate(t *testing.T) {
	// Test CPU info aggregate with varied data.
	info := []procfs.CPUInfo{
		{
			VendorID:   "GenuineIntel",
			CPUFamily:  "6",
			Model:      "85",
			ModelName:  "Intel(R) Xeon(R) Gold 6154 CPU @ 3.00GHz",
			Microcode:  "0x2006e05",
			Stepping:   "4",
			CacheSize:  "25344 KB",
		},
		{
			VendorID:   "GenuineIntel",
			CPUFamily:  "6",
			Model:      "85",
			ModelName:  "Intel(R) Xeon(R) Gold 6154 CPU @ 3.00GHz",
			Microcode:  "0x2006e05",
			Stepping:   "4",
			CacheSize:  "25344 KB",
		},
	}

	c := &cpuCollector{
		logger:                  slog.New(slog.NewTextHandler(io.Discard, nil)),
		cpuInfoAllCoreAggregate: nodeCPUInfoAllCoreAggregateDesc,
	}

	// Should emit 1 metric with modal values
	ch := make(chan prometheus.Metric, 1)
	c.emitCPUInfoAllCoreAggregate(ch, info)
	close(ch)

	count := 0
	for range ch {
		count++
	}
	if count != 1 {
		t.Fatalf("expected 1 metric, got %d", count)
	}
}

