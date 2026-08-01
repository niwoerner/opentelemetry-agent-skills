# Synthetic release evidence for evaluation

This fixture is a minimal, source-labeled simulation of facts that can disagree in upstream
retrieval. It is test data, not current compatibility guidance.

## Generic Language Support Status excerpt

- Language: Java
- Latest supported file format: `1.0.0-rc.3`
- Meaning: schema coverage metadata

## Javaagent 2.30.0 released parser and smoke-fixture excerpt

- Parser accepts schema family `1.*` and prefers `1.1`.
- Released smoke fixture begins with `file_format: "1.1"`.
- Declarative duration values are milliseconds.

## Configuration schema release v1.1.0 excerpt

- Canonical released example begins with `file_format: "1.1"`.
- OTLP/HTTP exporter endpoints are signal endpoints. For a base gateway at
  `http://otel-gateway:4318`, use `/v1/traces`, `/v1/metrics`, and `/v1/logs` for the respective
  exporters unless the selected runtime documentation proves it appends paths itself.
