---
name: otel-ottl
description: OpenTelemetry Transformation Language (OTTL) expert for writing and debugging telemetry transformations in the OpenTelemetry Collector. Use when authoring or reviewing `transform`, `filter`, `tail_sampling` processor configs or `routing` connector configs, debugging OTTL syntax or semantics, transforming traces, metrics, logs, or profiles, or converting data-processing requirements into OTTL statements.
---

# OpenTelemetry Transformation Language (OTTL)

OTTL transforms or selects telemetry inside Collector components. This skill is pinned to
collector-contrib **v0.157.0**. Function, path, default, and feature-gate availability varies by
release; when the user's version differs, verify against the matching upstream tag.

## Workflow

1. **Choose the component.** `transform` rewrites, `filter` drops, `tail_sampling` decides whether
   to retain traces, and `routing` sends telemetry to pipelines. A component controls its available
   contexts and functions.
2. **Choose the lowest usable context.** Lower contexts can read their parents (for example, a span
   can read `resource.attributes`), but parents cannot read children. Use `datapoint` for point
   attributes instead of traversing `metric.data_points`.
3. **Write the statement.** An editor such as `set` or `delete_key` mutates data and may have a
   `where` condition. Converters such as `ParseJSON` and `IsMatch` return values; they do not mutate.
4. **Set error behavior deliberately.** `ignore` logs statement errors and continues; `silent`
   continues without logging; `propagate` returns the error and can cause the component to drop the
   payload. In v0.157, transform and filter default to `ignore`; routing also defaults to `ignore`
   while its beta default-error feature gate is enabled. For routing, `ignore` sends an errored
   payload to `default_pipelines`; configure that fallback or the payload is dropped.
5. **Verify end to end.** Validate the exact Collector version, then send known telemetry and inspect
   file-exporter output. Use the
   [telemetrygen recipe](../otel-telemetrygen/SKILL.md#verifying-a-collector-config).

```ottl
set(span.attributes["env"], "prod") where resource.attributes["env"] == nil
```

## Load only what the task needs

- [Contexts](references/contexts.md) — exact paths, hierarchy, enums, and request metadata.
- [Functions](references/functions.md) — editor/converter signatures and release availability.
- [Quick reference](references/quick-reference.md) — component YAML, recipes, escaping,
  troubleshooting, and safe skeletons.

For a single path or function, read only the relevant section instead of loading the full catalogs.

## Safety and correctness gates

- Guard optional or polymorphic input before conversion: `where x != nil`, `IsString(x)`, or the
  appropriate type check.
- For JSON-object-only work, guard both the type and shape before calling `ParseJSON`, for example
  `IsString(log.body) and IsMatch(log.body.string, "(?s)^\\s*\\{.*\\}\\s*$")`. The RE2 `(?s)`
  flag admits pretty-printed objects containing newlines. Checking `IsMap` after
  parsing does not prevent arrays or scalar JSON from being parsed.
- On a version-pinned request, confirm every chosen path and function against that release tag;
  do not assume a function listed for this skill's v0.157 anchor exists in an older release.
  For v0.156 JSON-object parsing, `ParseJSON`, `IsString`, and `IsMatch` are available without the
  v0.157 alpha lambda feature gate.
- Request metadata is read-only and may contain credentials. Copy only explicitly allowlisted,
  non-sensitive keys. OTLP metadata routing requires `include_metadata: true` on the receiver.
  HTTP/client header spelling may retain its form (`otelcol.client.metadata["X-Tenant"][0]`);
  gRPC metadata keys are lowercase (`otelcol.grpc.metadata["x-tenant"][0]`).
- The routing `request` context is deprecated as of v0.156; use `otelcol.client.metadata` or
  `otelcol.grpc.metadata`.
- Log-record-specific rewrites of shared resource or scope data require `flatten_data: true` and the
  alpha `transform.flatten.logs` gate. This copies and regroups data; do not enable it accidentally.
- Hashing an identifier does not necessarily anonymize it. Apply the organization's data-handling
  policy before retaining deterministic hashes of personal data.

## Frequent syntax traps

- In Collector YAML, write an OTTL replacement backreference `${1}` as `$${1}`. A replacement such
  as `$1REDACTED` is literal and silently fails to substitute the capture.
- Go RE2 rejects large counted repetitions such as `(.{1024}).*`; use `Substring` with a nil/type
  guard and `Len`, or `truncate_all` for a map.
- Current span-event paths use `spanevent.*`, not `span_event.*`. Cache paths are context-qualified,
  such as `span.cache["parsed"]`.
- Use `Decode(value, "base64")`; `Base64Decode` is deprecated.
- Regex escapes inside OTTL strings are doubled (`\\d`, `\\s`, `\\.`).

## Upstream sources

- [OTTL package](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/v0.157.0/pkg/ottl)
- [Transform processor](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/v0.157.0/processor/transformprocessor)
- [Filter processor](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/v0.157.0/processor/filterprocessor)
- [Routing connector](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/v0.157.0/connector/routingconnector)
