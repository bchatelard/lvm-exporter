package main

import (
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// LVM collector, listen to port 9080 path /metrics
func main() {
	node, err := os.Hostname()
	if err != nil {
		node = "Unkown"
	}
	lvmVgCollector := newLvmVgCollector(node)
	prometheus.MustRegister(lvmVgCollector)

	lvmLvCollector := newLvmLvCollector(node)
	prometheus.MustRegister(lvmLvCollector)

	http.Handle("/metrics", promhttp.Handler())
	port := os.Getenv("PORT")
	if port == "" {
		port = ":9080"
	}
	log.Printf("Beginning to serve on port %v", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
