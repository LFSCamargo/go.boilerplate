#!/usr/bin/env bash
# Download Go Boilerplate email assets from Figma MCP export URLs.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
OUT="$ROOT/emails/static/go_boilerplate"
URLS_FILE="${1:-}"

if [[ -n "$URLS_FILE" && -f "$URLS_FILE" ]]; then
  mkdir -p "$OUT"
  while IFS='=' read -r name url; do
    [[ -z "$name" || -z "$url" ]] && continue
    echo "→ $name"
    curl -sL "$url" -o "$OUT/$name"
  done < "$URLS_FILE"
  echo "Saved assets to $OUT"
  exit 0
fi

echo "Assets directory: $OUT"
echo "Committed assets (verify present):"
for f in hero_image.png logo_mark.svg social_x.svg social_linkedin.svg social_youtube.svg social_github.svg; do
  if [[ -f "$OUT/$f" ]]; then
    echo "  ✓ $f"
  else
    echo "  ✗ missing $f" >&2
    exit 1
  fi
done
