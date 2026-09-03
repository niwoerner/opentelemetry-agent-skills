# `tail_sampling`: configuration

## Typical config

```yaml
processors:
  tail_sampling:
    decision_wait: 10s
    num_traces: 50000
    policies:
      # Keep all error traces
      - name: errors
        type: status_code
        status_code:
          status_codes: [ERROR]
      # Sample 10% of everything else
      - name: baseline
        type: probabilistic
        probabilistic:
          sampling_percentage: 10

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [tail_sampling]
      exporters: [otlphttp]
```

## Configuration reference (top-level)

| Key | Type | Default | Notes |
|-----|------|---------|-------|
| `sampling_strategy` | enum | `trace-complete` | `trace-complete`: evaluate the accumulated trace on the timer path (most flexible, higher memory). `span-ingest`: evaluate each incoming batch at ingest — terminal `sampled`/`dropped` outcomes finalize immediately, non-terminal traces finalize as not-sampled on cleanup (rejects stateful policies). Invalid values fail validation. |
| `decision_wait` | duration | `30s` | Time from first span arrival before the decision is made. Buffers spans for this long. Under `span-ingest` it instead controls pending-cleanup finalization timing. |
| `decision_wait_after_root_received` | duration | `0s` | Decide this long after the root span arrives; `0` disables (only `decision_wait` is used). |
| `num_traces` | int | `50000` | Max traces kept in memory. When full, the oldest are evicted (before decision) unless `block_on_overflow`. |
| `num_shards` | int | `1` | Parallel in-process event loops, with traces assigned by trace-ID hash. Maximum `256`; values >1 cannot be combined with `tail_storage`. Aggregate capacities and per-second policy limits are divided across shards. Added in v0.159.0. |
| `expected_new_traces_per_sec` | int | `0` | Hint for pre-allocating the trace buffer; `0` disables pre-allocation. |
| `sample_on_first_match` | bool | `false` | Stop evaluating and sample as soon as one policy matches. Do not combine with tracestate handling: later policies may report a less strict threshold. |
| `block_on_overflow` | bool | `false` | Block ingest instead of dropping the oldest traces when `num_traces` is reached. |
| `drop_pending_traces_on_shutdown` | bool | `false` | On shutdown, drop pending traces instead of deciding with partial data. |
| `maximum_trace_size_bytes` | int | `0` | Traces larger than this are dropped immediately; `0` disables. |
| `decision_cache.sampled_cache_size` | int | `0` | LRU cache of sampled trace IDs so late spans inherit the decision; `0` disables. |
| `decision_cache.non_sampled_cache_size` | int | `0` | LRU cache of dropped trace IDs; `0` disables. |
| `policies` | list | (required) | Sampling policies. At least one is required. |

`tail_storage` (a component ID) offloads span buffering to a tail-storage extension instead of memory, but is behind the alpha `processor.tailsamplingprocessor.tailstorageextension` feature gate — setting it without the gate fails validation. It cannot be combined with `num_shards > 1`.

The alpha `processor.tailsamplingprocessor.usetracestate` gate (off by default) makes the
`probabilistic`, `rate_limiting`, and `bytes_limiting` policies report effective probability
sampling thresholds and rewrites outgoing `th` values. It falls back to the legacy trace-ID hash
when no probability sampling information is present. Do not combine it with
`sample_on_first_match`. See [Tracestate probability sampling](policies.md#tracestate-probability-sampling).

For the full catalog of policy types and their sub-fields, see [Policy types](policies.md).
