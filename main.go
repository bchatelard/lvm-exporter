package main

import (
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	klog "k8s.io/klog/v2"
)

// LVM collector, listen to port 9080 path /metrics
func main() {
	node := os.Getenv("KUBE_NODE_NAME")
	if len(node) == 0 {
		var err error
		node, err = os.Hostname()
		if err != nil {
			node = "Unkown"
		}
	}
	lvmVgCollector := newLvmVgCollector(node)
	prometheus.MustRegister(lvmVgCollector)

	lvmLvCollector := newLvmLvCollector(node)
	prometheus.MustRegister(lvmLvCollector)

	http.Handle("/metrics", promhttp.Handler())
	klog.Info("Beginning to serve on port :9080")
	klog.Fatal(http.ListenAndServe(":9080", nil))
}
