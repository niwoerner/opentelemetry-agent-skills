---
name: opentelemetry-manual-instrumentation
description: OpenTelemetry best practices for manual instrumentation across any language. Use when planning, adding, or reviewing instrumentation, choosing signals and semantic conventions, naming spans, metrics, and attributes, controlling cardinality, propagating context, or checking released SDK and semantic convention versions.
---

# OpenTelemetry Manual Instrumentation

Use this skill for application-level instrumentation design and review.

## Workflow

0. Prepare before planning or editing code.
- check `references/generated/otel-version-index.md` to identify the latest available SDK or package version for the project language
- use the latest version when it is compatible with the project
- if the latest version is not compatible, use the latest compatible version and note the compatibility reason
- perform this version lookup once during preparation and reuse the result unless the target language, package, or compatibility constraints change
- when implementation details, examples, or SDK behavior are unclear, check the official SDK documentation and source code instead of relying only on model memory

1. Define the instrumentation plan before editing code.
- identify the runtime boundaries to instrument
- choose the signal for each boundary
- propose names and key attributes
- note propagation requirements and cardinality risks

2. Find the runtime boundary.
Create spans at meaningful boundaries such as incoming requests, outgoing service calls, database calls, message publish or consume, cache or network interactions, and high-value business operations.

Do not create spans for helpers, loop iterations, getters, validation-only steps, or pure computation.

See `references/boundaries.md`.

3. Choose the signal before writing code.
- span: request or operation flow
- metric: counts, rates, latency distributions, utilization, saturation
- log or event: discrete diagnostic fact
- nothing: if the telemetry does not improve production understanding

See `references/signal-selection.md`.

4. Always prefer semantic conventions before inventing custom keys.

For any known boundary type such as `http`, `db`, `messaging`, `rpc`, or another released semantic convention group:
- query the released semantic convention group before choosing names or attributes
- do this once per boundary type in the change, not before every edit
- use the released naming and attribute guidance from that group in the implementation
- derive span names from the released semantic convention naming rule only
- do not prepend or append custom descriptors such as protocol labels, hostnames, operation summaries, product names, or business hints to semconv-governed span names
- put extra context into released attributes, custom attributes, events, or parent-child structure instead of the span name
- if the released semantic convention naming rule does not provide a low-cardinality target, use the simpler fallback allowed by that convention instead of inventing a custom name
- if the code does not match the released semantic conventions for the boundary, fix it before finishing unless a concrete compatibility limitation prevents it

Use the bundled semantic convention lookup workflow:
- list groups: `./scripts/query-otel-semantic-conventions.sh --groups`
- inspect one group: `./scripts/query-otel-semantic-conventions.sh http`
- inspect one attribute: `./scripts/query-otel-semantic-conventions.sh http http.request.method`

See `references/semconv-selection.md` and `references/otel-semantic-conventions.md`. Adhere strictly to the conventions for names and attributes.

5. Configure or reuse the SDK, then implement or update the code.

For SDK setup, prefer these defaults unless the project already requires something else:
- traces: OTLP exporter plus batch span processor
- metrics: OTLP exporter plus periodic exporting metric reader
- logs: OTLP exporter plus batching log record processor
- propagators: `tracecontext,baggage`
- protocol: prefer the SDK default transport; if choosing explicitly, prefer `http/protobuf` unless the SDK or project requires `grpc`


During implementation:
- if no SDK is set up yet, configure one
- if an SDK is already present, reuse and extend the existing setup instead of replacing it
- when setting up or updating the SDK, use the common default production pipeline unless the project already has an intentional alternative
- default to OTLP exporters for enabled signals
- use a batch span processor for traces
- use a periodic exporting metric reader for metrics
- use a batching log record processor for logs when logs are part of the SDK setup
- preserve existing project-specific exporter, processor, or transport choices when they are already intentional and compatible
- control cardinality
- assign error ownership to the span that owns the final failed outcome
- preserve context across network and async boundaries

6. Control cardinality.
Prefer bounded values such as method, route template, status code, operation name, destination name, region, deployment environment, or customer tier.

Always avoid unbounded values such as raw user IDs, full URLs, raw SQL text in metrics, free-form messages, and timestamps.

7. Handle errors by ownership and final outcome.
- record a failure on the span that owns the final failed outcome
- do not mark the final span as failed if retries succeed
- do not treat handled errors as terminal failures

See `references/boundaries.md`.

8. Preserve context across boundaries.
At every network or async boundary:
- extract incoming context
- inject outgoing context
- use baggage only when intentional and bounded

See `references/propagation.md`.

9. Run the required completion loop before finishing.

## Required Completion Loop

Follow this loop every time:
1. define the instrumentation plan before editing code
2. implement or update the instrumentation
3. review the changed code against the checklist below
4. re-open the changed files and confirm each checklist item with codebase evidence
5. if any item is unresolved, patch the code or mark it not applicable with a reason, then repeat the review
6. do not finish until every checklist item is completed or explicitly marked not applicable

Do not mark a checklist item complete based on intent alone. Mark it complete only after confirming it in the current codebase.

## Instrumentation Checklist

For every item, report one of these statuses in the final answer:
- `[x]` completed
- `[~]` not applicable, with a reason
- `[ ]` unresolved

Include file references as evidence for every completed item.
If you cannot cite codebase evidence for an item, leave it unresolved and continue the loop.

- `[ ]` Instrumentation is attached to a meaningful runtime boundary, not a helper, loop body, or other low-value location.
- `[ ]` Preparation used the version index to choose the latest compatible SDK or package version, and official SDK docs or source were checked when needed.
- `[ ]` SDK setup uses the common default production pipeline for enabled signals, or preserves an intentional existing project-specific alternative.
- `[ ]` Released semantic convention groups were queried for each known boundary type in the change.
- `[ ]` Implemented names and attributes match the released semantic conventions for each known boundary type, or any deviation is explained by a concrete compatibility limitation.
- `[ ]` For each known boundary type, every semconv-governed span name can be explained directly from the released naming rule.
- `[ ]` Semconv-governed span names do not include custom prose, protocol prefixes, hostnames, or business labels.
- `[ ]` The signal choice is intentional for each instrumentation point: span, metric, event or log, or nothing.
- `[ ]` Attributes use bounded cardinality and avoid raw high-cardinality values.
- `[ ]` Error ownership is recorded on the span that owns the final failed outcome.
- `[ ]` Context propagation is handled across relevant network or async boundaries.
- `[ ]` The changed files were re-read after implementation to verify the final state.
- `[ ]` The final answer includes this checklist, file evidence, and any remaining risks or gaps.

## Final Review Format

In the final answer, include the checklist in this format:
- `[x]` Boundary choice. Evidence: `path/to/file:line` and a short reason.
- `[~]` Propagation. Reason: no network or async boundary exists in this change.
- `[ ]` Semantic conventions. Missing evidence; re-check required.
