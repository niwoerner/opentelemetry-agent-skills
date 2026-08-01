# Synthetic SPA setup excerpt

This fixture is documentation-only. Do not install packages or contact endpoints.

## Pinned dependencies

The Browser SDK entry is a local build from the upstream GitHub `browser-sdk-v0.2.0` release tag,
not an npm-published 0.2.0 package. The other entries are npm packages.

```json
{
  "@opentelemetry/browser-sdk": "0.2.0",
  "@opentelemetry/browser-instrumentation": "0.7.0",
  "@opentelemetry/instrumentation-fetch": "0.221.0",
  "@opentelemetry/sdk-trace-web": "2.0.1"
}
```

The application origin is `https://shop.example.test`. Its only application API origin is
`https://api.example.test`. Browser telemetry should go to the local development Collector at
`http://127.0.0.1:4318` during verification.

## Current excerpt

```typescript
startBrowserSdk({
  serviceName: 'shop-web',
  exportUrl: 'https://telemetry-backend.example.test',
  exportHeaders: { authorization: 'browser-must-not-hold-backend-credentials' },
  traces: { sampler: new TraceIdRatioBasedSampler(0.1) },
});

new FetchInstrumentation({
  propagateTraceHeaderCorsUrls: [/.*/],
});
```

An internal note written for Browser SDK 0.1.0 says: “The `sampler` field is accepted by the type
but ignored by `startTracesSdk`, so sampling can only happen downstream.” The pinned 0.2.0 release
notes state that trace sampler configuration is now passed to the provider.

The API's development CORS policy currently allows `content-type` only. Production must not be
changed or contacted as part of this review.
