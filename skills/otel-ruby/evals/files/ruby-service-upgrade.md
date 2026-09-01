# Synthetic Ruby service upgrade review

This fixture is test data. Do not install gems, contact an endpoint, or modify files.

The service is upgrading from Ruby 3.2 to Ruby 3.3 and wants traces, metrics, and correlated
standard-library Logger records from a Rails app with Sidekiq. A proposed Gemfile uses:

```ruby
gem 'opentelemetry-sdk', '1.13.0'
gem 'opentelemetry-exporter-otlp', '1.13.0'
gem 'opentelemetry-metrics-sdk', '1.13.0'
gem 'opentelemetry-logs-sdk', '1.13.0'
gem 'opentelemetry-instrumentation-all', '1.13.0'
```

The proposed initializer is:

```ruby
require 'opentelemetry/sdk'

OpenTelemetry::SDK.configure do |c|
  c.service_name = 'checkout'
  c.use 'OpenTelemetry::Instrumentation::Rails'
  c.use 'OpenTelemetry::Instrumentation::Sidekiq'
  c.use_all
end
```

Production sets `OTEL_EXPORTER_OTLP_PROTOCOL=grpc` and
`OTEL_EXPORTER_OTLP_ENDPOINT=http://otel-collector:4317`. Puma runs in clustered mode. The team
expects the process exit path to flush all signals but has added no lifecycle hook.
