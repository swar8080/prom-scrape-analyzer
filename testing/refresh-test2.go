package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	registry     = prometheus.NewRegistry()
	metrics      = make(map[string]*prometheus.GaugeVec)
	metricsMutex sync.Mutex
)

func createNewMetric(name string) {
	metricsMutex.Lock()
	defer metricsMutex.Unlock()

	g := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: name,
		Help: "Dynamically generated metric",
	}, []string{"label"})

	metrics[name] = g
	registry.MustRegister(g)
	fmt.Println("Added new metric:", name)
}

func increaseCardinality() {
	metricsMutex.Lock()
	defer metricsMutex.Unlock()

	for _, g := range metrics {
		label := fmt.Sprintf("val_%d", rand.Intn(100000))
		g.WithLabelValues(label).Set(rand.Float64() * 100)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	baseMetricNames := []string{"metric_one", "metric_two", "metric_three", "metric_four", "metric_five"}

	// Initialize 5 base metrics
	for _, name := range baseMetricNames {
		createNewMetric(name)
	}

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	go func() {
		for {
			increaseCardinality()
			time.Sleep(time.Duration(rand.Intn(2500)+500) * time.Millisecond)
		}
	}()

	go func() {
		counter := 1
		for {
			createNewMetric(fmt.Sprintf("dynamic_metric_%d", counter))
			counter++
			time.Sleep(time.Duration(rand.Intn(1000)+1000) * time.Millisecond)
		}
	}()

	//go func() {
	//	for {
	//		time.Sleep(6 * time.Second)
	//		metricsMutex.Lock()
	//		registry := prometheus.NewRegistry()
	//		metrics = make(map[string]*prometheus.GaugeVec)
	//		for _, name := range baseMetricNames {
	//			createNewMetric(name)
	//		}
	//
	//		http.DefaultServeMux = http.NewServeMux()
	//		http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	//
	//		metricsMutex.Unlock()
	//		fmt.Println("Registry cleared and reset")
	//	}
	//}()

	fmt.Println("Serving metrics on :8080/metrics")
	http.ListenAndServe(":8080", nil)
}
