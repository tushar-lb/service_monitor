# INTERNAL SERVICE MONITOR

Monitor internal request urls and expose prometheus metrics and build visualizations, dashboards.

## Structure of data:
```
sample_external_url_response_ms{url="https://httpstat.us/200"} 1181
sample_external_url_response_ms{url="https://httpstat.us/503"} 335

sample_external_url_up{url="https://httpstat.us/200"} 1
sample_external_url_up{url="https://httpstat.us/503"} 0
```

## Prerequisites:
- go version 1.14 and above

## Build Instructions

1. Clone repository and compile source code:

```console
$ git clone https://github.com/tusharraut1994/service_monitor.git
$ cd service_monitor
$ make
```

2. Build and push docker image:

```console
$ make container deploy
```

## Deploy Instructions - Existing Prometheus and Grafana

1. Create internal-service-monitor deployment and service on kubernetes
```console
$ kubectl apply -f specs/internal_service_monitor.yaml --namespace test
```

2. Create prometheus service monitor
```console
$ kubectl apply -f specs/prometheus_service_monitor.yaml --namespace test
```

3. Update the prometheus instance and add serviceMonitorSelector label.
ex:
```
serviceMonitorSelector:
  matchExpressions:
  - key: prometheus
    operator: In
    values:
    - internal-service-monitor
```

## Deploy Instructions - Create Prometheus and Grafana

NOTE: Used `test` namespace in all specs.

1. Create internal-service-monitor deployment and service on kubernetes
```console
$ kubectl apply -f specs/internal_service_monitor.yaml
```

2. Deploy prometheus
```console
$ kubectl apply -f specs/prometheus-operator.yaml
$ kubectl apply -f specs/prometheus-instance.yaml
```

3. Deploy Grafana
```console
$ kubectl apply -f specs/grafana.yaml
```

4. Create prometheus service monitor
```console
$ kubectl apply -f specs/prometheus_service_monitor.yaml
```

5. Verify all components are up:
```
[root@tushar-dev ~]# kubectl get po,svc,prometheus,servicemonitor --namespace test
NAME                                                   READY   STATUS    RESTARTS   AGE
pod/internal-service-monitor-7f67fbd579-2k8t2          1/1     Running   0          8m4s
pod/prometheus-internal-service-monitor-prometheus-0   3/3     Running   0          7m40s
pod/prometheus-internal-service-monitor-prometheus-1   3/3     Running   0          7m40s
pod/prometheus-operator-c7d946cd7-4x7g7                1/1     Running   0          7m56s
pod/test-grafana-65c7ccdcd7-vl2ln                      1/1     Running   0          8m8s

NAME                               TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
service/internal-service-monitor   ClusterIP   10.108.37.3      <none>        10001/TCP        8m6s
service/prometheus                 NodePort    10.107.119.242   <none>        9090:31052/TCP   8m
service/prometheus-operated        ClusterIP   None             <none>        9090/TCP         7m40s
service/test-grafana               NodePort    10.110.61.214    <none>        3000:30728/TCP   8m7s

NAME                                                                   VERSION   REPLICAS   AGE
prometheus.monitoring.coreos.com/internal-service-monitor-prometheus             2          8m1s

NAME                                                               AGE
servicemonitor.monitoring.coreos.com/internal-service-monitor-sm   7m55s
```

6. Prometheus UI will be accessible on: `http://NODE_IP:NODE_PORT` nodePort of `prometheus` service.

7. Login to Grafana UI using: `http://NODE_IP:NODE_PORT` nodePort of `grafana` service. Login with default credentials: `admin/admin`. After successful login upload test dashboard json template.

## Sample Grafana Dashborads & Prometheus Queries Available Here:
https://github.com/tusharraut1994/service_monitor/blob/master/dashboards.md