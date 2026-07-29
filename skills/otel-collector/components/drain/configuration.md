# `drain`: configuration

The Drain algorithm tokenizes each log line and walks a **fixed-depth parse tree** (depth set by `tree_depth`). Lines with similar token structure land in the same **cluster**; each cluster carries a template. As more logs arrive, the template is refined — tokens that vary across a cluster's members become `<*>` wildcards while stable tokens are kept verbatim. The algorithm is deterministic: given the same configuration and a representative sample of log forms, the same token structure always yields the same template. Because of this, independent collector instances converge on identical templates once each has seen enough lines.

## Typical config

```yaml
processors:
  drain:
    template_attribute: log.record.template

service:
  pipelines:
    logs:
      receivers: [otlp]
      processors: [memory_limiter, drain]
      exporters: [otlphttp]
```

## Configuration reference

| Key | Type | Default | Notes |
|-----|------|---------|-------|
| `tree_depth` | int | `4` | Max depth of the Drain parse tree (`depth` in the paper). Higher = more specific templates. **Minimum: 3.** |
| `merge_threshold` | float | `0.4` | Minimum fraction of tokens that must match an existing cluster for a line to merge into it rather than form a new cluster (`st` in the paper). Range `[0.0, 1.0]`. |
| `max_node_children` | int | `100` | Maximum children per internal parse-tree node (`maxChild` in the paper). Bounds memory on high-cardinality token positions. |
| `max_clusters` | int | `0` | Maximum clusters tracked; when exceeded, the least-recently-used cluster is evicted. `0` = unlimited. |
| `extra_delimiters` | []string | `[]` | Additional token delimiters beyond whitespace (e.g. `[",", ":"]`). |
| `body_field` | string | `""` | If set and the body is a map, the value of this **single top-level key** is templated instead of the full body. Full OTTL paths are not supported. `""` = use the full body string. |
| `template_attribute` | string | `"log.record.template"` | Attribute key the derived template string is written to. |
| `extract_parameters` | bool | `false` | When `true`, write the body tokens occupying `<*>` positions in the matched template as a string slice. Added in v0.157.0. |
| `params_attribute` | string | `"log.record.template.params"` | Attribute key for the extracted string slice. Consulted only when `extract_parameters: true`. |
| `warmup_min_clusters` | int | `0` | Distinct clusters that must be observed before annotation starts. `0` disables warmup. |
| `storage` | component ID | `""` | ID of a storage extension used to persist the tree across restarts. `""` = disabled. |
| `save_interval` | duration | `0s` | Interval between periodic snapshot saves. `0s` = save on shutdown only. Requires `storage`. |
| `seed_templates` | []string | `[]` | Template strings pre-loaded at startup. Empty/whitespace-only entries are skipped. |
| `seed_logs` | []string | `[]` | Raw example log lines trained on at startup. Empty/whitespace-only entries are skipped. |

## Validation rules

- `tree_depth` must be `>= 3`.
- `merge_threshold` must be in `[0.0, 1.0]`.
- `warmup_min_clusters` must be `>= 0`.
- `save_interval` must be `>= 0`.
- `save_interval > 0` requires `storage` to be set.

## Structured (map) body behavior

If the body is a map, the full serialized form is templated by default, which yields low-value templates. Set `body_field` to the **single top-level** message key (full OTTL paths like `body["event"]["message"]` are not supported), or promote that field to a plain string body upstream (e.g. a `move` operator in the filelog receiver). The full body is used unchanged if the key is absent or the body is not a map.

## Positional parameter extraction

With `extract_parameters: true`, the processor writes the body token at each `<*>` position in template order:

```yaml
processors:
  drain:
    extract_parameters: true
    params_attribute: log.record.template.params
```

For body `user alice logged in from 10.0.0.1` and template `user <*> logged in from <*>`, the parameter attribute is `["alice", "10.0.0.1"]`. Extraction uses the same tokenization as the parse tree (whitespace plus `extra_delimiters`), is positional and unnamed, and writes no slice when the matched template contains no `<*>`. During warmup suppression, neither the template nor parameter attribute is written.
