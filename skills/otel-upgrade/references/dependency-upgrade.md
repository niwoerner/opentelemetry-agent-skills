# Dependency Upgrade

Use this workflow for OpenTelemetry APIs, SDKs, instrumentation libraries, exporters, plugins, and related packages in any ecosystem.

## Assess the dependency change

1. Inspect relevant manifests, lockfiles, imports, toolchain files, CI, configuration, generated code, deployment files, and repository status.
2. Identify the direct dependencies involved, relevant transitive constraints, actual usage, wrappers, plugins, independently versioned packages, and packages that must move together. Do not assume related packages share a version.
3. Query the authoritative package registry and primary upstream release sources for each relevant package. When `otel-sdk-versions` is available, use it as an initial lookup index, then confirm package-specific versions at their release sources.
4. Resolve the candidate dependency family in an isolated environment. Inspect material changes to the resolved graph and generated lockfiles rather than treating successful resolution as sufficient proof.
5. Run the available project-supported checks relevant to the change, such as dependency resolution, lockfile generation, compilation or type checking, tests, configuration validation, and integration checks.
6. Review release notes and migration guidance for changes connected to APIs, features, configuration, protocols, deployment modes, or externally consumed behavior used by the repository. Inspect the corresponding upstream source diff when risk or incomplete documentation warrants it.
7. Check available downstream contracts, including APIs, schemas, persisted data, events, telemetry, dashboards, and alerts, and identify any required rollout ordering.
8. When SDK instrumentation is involved and `otel-telemetry-emissions` is available, compare exact current and candidate telemetry. Report missing component or version coverage as unverified; do not infer emissions from semantic conventions.

## Interpret the evidence

- Dependency resolution shows that the resolver found a graph; it does not prove source or behavioral compatibility.
- Successful compilation covers only the configurations and code paths compiled.
- Passing tests increase confidence only for the behavior and environments they exercise.
- An untested version is a candidate, not the highest validated compatible version.
