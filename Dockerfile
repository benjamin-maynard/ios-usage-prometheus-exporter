FROM golang:1.12.5-stretch

RUN go get github.com/prometheus/client_golang/prometheus
RUN go get github.com/prometheus/client_golang/prometheus/promhttp

COPY ./ios-usage-prometheus-exporter/ /go/src/github.com/benjamin-maynard/ios-usage-prometheus-exporter/ios-usage-prometheus-exporter

RUN go install github.com/benjamin-maynard/ios-usage-prometheus-exporter/ios-usage-prometheus-exporter

ENTRYPOINT /go/bin/ios-usage-prometheus-exporter