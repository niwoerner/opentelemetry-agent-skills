#!/usr/bin/env bash

set -euo pipefail

script_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
repo=${OTEL_SEMCONV_REPO:-}

if [[ -z "$repo" ]]; then
  printf 'OTEL_SEMCONV_REPO must point to a semantic-conventions checkout\n' >&2
  exit 2
fi

query=(env OTEL_SEMCONV_REPO="$repo" OTEL_SEMCONV_TAG=v1.44.0 "$script_dir/query-otel-semantic-conventions.sh")

messaging_spans=$("${query[@]}" messaging spans)
grep -q '^span[.]messaging[.]aws[.]sqs[.]send[.]producer' <<<"$messaging_spans"

messaging_common=$("${query[@]}" messaging common)
grep -q '^messaging[.]kafka[.]attributes[.]common' <<<"$messaging_common"
grep -q '^messaging[.]rabbitmq[.]attributes[.]destination' <<<"$messaging_common"

exact_entry=$("${query[@]}" messaging span.messaging.aws.sqs.send.producer)
grep -q '^source: .*semantic-conventions/blob/v1[.]44[.]0/model/messaging/aws[.]yaml#L' <<<"$exact_entry"

legacy_http=$(env OTEL_SEMCONV_REPO="$repo" OTEL_SEMCONV_TAG=v1.43.0 "$script_dir/query-otel-semantic-conventions.sh" http spans)
grep -q '^span[.]http[.]client' <<<"$legacy_http"

legacy_messaging=$(env OTEL_SEMCONV_REPO="$repo" OTEL_SEMCONV_TAG=v1.43.0 "$script_dir/query-otel-semantic-conventions.sh" messaging spans)
grep -q '^messaging[.]kafka' <<<"$legacy_messaging"
