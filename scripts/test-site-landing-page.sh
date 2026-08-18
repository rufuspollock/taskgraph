#!/usr/bin/env bash
set -euo pipefail

file="site/homepage-a.html"
file2="site/index.md"

if [ ! -f "$file" ]; then
  echo "FAIL: $file does not exist"
  exit 1
fi

checks=(
  "TaskGraph"
  "curl -fsSL"
  "id=\"demo\""
  "tg add"
)

for pattern in "${checks[@]}"; do
  if ! grep -q "$pattern" "$file"; then
    echo "FAIL: missing pattern '$pattern' in $file"
    exit 1
  fi
done

echo "PASS: homepage-a landing page contains required sections"

if [ ! -f "$file2" ]; then
  echo "FAIL: $file2 does not exist"
  exit 1
fi

checks2=(
  "TaskGraph"
  "curl -fsSL"
  "id=\"quickstart\""
  "id=\"vision\""
  "<svg"
  "graph-svg"
  "story-path"
  "story-cursor"
  "story-caption"
  "needs breakdown"
  "review notes"
  "prefers-reduced-motion"
  "quickstart-grid"
  "min-width: 0"
)

for pattern in "${checks2[@]}"; do
  if ! grep -q "$pattern" "$file2"; then
    echo "FAIL: missing pattern '$pattern' in $file2"
    exit 1
  fi
done

for disallowed in "Project A" "Project B" "Project C" "svg-edge-focus" "svg-node-focus"; do
  if grep -q "$disallowed" "$file2"; then
    echo "FAIL: found obsolete label '$disallowed' in $file2"
    exit 1
  fi
done

if awk '
  /<svg/ { inside_svg = 1 }
  inside_svg && /^[[:space:]]*$/ { found_blank = 1 }
  /<\/svg>/ { inside_svg = 0 }
  END { exit(found_blank ? 0 : 1) }
' "$file2"; then
  echo "FAIL: blank line inside SVG will terminate the Markdown raw HTML block"
  exit 1
fi

if grep -q 'class="edge active"' "$file2"; then
  echo "FAIL: found obsolete HTML edge hero markup in $file2"
  exit 1
fi

echo "PASS: primary Markdown landing page contains the narrative graph animation"
