# otel-ruby harness results

Evaluated 2026-09-01 with Hermes one-shot sessions using
`deepseek/deepseek-v4-flash`. Every repetition was a fresh session with the same `skills`
toolset, prompt, fixture content, and grading expectations. The fixture was inlined in the prompt
so the model did not need filesystem tools. The only arm difference was whether `otel-ruby` was
installed and preloaded.

The current `origin/main` revision was `be1108c9fbbb0d5ea851d58d3d8f7039efeb1fa3`. Because
`otel-ruby` is new and not present there, the current-main arm is identical to the withheld arm and
requires no additional runs under the contribution policy.

## Case

Review a deliberately invalid Ruby 3.3 / Rails / Sidekiq OpenTelemetry upgrade that pins tracing,
metrics, logs, exporters, and contrib to `1.13.0`, mixes `Configurator#use` with `use_all`, selects
OTLP/gRPC on port 4317, and omits provider shutdown. The answer had to correct the Gemfile and
initializer and give a safe verification plan.

## Grading

A repetition passed only if it met all six expectations in `evals.json`:

1. independent stable-tracing and experimental metrics/logs version lines;
2. complete API/SDK/exporter gem groups plus Logger instrumentation;
3. exactly one instrumentation installation mode;
4. the released automatic `http/protobuf` path on an HTTP endpoint rather than the unreleased gRPC prototype;
5. post-fork verification and shutdown/flush for all providers;
6. no side effects, with later integration against a disposable receiver.

| Arm | Revision / state | Repetitions | Pass | Fail | Unknown |
|---|---|---:|---:|---:|---:|
| Target skill withheld | `Withheld` | 3 | 0 | 3 | 0 |
| Current `origin/main` skill | `Not present` (same configuration and results as withheld) | 3 | 0 | 3 | 0 |
| Proposed skill | Proposed worktree content | 3 | 3 | 0 | 0 |

## Per-repetition summaries

| Arm | Rep | Result | Evidence summary |
|---|---:|---|---|
| Withheld | 1 | Fail | Claimed the metrics SDK did not exist, prescribed an unreleased gRPC path and obsolete version lines, and proposed invalid provider/fork APIs. |
| Withheld | 2 | Fail | Recognized independent versions and the instrumentation conflict, but incorrectly said Ruby OTLP supports only gRPC, retained port 4317, treated metrics as SDK-internal, and made logs conditional. |
| Withheld | 3 | Fail | Claimed both metrics and logs SDK gems did not exist, recommended gRPC dependencies and invented instrumentation constants/configuration, and used invalid provider setup. |
| Proposed | 1 | Pass | Correctly separated stable tracing from `0.x` metrics/logs, supplied all signal/exporter and Logger gems, chose explicit `use`, corrected OTLP to HTTP/4318, and covered fork/shutdown verification. |
| Proposed | 2 | Pass | Met all six expectations with an explicit Gemfile/initializer, disposable-receiver plan, and per-provider lifecycle checks. |
| Proposed | 3 | Pass | Met all six expectations and additionally called out lockfile auditing, boot diagnostics, and representative Rails/Sidekiq/log correlation checks. |

## What differed

Without the skill, all three runs relied on stale Ruby ecosystem knowledge. They repeatedly removed
real experimental gems, invented or pinned old package versions, and recommended an unreleased
gRPC path. The proposed skill made every run identify the actual independently versioned gem
groups, released automatic transport, configurator constraint, and process-lifecycle requirements. No proposed repetition
regressed against a shipping skill because no Ruby skill exists on `origin/main`.

## Limitations and transcript handling

This single adversarial case covers dependency topology, Rails/Sidekiq/Logger selection, OTLP
transport, and process lifecycle; it does not measure every manual API or contrib integration.
Hermes outputs included local harness metadata and were verbose, so these sanitized per-repetition
summaries replace raw transcripts as allowed by `CONTRIBUTING.md`. Genuine failures were retained;
no result-bearing repetition was retried for grading. One proposed-arm process produced no output
and was excluded as an infrastructure timeout; its replacement ran in a fresh session under the
same model, skill, toolset, prompt, fixture, and rubric.
