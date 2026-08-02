# Synthetic runtime-upgrade evidence

This fixture is test data. Version and validation records are synthetic and must not be treated as
current OpenTelemetry compatibility guidance.

## Target deployment identity

```text
runtime: Java 21
agent package: example.synthetic/otel-javaagent
agent version: 2.31.0
bootstrap: declarative loader enabled
identified official OpenTelemetry source repository: open-telemetry/opentelemetry-java-instrumentation
startup diagnostic: configuration rejected after upgrade; schema mismatch at file_format "1.1"
```

## Retained compatibility cache record

```text
runtime: Java 21
agent package: example.synthetic/otel-javaagent
agent version: 2.30.0
selected schema tag: v1.1.0
runtime source revision: refs/tags/v2.30.0
accepted file_format: "1.1"
release-schema validation: passed for the cached file
selected-parser validation: passed with agent 2.30.0
live validation: not run
```

## Schema release discovery output

```text
v1.3.0 (latest discovered release)
v1.2.0
v1.1.0
```

The discovery output contains no mapping from agent 2.31.0 to a compatible schema release or
accepted `file_format` literal.
