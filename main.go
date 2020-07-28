package main

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var (
	sampleExternalURLUp = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "sample_external_url_up",
	}, []string{"url"},
	)
	sampleExternalURLResponseMs = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "sample_external_url_response_ms",
	}, []string{"url"},
	)
)

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(sampleExternalURLUp)
	prometheus.MustRegister(sampleExternalURLResponseMs)
}

func getURLDetails(url string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		logrus.Errorf("error while checking url [%s] status, %v", url, err)
	}
	responseTime := float64(time.Since(start).Milliseconds())
	statusCode := resp.StatusCode
	logrus.Infof("URL: [%s], Response time:[%d], Status code:[%d]", url, responseTime, statusCode)
	if statusCode == 503 {
		sampleExternalURLUp.With(prometheus.Labels{"url": url})
	} else {
		sampleExternalURLUp.With(prometheus.Labels{"url": url}).Inc()
	}
	sampleExternalURLResponseMs.With(prometheus.Labels{"url": url}).Add(responseTime)
}

func main() {
	logrus.Infof("internal-service-monitor started..")
	// Define list of url's to check
	URLsList := []string{"https://httpstat.us/200", "https://httpstat.us/503"}
	logrus.Infof("URLs list: %s", URLsList)
	for _, url := range URLsList {
		getURLDetails(url)
	}
	// register /metrics http handler with promhttp
	http.Handle("/metrics", promhttp.Handler())
	// start listening and serving requests on port 100001
	http.ListenAndServe(":10001", nil)
}
