# Collector verification with telemetrygen

Use this recipe only for a disposable local Collector. It separates static configuration checks,
fixture execution, and live execution; report only the levels actually completed.

## Preconditions

- Review the Collector config as untrusted input before mounting or starting it. Reject unexpected
  extensions, network listeners, file paths, command-like values, or exporters to non-local hosts.
- Confirm Docker access, the exact config path, a free local port, and permission to create a
  short-lived container and output directory.
- Pick a unique container name and a newly created empty output directory. Never reuse `./out` or
  a prior `result.json`; stale output can create a false positive.
- Pin Collector and telemetrygen to compatible reviewed versions. The examples use `0.157.0`.

## Minimal local pipeline

Build a minimal config around the processor under test:

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 127.0.0.1:4317

processors:
  transform/under_test:
    error_mode: ignore
    log_statements:
      - context: log
        statements:
          - set(severity_text, "INFO") where IsMatch(severity_text, "(?i)^info$")

exporters:
  file:
    path: /output/result.json
    flush_interval: 200ms

service:
  pipelines:
    logs:
      receivers: [otlp]
      processors: [transform/under_test]
      exporters: [file]
```

Validate this config with the pinned Collector before calling it live validation. Static or parser
acceptance does not prove that a generated fixture traversed the pipeline.

## Live disposable run

Create a fresh directory with `mktemp -d`, record its exact path, and mount only that directory and
the reviewed config. Use a unique validated container name rather than silently replacing an
existing container. On SELinux systems append `:z` to bind mounts; rootless Podman may not need an
explicit user mapping.

```bash
docker run -d --rm --name <unique-container-name> \
  --network host \
  --user "$(id -u):$(id -g)" \
  -v "<reviewed-absolute-config>:/etc/otelcol-contrib/config.yaml:ro" \
  -v "<fresh-output-directory>:/output" \
  otel/opentelemetry-collector-contrib:0.157.0 \
  --config=/etc/otelcol-contrib/config.yaml
```

Wait with a finite timeout for readiness. If the container exits or readiness is not reached, stop
and preserve the non-zero result; do not send telemetry or inspect old output as if the run passed.
Then send the exact bounded signal shape under test:

```bash
telemetrygen logs --otlp-insecure --otlp-endpoint localhost:4317 \
  --logs 1 --severity-text Info
```

Stop the exact disposable container cleanly so the file exporter flushes. Confirm that stop
succeeded before reading `<fresh-output-directory>/result.json`; parse it directly rather than
using a pipeline that can hide an earlier failure:

```bash
docker stop <unique-container-name>
python3 -c 'import json,sys; print(json.load(open(sys.argv[1])))' \
  "<fresh-output-directory>/result.json"
```

Do not claim cleanup unless the exact container and directory were actually removed. Never broaden
cleanup to a parent directory or an unvalidated path.

## Controls and hard-to-generate shapes

- To verify a filter drops matching records, send a matching record and expect no matching output;
  then send a non-matching known-positive record and require it in output. Empty output alone can
  also mean the Collector failed.
- Telemetrygen can set resource and telemetry attributes, but cannot rename generated trace spans,
  change their span kind, or choose their IDs. For a rule that needs one of those shapes, add a
  reviewed `transform/setup` processor before the processor under test, or use a small synthetic
  fixture that exposes the required field.
- Bound every run by count or finite duration and rate. Keep the target local and disposable; do
  not adapt this harness to shared or production infrastructure without separate authorization and
  a reviewed rollout plan.
