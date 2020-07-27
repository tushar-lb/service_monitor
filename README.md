# INTERNAL SERVICE MONITOR

Monitor internal request urls and expose prometheus metrics and build visualizations, dashboards.

## Structure of data:
```
sample_external_url_response_ms{url="https://httpstat.us/200"} 1181
sample_external_url_response_ms{url="https://httpstat.us/503"} 335

sample_external_url_up{url="https://httpstat.us/200"} 1
sample_external_url_up{url="https://httpstat.us/503"} 0
```

## Compile source code:

```console
$ make
```

## Build and push docker image:

```console
$ make container deploy
```

## Deploy internal-service-monitor deployment and service on kubernetes
```console
$ kubectl apply -f specs/internal_service_monitor.yaml
```

## Create prometheus service monitor
```console
$ kubectl apply -f specs/prometheus_service_monitor.yaml
```

