package main

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

const (
	request200 = "https://httpstat.us/200"
	request503 = "https://httpstat.us/503"
	invalidURL = "https://tushartest:8080"
)

var (
	sampleExternalURLUp = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "sample_external_url_up",
	}, []string{"url"},
	)
	sampleExternalURLResponseMs = prometheus.NewSummaryVec(prometheus.SummaryOpts{
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
	responseTime := int64(time.Since(start).Milliseconds())
	statusCode := resp.StatusCode
	if statusCode == 503 {
		sampleExternalURLUp.With(prometheus.Labels{"url": url})
	} else {
		sampleExternalURLUp.With(prometheus.Labels{"url": url}).Inc()
	}
	sampleExternalURLResponseMs.With(prometheus.Labels{"url": url}).Observe(float64(responseTime))
	logrus.Infof("URL: [%s], Response time:[%d], Status code:[%d]", url, responseTime, statusCode)
}

func main() {
	logrus.Infof("internal-service-monitor started..")
	logrus.Infof("200 success request:")
	getURLDetails(request200)
	logrus.Infof("503 service unreachable request:")
	getURLDetails(request503)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":10001", nil)
}
