# ios-usage-prometheus-exporter

ios-usage-prometheus-exporter is a Prometheus Exporter, written in Golang that allows metrics from certain actions in iOS 13 to be exported into Prometheus. It centres around a HTTP RESTful API, that can be called using the new Shortcuts application available in iOS 13 based on certain triggers. Metrics are exported via the `/metrics` endpoint, that listens on a different port to the main API.

It operates in a very similar fashion to Prometheus Pushgateway, however, due to the fact that Pushgateway is only "vaguely REST-like", the iOS `Get Contents of URL` doesn't play nicely with it. It also has no authentication which makes it challenging to expose externally (which is essential for a roaming iOS device).

## What possible use case could you have for this?

I created this as a fun little side project, primarily because I was bored. But I also wanted to track my usage of social media and make a conscious effort to reduce the amount of times I was opening my social media apps. I use these metrics to graph my usage in Grafana.

## Supported Functions

Currently there is only one supported function type which is the `incTotalAppOpens` endpoint that tracks how many times an application is opened. This will be expanded with other useful metrics over time.

### `incTotalAppOpens` Endpoint

The `incTotalAppOpens` endpoint is used for incrementing a Prometheus Counter each time an application (or group of applications) are opened. Counter increments are performed by making a HTTP Get request via the iOS 13 `Get Contents of URL` action.

These are exported via the `ios_app_open_total` Prometheus Counter Metric. This metric has two labels:
1. `deviceName`: Derived from the `deviceName` HTTP header, used for aggregating stats across multiple devices.
2. `appName`: Derived from th `appName` HTTP header, used for identifying individual, or groups of applications.

### HTTP Request
`GET https://<YOUR-FQDN>/api/v1.0/incTotalAppOpens/`

### Headers

| Header      | Purpose                                                                                                                                             | Default   | Required  |
| ---         | ---                                                                                                                                                 | ---       | ---       |
| apiKey      | The API key that you defined in your configuration for authentication.                                                                              | None      | true      |
| deviceName  | The Device Name that the metrics are for. Used in the `deviceName` metric label. Alphanumeric only, special characters and spaces will be removed.  | None      | true      |
| appName     | The App Name (e.g Instagram) or group name (e.g. Social). Used in the `appName` metric label. Alphanumeric only, special characters and spaces will be removed.                                                                                                                                                                                          | None      | true      |

### Example Usage - Instagram

1. Create a new Personal Automation in iOS 13
2. Define a trigger of "Open App" and select Instagram
3. Add the `Get contents of URL` Action
4. Enter the FQDN of your instance of ios-usage-prometheus-exporter, followed by the `incTotalAppOpens` API endpoint. For example: `https://ios-metrics-reporter.maynard.io/api/v1.0/incTotalAppOpens/`
5. Add the 3 headers: `apiKey`, `deviceName` and `appName` and their associated values.
6. Save and ensure the "Ask Before Running" option is not selected.
7. Open Instagram
8. Validate success by running the `ios_app_open_total` Prometheus Query

## Deployment and Configuration

ios-usage-prometheus-exporter is designed to be deployed in Kubernetes, 

The API Webserver Port (defined by the `WEBSERVER_PORT` environment variable) should be exposed externally so that your iOS device can make API calls based on your defined triggers. It is strongly recommended that expose ios-usage-prometheus-exporter using HTTPS. SSL/TLS should be configured on your Load Balancer (or ingress if using something like ingress-nginx).

Service discovery in Kubernetes should be used for scraping the metrics, which are exposed via a separate port (defined by the `PROMETHEUS_PORT` environment variable). This should not be exposed externally. 

The following environment variables are used by ios-usage-prometheus-exporter for configuration:

| Environment Variable  | Purpose                                                                                                                       | Default   | Required  |
| ---                   | ---                                                                                                                           | ---       | ---       |
| API_KEY               | A secure API Key that will be validated against the `apiKey` header in requests. **This should be created as a K8s secret.**  | None      | true      |
| PROMETHEUS_PORT       | The port to use for the webserver that exposes the `/metrics` Prometheus endpoint.                                            | 9090      | false     |
| WEBSERVER_PORT        | The port to use for the webserver that exposes the ios-usage-prometheus-exporter REST API                                     | 80        | false     |