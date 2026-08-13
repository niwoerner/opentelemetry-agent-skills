# OpenTelemetry Go upgrade review

The service is upgrading from core OpenTelemetry Go v1.44.0 and log API/SDK v0.20.0 to
core v1.45.0 and log API/SDK v0.21.0. It currently emits a log record like this:

```go
var rec log.Record
rec.SetBody(log.StringValue("checkout failed"))
rec.AddAttributes(log.String("component", "checkout"))
logger.Emit(ctx, rec)
```

Its trace exporter is configured as follows, and the Collector expects traces at
`http://otel-gateway:4318/v1/traces`:

```go
exporter, err := otlptracehttp.New(ctx,
    otlptracehttp.WithEndpointURL("http://otel-gateway:4318"),
)
```

The Collector may return either `Retry-After: 5` or an HTTP-date value while throttling.
The team needs to know whether the v1.45.0 exporter honors both forms and how that differs
from v1.44.0. This is a static review only: do not contact the example endpoint.
