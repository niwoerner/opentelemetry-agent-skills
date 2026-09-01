# Ruby upgrade and breaking-change review

## Inventory the locked graph

```sh
bundle list | rg 'opentelemetry'
bundle outdated --strict | rg 'opentelemetry'
bundle info opentelemetry-sdk
```

Review the `Gemfile.lock`, not only direct dependencies. Core, semantic conventions, exporters,
experimental signals, and each contrib instrumentation gem publish independently.

## Classify each changed gem

| Group | Typical gems | Review focus |
|---|---|---|
| Stable tracing | `opentelemetry-api`, `opentelemetry-sdk` | API/SDK changelogs, Ruby minimum, sampler/processor behavior |
| Trace export | `opentelemetry-exporter-otlp` | protocol, endpoint construction, TLS, headers, retries |
| Experimental metrics | `opentelemetry-metrics-*`, `opentelemetry-exporter-otlp-metrics` | every minor release may break API, aggregation, cardinality, exemplars, temporality |
| Experimental logs | `opentelemetry-logs-*`, `opentelemetry-exporter-otlp-logs` | every minor release may break API, limits, processors, event fields |
| Contrib | `opentelemetry-instrumentation-*`, propagators, resource detectors | target-gem compatibility, Ruby minimum, span names/attributes, options |
| Semantic conventions | `opentelemetry-semantic_conventions` | generated constants, stability namespace, emitted-schema migration |

Do not align versions by number across groups. A valid bundle can contain `1.x` tracing gems and
different `0.x` metrics, logs, exporter, and instrumentation versions.

## Ruby and target-library compatibility

Check both core and contrib requirements. Current contrib development requires Ruby 3.3 or newer,
but an older locked instrumentation release may have a different range. Each instrumentation
gemspec also constrains the Rails/framework/client versions it supports. Upgrade Ruby, Rails, or a
client library and its OTel instrumentation as one compatibility decision.

## Behavioral review

Read every changed gem's changelog from the old locked version through the proposed one. Search for:

- `BREAKING CHANGE`, minimum Ruby changes, and removed/deprecated APIs;
- new default exporters or protocols;
- sampler, batch processor, fork, shutdown, and timeout changes;
- metric aggregation, temporality, exemplar, or cardinality changes;
- log attribute-limit and event-field changes;
- semantic-convention migrations and span-name/attribute changes;
- expanded or narrowed target-library version constraints.

For contrib, inspect the released gem's README, gemspec, and tests when the changelog does not fully
describe emitted telemetry.

## Verification path

1. Run `bundle update <explicit OTel gems>` rather than a broad unrelated update.
2. Inspect `git diff -- Gemfile Gemfile.lock` and explain every changed OTel/transitive gem.
3. Run the repository's tests and boot checks on each supported Ruby/framework combination.
4. Exercise representative requests, database calls, and jobs against a disposable OTLP receiver.
5. Compare span names, attributes, parentage, metrics temporality/cardinality, logs, and shutdown.
6. Audit dashboards, alerts, sampling rules, and Collector processors for renamed telemetry.

Sources: [core changelogs](https://github.com/open-telemetry/opentelemetry-ruby),
[contrib changelogs](https://github.com/open-telemetry/opentelemetry-ruby-contrib), and
[OpenTelemetry Ruby releases](https://github.com/open-telemetry/opentelemetry-ruby/releases).
