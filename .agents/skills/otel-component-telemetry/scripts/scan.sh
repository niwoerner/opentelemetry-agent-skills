#!/usr/bin/env bash
set -euo pipefail

# One headless pi agent per component in config.json, then rebuild index.md.
# The prompt is SKILL.md itself, rendered per component with envsubst.
# Versions that already have a file are skipped unless options.force is true.

SKILL_DIR="$(cd "$(dirname "$0")/.." && pwd)"
OUT_DIR="$(cd "$SKILL_DIR/../../.." && pwd)/skills/otel-components-telemetry"
CONFIG="$SKILL_DIR/config.json"
PI_CMD="${PI_CMD:-pi -p --provider openai-codex --model gpt-5.6-luna --thinking medium}"
WORK_DIR="/tmp/otel-component-telemetry"
mkdir -p "$WORK_DIR"

parallel_agents="$(jq -r '.options.parallel_agents' "$CONFIG")"
last_versions="$(jq -r '.options.last_versions' "$CONFIG")"
force="$(jq -r '.options.force // false' "$CONFIG")"
cleanup="$(jq -r '.options.cleanup // false' "$CONFIG")"

if [ "$cleanup" = true ]; then
  trap 'rm -rf "$WORK_DIR"' EXIT
fi

# Clone each unique repo once (blobless; file contents are fetched lazily on
# checkout), or update tags if it is already there from a previous run.
jq -r '.components[].repo' "$CONFIG" | sort -u | while read -r repo; do
  dir="$WORK_DIR/$(basename "$repo")"
  if [ -d "$dir" ]; then
    git -C "$dir" fetch --tags
  else
    git clone --filter=blob:none "https://github.com/$repo" "$dir"
  fi
done

jq -c '.components[]' "$CONFIG" | {
  while read -r entry; do
    repo="$(jq -r .repo <<<"$entry")"
    path="$(jq -r .path <<<"$entry")"

    # Skip when all of the last N repo tags already have a file. In monorepos
    # the package version differs from the repo tag, so this never matches
    # there and the agent decides (it applies the same rule per version).
    if [ "$force" != true ]; then
      missing=0
      for tag in $(git -C "$WORK_DIR/$(basename "$repo")" tag | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$' | sort -V | tail -n "$last_versions"); do
        [ -f "$OUT_DIR/$(basename "$repo")/$path/$tag.md" ] || missing=1
      done
      if [ "$missing" = 0 ]; then
        echo "skip $(basename "$repo")/$path — all versions scanned" >&2
        continue
      fi
    fi

    # Strip the frontmatter — pi would parse a prompt starting with "---" as a
    # CLI option — then fill in the placeholders.
    prompt="$(awk 'NR==1 && $0=="---" {skip=1; next} skip && $0=="---" {skip=0; next} !skip' "$SKILL_DIR/SKILL.md" |
      Repo="$repo" Path="$path" Version="$last_versions" Force="$force" OutDir="$OUT_DIR" \
      envsubst '${Repo} ${Path} ${Version} ${Force} ${OutDir}')"
    $PI_CMD "$prompt" &

    while [ "$(jobs -rp | wc -l)" -ge "$parallel_agents" ]; do sleep 1; done
  done
  wait
}

{
  echo "# otel-components-telemetry index"
  echo
  echo "| File | Version | Last verified |"
  echo "|---|---|---|"
  find "$OUT_DIR" -name 'v*.md' | sort | while read -r f; do
    fm() { awk -F': *' "/^$1:/{print \$2; exit}" "$f"; }
    echo "| ${f#"$OUT_DIR"/} | $(fm version) | $(fm last_verified) |"
  done
} >"$OUT_DIR/index.md"
