# ios-usage-prometheus-exporter

ios-usage-prometheus-exporter is a Prometheus Exporter, written in Golang that allows metrics from certain actions in iOS 13 to be exported into Prometheus. It centers around a HTTP RESTful API, that can be called using the new Shortcuts application available in iOS 13 based on certain triggers. Metrics are exported via the `/metrics` endpoint, that listens on a different port to the main API.

## Deployment and Configuration

ios-usage-prometheus-exporter is designed to be deployed in Kubernetes, 

The API Webserver Port (defined by the `WEBSERVER_PORT` environment variable) should be exposed externally so that your iOS device can make API calls based on your defined triggers. It is strongly recommended that expose ios-usage-prometheus-exporter using HTTPS. SSL/TLS should be configured on your Load Balancer (or ingress if using something like ingress-nginx).

Service discovery in Kubernetes should be used for scraping the metrics, which are exposed via a seperate port (defined by the `PROMETHEUS_PORT` environment variable). This should not be exposed externally. 

The following environment variables are used by ios-usage-prometheus-exporter for configuration:

| Environment Variable  | Purpose                                                                                                                       | Default   | Required  |
| ---                   | ---                                                                                                                           | ---       | ---       |
| API_KEY               | A secure API Key that will be validated against the `apiKey` header in requests. **This should be created as a K8s secret.**  | None      | true      |
| PROMETHEUS_PORT       | The port to use for the webserver that exposes the `/metrics` Prometheus endpoint.                                            | 9090      | false     |
| WEBSERVER_PORT        | The port to use for the webserver that exposes the ios-usage-prometheus-exporter REST API                                     | 80        | false     |


## Example Usage - Measuring App Usage

The new iOS 13 Shortcuts application enables users to define Personal Automations, which are actions that are performed on the back of certain activities, for example opening the Instagram App.

This application works by using the `Get Contents of URL` action in iOS 13 to perform API calls when these actions are performed.

For example - to increment a Prometheus Counter each time the Instagram application is opened, you would do the following:

1. Create a new Personal Automation in iOS 13
2. Define a trigger of "Open App" and select Instagram
3. Add the `Get contents of URL` Trigger
4. Enter the FQDN of your instance of ios-usage-prometheus-exporter, followed by the `incTotalAppOpens` API endpoint. For example: `https://ios-metrics-reporter.maynard.io/api/v1.0/incTotalAppOpens/`
5. Add 3 headers: `apiKey`, `deviceName` and `appName` and their associated values.
6. Save and ensure the "Ask Before Running" option is not selected.