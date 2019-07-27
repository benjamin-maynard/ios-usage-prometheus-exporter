# ios-usage-prometheus-exporter

ios-usage-prometheus-exporter is a Prometheus Exporter, written in Golang that allows metrics from certain actions in iOS 13 to be exported into Prometheus. It centers around a HTTP RESTful API, that can be called using the new Shortcuts application available in iOS 13 based on certain triggers. Metrics are exported via the `/metrics` endpoint, that listens on a different port to the main API.

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