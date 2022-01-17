package main

import (
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type lvmLvCollector struct {
	lvTotalSizeMetric       *prometheus.Desc
	lvDataPercentMetric     *prometheus.Desc
	lvSnapPercentMetric     *prometheus.Desc
	lvMetadataPercentMetric *prometheus.Desc
	node                    string
}

func newLvmLvCollector(node string) *lvmLvCollector {
	return &lvmLvCollector{
		lvTotalSizeMetric: prometheus.NewDesc("lvm_lv_total_size_bytes",
			"Shows LVM LV total size in Bytes",
			[]string{"lv_name", "vg_name", "node"}, nil,
		),
		lvDataPercentMetric: prometheus.NewDesc("lvm_lv_data_percent",
			"Shows LVM LV data percent",
			[]string{"lv_name", "vg_name", "node"}, nil,
		),
		lvSnapPercentMetric: prometheus.NewDesc("lvm_lv_snap_percent",
			"Shows LVM LV snap percent",
			[]string{"lv_name", "vg_name", "node"}, nil,
		),
		lvMetadataPercentMetric: prometheus.NewDesc("lvm_lv_metadata_percent",
			"Shows LVM LV metadata percent",
			[]string{"lv_name", "vg_name", "node"}, nil,
		),
		node: node,
	}
}

func (collector *lvmLvCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.lvTotalSizeMetric
}

func (collector *lvmLvCollector) Collect(ch chan<- prometheus.Metric) {
	out, err := exec.Command("/sbin/lvs", "--units", "B", "--separator", ",", "-o", "lv_size,lv_name,vg_name,data_percent,snap_percent,metadata_percent", "--noheadings").Output()
	if err != nil {
		log.Print(err)
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		values := strings.Split(strings.TrimSpace(line), ",")
		if len(values) < 6 {
			continue
		}

		logicalVolumeName := values[1]
		volumeGroupName := values[2]
		size, err := strconv.ParseFloat(strings.Trim(values[0], "B"), 64)
		if err != nil {
			continue
		}
		dataPercent, err := strconv.ParseFloat(values[3], 64)
		if err != nil {
			continue
		}
		snapPercent, err := strconv.ParseFloat(values[4], 64)
		if err != nil {
			continue
		}
		metadataPercent, err := strconv.ParseFloat(values[5], 64)
		if err != nil {
			continue
		}

		ch <- prometheus.MustNewConstMetric(collector.lvTotalSizeMetric, prometheus.GaugeValue, size, logicalVolumeName, volumeGroupName, collector.node)
		ch <- prometheus.MustNewConstMetric(collector.lvDataPercentMetric, prometheus.GaugeValue, dataPercent, logicalVolumeName, volumeGroupName, collector.node)
		ch <- prometheus.MustNewConstMetric(collector.lvSnapPercentMetric, prometheus.GaugeValue, snapPercent, logicalVolumeName, volumeGroupName, collector.node)
		ch <- prometheus.MustNewConstMetric(collector.lvMetadataPercentMetric, prometheus.GaugeValue, metadataPercent, logicalVolumeName, volumeGroupName, collector.node)
	}
}
