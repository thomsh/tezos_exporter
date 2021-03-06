package collector

import (
	"runtime/debug"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	GitRevision string
	GitBranch   string
)

type constCollector struct {
	metric prometheus.Metric
}

// NewBuildInfoCollector returns a collector collecting a single metric "go_build_info"
func NewBuildInfoCollector(prefix string) prometheus.Collector {
	var path, version, sum, revision, branch = "(unknown)", "(unknown)", "(unknown)", "(unknown)", "(unknown)"
	if bi, ok := debug.ReadBuildInfo(); ok {
		path = bi.Main.Path
		version = bi.Main.Version
		sum = bi.Main.Sum
	}

	if GitRevision != "" {
		revision = GitRevision
	}

	if GitBranch != "" {
		branch = GitBranch
	}

	if prefix == "" {
		prefix = "go"
	}

	return &constCollector{
		metric: prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prefix+"_build_info",
				"Build information about the main Go module.",
				nil, prometheus.Labels{
					"path":     path,
					"version":  version,
					"checksum": sum,
					"revision": revision,
					"branch":   branch,
				},
			),
			prometheus.GaugeValue, 1),
	}
}

// Describe implements prometheus.Collector
func (c *constCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.metric.Desc()
}

// Collect implements prometheus.Collector
func (c *constCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- c.metric
}
