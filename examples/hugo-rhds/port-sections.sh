#!/bin/bash
# Port RHDS docs sections (patterns, foundations, get-started, theming, etc.) from 11ty to Hugo
RHDS_DOCS="$HOME/Sites/_rhds/red-hat-design-system/docs"
HUGO_CONTENT="content"

# Function to port a markdown file, stripping Nunjucks
port_md() {
  local src="$1"
  local dest="$2"
  local title="$3"

  mkdir -p "$(dirname "$dest")"

  local has_fm=false
  if head -1 "$src" | grep -q '^---'; then
    has_fm=true
  fi

  {
    if $has_fm; then
      # Pass through existing frontmatter, content follows
      cat "$src"
    else
      echo "---"
      echo "title: \"$title\""
      echo "---"
      echo ""
      cat "$src"
    fi
  } | \
  sed 's/{%[^%]*%}//g' | \
  sed 's/{{[^}]*}}//g' | \
  sed 's/<uxdot-pattern[^>]*>/<div class="example">/g' | \
  sed 's/<\/uxdot-pattern>/<\/div>/g' | \
  sed 's/<uxdot-example[^>]*>/<div class="example">/g' | \
  sed 's/<\/uxdot-example>/<\/div>/g' | \
  sed '/<uxdot-feedback>/,/<\/uxdot-feedback>/d' | \
  sed '/<script[^>]*data-helmet[^>]*>/,/<\/script>/d' | \
  sed '/<link[^>]*data-helmet[^>]*>/d' \
  > "$dest"
}

# Function to port an HTML demo file
port_demo() {
  local src="$1"
  local dest="$2"
  local title="$3"

  mkdir -p "$(dirname "$dest")"

  {
    echo "---"
    echo "title: \"$title\""
    echo "---"
    echo ""
    sed '/<script[^>]*type="module"[^>]*>/,/<\/script>/d' "$src" | \
    sed '/<script[^>]*>/,/<\/script>/d'
  } > "$dest"
}

# ============================================================
# PATTERNS
# ============================================================
echo "=== Porting patterns ==="
cat > "$HUGO_CONTENT/patterns/_index.md" << 'EOF'
---
title: "Patterns"
---
Common UI patterns built with RHDS elements.
EOF

for pattern_dir in "$RHDS_DOCS"/patterns/*/; do
  [ -d "$pattern_dir" ] || continue
  pname=$(basename "$pattern_dir")
  target="$HUGO_CONTENT/patterns/$pname"
  mkdir -p "$target"

  # Get title from first md file's frontmatter
  ptitle=$(echo "$pname" | sed 's/-/ /g' | awk '{for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) tolower(substr($i,2))}1')

  cat > "$target/_index.md" << PEOF
---
title: "$ptitle"
isPatternIndex: true
---
PEOF

  # Port markdown files
  for md in "$pattern_dir"/*.md; do
    [ -f "$md" ] || continue
    bn=$(basename "$md")
    page_title=$(echo "${bn%.md}" | sed 's/-/ /g' | awk '{for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) tolower(substr($i,2))}1')
    port_md "$md" "$target/$bn" "$page_title"
  done

  # Port pattern demo HTML files
  if [ -d "$pattern_dir/patterns" ]; then
    for html in "$pattern_dir/patterns/"*.html; do
      [ -f "$html" ] || continue
      bn=$(basename "$html")
      demo_title=$(echo "${bn%.html}" | sed 's/-/ /g' | awk '{for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) tolower(substr($i,2))}1')
      port_demo "$html" "$target/${bn%.html}-demo.html" "$demo_title"
    done
  fi

  count=$(find "$target" -type f | wc -l)
  echo "  $pname: $count files"
done

# ============================================================
# SIMPLE SECTIONS (foundations, get-started, theming, etc.)
# ============================================================
for section in foundations get-started theming tokens about release-notes accessibility support personalization design-code-status; do
  src_dir="$RHDS_DOCS/$section"
  [ -d "$src_dir" ] || continue

  echo "=== Porting $section ==="
  target="$HUGO_CONTENT/$section"
  mkdir -p "$target"

  stitle=$(echo "$section" | sed 's/-/ /g' | awk '{for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) tolower(substr($i,2))}1')

  # Create section index if not present
  if [ ! -f "$target/_index.md" ]; then
    cat > "$target/_index.md" << SEOF
---
title: "$stitle"
---
SEOF
  fi

  # Port all markdown files recursively
  find "$src_dir" -name '*.md' -type f | while read -r md; do
    rel=$(realpath --relative-to="$src_dir" "$md" 2>/dev/null || python3 -c "import os.path; print(os.path.relpath('$md', '$src_dir'))")
    # Create subdirectories if nested
    subdir=$(dirname "$rel")
    if [ "$subdir" != "." ]; then
      mkdir -p "$target/$subdir"
      # Create sub-section _index.md if needed
      if [ ! -f "$target/$subdir/_index.md" ]; then
        sub_title=$(basename "$subdir" | sed 's/-/ /g' | awk '{for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) tolower(substr($i,2))}1')
        cat > "$target/$subdir/_index.md" << SUBEOF
---
title: "$sub_title"
---
SUBEOF
      fi
    fi
    page_title=$(basename "${md%.md}" | sed 's/-/ /g' | awk '{for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) tolower(substr($i,2))}1')
    port_md "$md" "$target/$rel" "$page_title"
  done

  # Port any HTML demo files
  find "$src_dir" -name '*.html' -type f | while read -r html; do
    rel=$(realpath --relative-to="$src_dir" "$html" 2>/dev/null || python3 -c "import os.path; print(os.path.relpath('$html', '$src_dir'))")
    subdir=$(dirname "$rel")
    [ "$subdir" != "." ] && mkdir -p "$target/$subdir"
    demo_title=$(basename "${html%.html}" | sed 's/-/ /g' | awk '{for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) tolower(substr($i,2))}1')
    port_demo "$html" "$target/$rel" "$demo_title"
  done

  count=$(find "$target" -type f | wc -l)
  echo "  $section: $count files"
done

echo ""
echo "=== Summary ==="
echo "Total content files: $(find "$HUGO_CONTENT" -type f | wc -l)"
