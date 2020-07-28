package main

import (
	"reflect"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestURLs(t *testing.T) {

	URLsList := getURLsList()
	if len(URLsList) == 0 {
		t.Errorf("No urls specified for status check")
	}

	if URLsList[0] != "https://httpstat.us/200" {
		t.Errorf("[https://httpstat.us/200] should be present in input urls list")
	}

	if URLsList[1] != "https://httpstat.us/503" {
		t.Errorf("[https://httpstat.us/503] should be present in input urls list")
	}
}

func TestURLsStatusCode(t *testing.T) {
	URLsList := getURLsList()
	if len(URLsList) == 0 {
		t.Errorf("No urls specified for status check")
	}

	status, _, _ := getURLDetails(URLsList[0])
	if status != 200 {
		t.Errorf("URL : [https://httpstat.us/200] status should be : 200")
	}

	status, _, _ = getURLDetails(URLsList[1])
	if status != 503 {
		t.Errorf("URL : [https://httpstat.us/503] status should be : 503")
	}
}

func TestURLsUPStatus(t *testing.T) {
	URLsList := getURLsList()
	if len(URLsList) == 0 {
		t.Errorf("No urls specified for status check")
	}

	_, status, _ := getURLDetails(URLsList[0])
	if status != 1 {
		t.Errorf("URL : [https://httpstat.us/200] UP status should be : 1")
	}

	_, status, _ = getURLDetails(URLsList[1])
	if status != 0 {
		t.Errorf("URL : [https://httpstat.us/503] UP status should be : 0")
	}
}

func TestURLsResponseTimeInMS(t *testing.T) {
	URLsList := getURLsList()
	if len(URLsList) == 0 {
		t.Errorf("No urls specified for status check")
	}

	_, _, responseTime := getURLDetails(URLsList[0])
	if responseTime <= 0 {
		t.Errorf("URL : [https://httpstat.us/200] response time should be > 0")
	}

	_, _, responseTime = getURLDetails(URLsList[1])
	if responseTime <= 0 {
		t.Errorf("URL : [https://httpstat.us/503] response time should be > 0")
	}
}

func TestSampleQueryVariables(t *testing.T) {
	var testsampleExternalURLUp *prometheus.GaugeVec
	var testsampleExternalURLResponseMs *prometheus.GaugeVec

	if reflect.TypeOf(sampleExternalURLUp) != reflect.TypeOf(testsampleExternalURLUp) {
		t.Errorf("sample_external_url_up is not defined")
	}

	if reflect.TypeOf(sampleExternalURLResponseMs) != reflect.TypeOf(testsampleExternalURLResponseMs) {
		t.Errorf("sample_external_url_response_ms is not defined")
	}
}

func TestSampleQueryName(t *testing.T) {
	if externalURLResponseQueryName != "sample_external_url_response_ms" {
		t.Errorf("External URL response time query name should be set to: sample_external_url_response_ms")
	}

	if externalURLUPQueryName != "sample_external_url_up" {
		t.Errorf("External URL UP status query name should be set to: sample_external_url_up")
	}
}

func TestDefaultMetricsHandler(t *testing.T) {
	defaultMetricsHandler := "/metrics"
	if defaultMetricsHandler != metricsHandler {
		t.Errorf("Default metrics handler is /metrics")
	}
}

func TestDefaultServiceMonitorListenPort(t *testing.T) {
	defaultListenPort := 10001
	if defaultListenPort != serviceMonitorPort {
		t.Errorf("Default listener port is not set to 10001")
	}
}
