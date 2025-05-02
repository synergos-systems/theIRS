#!/usr/bin/env bash
set -euo pipefail

SRC_ROOT="./data/990_xsd/output/generated_templates"
DST="./models"

# 1. Wipe & recreate destination
rm -rf "$DST"
mkdir -p "$DST"

# 2. Find every model.go two levels down, adjust package line,
#    and copy to DST as <ParentDir>.go
find "$SRC_ROOT" -mindepth 2 -maxdepth 2 -type f -name 'models.go' | while read -r file; do
  parent="$(basename "$(dirname "$file")")"
  # change `package whatever` → `package models`
  sed '1s/^package .*/package models/' "$file" > "$DST/${parent}.go"
done

echo "✅ Merged into $DST with $(ls -1 "$DST" | wc -l) files"

