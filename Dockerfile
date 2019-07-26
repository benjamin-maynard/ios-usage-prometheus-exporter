FROM golang:1.12.5-stretch

RUN go get github.com/prometheus/client_golang/prometheus
RUN go get github.com/prometheus/client_golang/prometheus/promhttp

COPY ./ios-usage-prometheus-exporter/ /go/src/github.com/benjamin-maynard/ios-usage-prometheus-exporter/ios-usage-prometheus-exporter

RUN go install github.com/benjamin-maynard/ios-usage-prometheus-exporter/ios-usage-prometheus-exporter

ENV PROMETHEUS_PORT=9090
ENV WEBSERVER_PORT=80

ENTRYPOINT /go/bin/ios-usage-prometheus-exporter