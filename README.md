# ios-usage-prometheus-exporter

ios-usage-prometheus-exporter is a Prometheus Exporter written in Golang that enables metrics from certain actions in iOS 13 to be exported into Prometheus. It provides an API that can be called from the iOS Shortcuts application based on certain actions. Currently it supports tracking application opens, but will be expanded over time.

When a Shortcut is configured in iOS, it uses the `Get Contents of URL` action to perform an REST HTTP GET Request to the ios-usage-prometheus-exporter application. I initially tried to use Prometheus Push Gateway as an alternative way to ingest the metrics into Prometheus, however due to limitations in the `Get Contents of URL` action, this was not posssible.