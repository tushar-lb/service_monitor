package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

const (
	externalURLUPQueryName       = "sample_external_url_up"
	externalURLResponseQueryName = "sample_external_url_response_ms"
	serviceMonitorPort           = 10001
	metricsHandler               = "/metrics"
	interval                     = 5 * time.Second
	serviceURLOne                = "https://httpstat.us/200"
	serviceURLtwo                = "https://httpstat.us/503"
)

var (
	sampleExternalURLUp = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: externalURLUPQueryName,
	}, []string{"url"},
	)
	sampleExternalURLResponseMs = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: externalURLResponseQueryName,
	}, []string{"url"},
	)
)

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(sampleExternalURLUp)
	prometheus.MustRegister(sampleExternalURLResponseMs)
}

func getURLDetails(url string) (int, int, float64) {
	setURLUPStatus := 0
	statusCode := 0
	responseTime := 0.0

	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		logrus.Errorf("error while checking url [%s] status, %v", url, err)
		return 0, 0, 0
	}

	responseTime = float64(time.Since(start).Milliseconds())
	statusCode = resp.StatusCode

	if statusCode == 200 {
		setURLUPStatus = 1
	}
	sampleExternalURLUp.With(prometheus.Labels{"url": url}).Set(float64(setURLUPStatus))
	sampleExternalURLResponseMs.With(prometheus.Labels{"url": url}).Set(responseTime)

	logrus.Infof("URL: [%s], Response time:[%f], Status code:[%d], URL up status set to: [%d]", url, responseTime, statusCode, setURLUPStatus)
	return statusCode, setURLUPStatus, responseTime
}

func getURLsList() []string {
	URLsList := make([]string, 2)
	URLsList[0] = serviceURLOne
	URLsList[1] = serviceURLtwo
	logrus.Infof("URLs list: %s", URLsList)
	return URLsList
}

func startListener() {
	mux := http.NewServeMux()
	mux.HandleFunc("/status", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprint(res, "running")
	})
	// register /metrics http handler with promhttp
	mux.Handle(metricsHandler, promhttp.Handler())
	// start listening and serving requests on port 10001
	listenAddress := fmt.Sprintf(":%d", serviceMonitorPort)
	http.ListenAndServe(listenAddress, mux)
}

func main() {
	logrus.Infof("internal-service-monitor started..")
	// Start /metrics handler and listen on port 10001
	go startListener()
	// Get list of urls
	for {
		URLsList := getURLsList()
		if len(URLsList) > 0 {
			for _, url := range URLsList {
				statusCode, setURLUPStatus, responseTime := getURLDetails(url)
				logrus.Debugf("URL: [%s], Response time:[%f], Status code:[%d], URL up status set to: [%d]", url, responseTime, statusCode, setURLUPStatus)
			}
		} else {
			panic("No urls available for status check")
		}
		time.Sleep(interval)
	}
}
