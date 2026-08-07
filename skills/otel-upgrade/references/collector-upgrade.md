# Collector Upgrade

Use this workflow for Collector distributions, components, images, deployment wrappers, configuration, custom builds, and runtime behavior.

## Establish the upgrade surface

Record each applicable layer separately:

| Layer | Evidence to identify |
|---|---|
| Deployment wrapper | Chart, Operator, or other packaging version and manifests |
| Runtime artifact | Distribution, binary version, image tag, and immutable digest when available |
| Build | OCB version, builder manifest, toolchain, and build inputs |
| Components | Compiled receivers, processors, exporters, connectors, extensions, and configuration providers |
| Configuration | Fully resolved configuration and feature gates used by the deployment |
| Embedded modules | Exact core, contrib, stable, custom, and third-party module versions |

Do not use one layer's version as proof of another.

## Use specialized skills when available

- Use `otel-collector` for component configuration, defaults, validation, signal support, stability, and version-specific gotchas.
- Use `otel-collector-builder` for OCB manifests, module alignment, generated builds, and custom distributions.
- Use `otel-telemetry-emissions` for version-pinned emitted telemetry.
- Use `otel-telemetrygen` to generate representative runtime traffic.

If a relevant skill is unavailable, use the exact-version upstream source, release artifact, and local repository evidence instead. State what could not be verified.

## Compare the distributions

1. Identify the exact current and candidate artifacts from repository and runtime evidence.
2. Compare their authoritative distribution or builder manifests and compiled component inventories. Check for added, removed, renamed, deprecated, or stability-changed components and providers, including per-signal support.
3. Confirm that every configured component and configuration provider exists in the candidate artifact. Do not infer a stock distribution's contents from another distribution at the same release.
4. Review release notes, migration guidance, and relevant source changes across the full version interval, but report only changes connected to the deployed distribution, components, configuration, or behavior.

## Check configuration compatibility

1. Resolve the same configuration inputs, includes, environment substitutions, and provider sources for both artifacts without exposing secrets.
2. Validate the unchanged configuration with the exact current and candidate binaries before proposing a migration. Compare fatal errors and warnings, removed or renamed keys, changed defaults and validation rules, component identifiers, pipeline wiring, provider availability, feature gates, and signal support.
3. Validate any proposed migrated configuration separately. Preserve the distinction between an unchanged configuration being compatible and a modified configuration making the upgrade viable.
4. Treat configuration acceptance as evidence that the candidate can load the checked configuration, not that it behaves equivalently at runtime.

## Check the build

For a published distribution, verify the selected artifact, platforms, image metadata, entrypoint, permissions, certificates, and deployment packaging relevant to the repository. Do not rebuild upstream artifacts unless the task requires it.

For a custom distribution:

1. Verify version alignment from the exact target release rather than assuming all modules share one version.
2. Inspect changes to the builder manifest, generated component registration, resolved module graph, and generated dependency files.
3. Build every relevant production target, including applicable toolchains, operating systems, architectures, CGO settings, build tags, and container stages.
4. Run the candidate binary's component inventory and configuration validation commands. Compilation proves only that the checked build target compiled.

## Compare runtime behavior

Run the current and candidate artifacts with equivalent resolved configuration, environment inputs, and representative telemetry. Define acceptable differences before judging the result, and normalize expected nondeterminism such as timestamps, identifiers, ordering, and batch boundaries.

Compare what is material to the deployment:

- exported telemetry content, attributes, cardinality, and protocol behavior
- dropped, duplicated, reordered, batched, or transformed data
- Collector self-telemetry, instrumentation scopes, logs, warnings, and errors
- startup readiness, health, graceful shutdown, and failure handling
- retry, queue, batching, timeout, and backpressure behavior
- persistent queues, checkpoints, write-ahead logs, and storage extensions
- resource use relative to configured limits and operational thresholds

## Check rollout and rollback

Check mixed-version operation, protocol compatibility, deployment ordering, persisted-state compatibility, and whether the previous version can read state written by the candidate. Treat production-only conditions and unavailable external consumers as unverified.
