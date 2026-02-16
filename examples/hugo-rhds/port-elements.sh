#!/bin/bash
# Port RHDS element docs from 11ty to Hugo
# Source: ~/Sites/_rhds/red-hat-design-system/elements/rh-*/
# Target: content/elements/rh-*/

RHDS_ROOT="$HOME/Sites/_rhds/red-hat-design-system"
HUGO_CONTENT="content/elements"

# Clean and recreate
rm -rf "$HUGO_CONTENT"/rh-*

for elem_dir in "$RHDS_ROOT"/elements/rh-*/; do
  elem=$(basename "$elem_dir")
  docs_dir="$elem_dir/docs"
  demo_dir="$elem_dir/demo"

  # Skip if no docs directory
  [ -d "$docs_dir" ] || continue

  target="$HUGO_CONTENT/$elem"
  mkdir -p "$target"

  # Read element name from data.yaml if it exists
  elem_name=""
  elem_status=""
  if [ -f "$docs_dir/data.yaml" ]; then
    elem_name=$(grep '^name:' "$docs_dir/data.yaml" | head -1 | sed 's/^name: *//' | tr -d '"')
    elem_status=$(grep '^overallStatus:' "$docs_dir/data.yaml" | head -1 | sed 's/^overallStatus: *//' | tr -d '"')
  fi
  [ -z "$elem_name" ] && elem_name=$(echo "$elem" | sed 's/^rh-//' | sed 's/-/ /g' | awk '{for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) tolower(substr($i,2))}1')

  # Create section _index.md
  cat > "$target/_index.md" << INDEXEOF
---
title: "$elem_name"
description: "RHDS $elem_name component ($elem)"
status: "${elem_status:-ready}"
isElementIndex: true
---
INDEXEOF

  # Port each doc markdown file
  for md_file in "$docs_dir"/*.md; do
    [ -f "$md_file" ] || continue
    basename_md=$(basename "$md_file")

    # Map numbered filenames to clean names
    case "$basename_md" in
      00-overview.md) out_name="overview.md" ; page_title="Overview" ;;
      10-style.md)    out_name="style.md"    ; page_title="Style" ;;
      20-guidelines.md) out_name="guidelines.md" ; page_title="Guidelines" ;;
      30-code.md)     out_name="code.md"     ; page_title="Code" ;;
      40-accessibility.md) out_name="accessibility.md" ; page_title="Accessibility" ;;
      90-demos.md)    out_name="demos.md"    ; page_title="Demos" ;;
      *)
        # Strip leading number prefix if present
        out_name=$(echo "$basename_md" | sed 's/^[0-9]*-//')
        page_title=$(echo "$out_name" | sed 's/\.md$//' | sed 's/-/ /g' | awk '{for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) tolower(substr($i,2))}1')
        ;;
    esac

    # Check if file has frontmatter
    has_fm=false
    if head -1 "$md_file" | grep -q '^---'; then
      has_fm=true
    fi

    # Process the file: strip Nunjucks, uxdot-* elements, data-helmet scripts
    {
      if $has_fm; then
        # Extract existing frontmatter and add our fields
        awk '
          BEGIN { in_fm=0; printed=0 }
          /^---$/ && !printed { in_fm++; if(in_fm==1) { print; printed=1; next } else { print "title: \"'"$page_title"'\""; print; next } }
          in_fm==1 { print; next }
          in_fm>=2 { print }
        ' "$md_file"
      else
        echo "---"
        echo "title: \"$page_title\""
        echo "---"
        echo ""
        cat "$md_file"
      fi
    } | \
    # Strip Nunjucks tags
    sed 's/{%[^%]*%}//g' | \
    # Strip Nunjucks variables (but not Hugo ones)
    sed 's/{{[^}]*}}//g' | \
    # Convert uxdot-pattern to div.example
    sed 's/<uxdot-pattern[^>]*src="[^"]*"[^>]*>/<div class="example">/g' | \
    sed 's/<\/uxdot-pattern>/<\/div>/g' | \
    # Convert uxdot-example to div.example
    sed 's/<uxdot-example[^>]*>/<div class="example">/g' | \
    sed 's/<\/uxdot-example>/<\/div>/g' | \
    # Strip uxdot-feedback
    sed '/<uxdot-feedback>/,/<\/uxdot-feedback>/d' | \
    # Strip uxdot-header (handled by layout)
    sed '/<uxdot-header[^>]*>/,/<\/uxdot-header>/d' | \
    # Strip data-helmet script/link tags
    sed '/<script[^>]*data-helmet[^>]*>/,/<\/script>/d' | \
    sed '/<link[^>]*data-helmet[^>]*>/d' \
    > "$target/$out_name"

  done

  # Copy demo HTML files
  if [ -d "$demo_dir" ]; then
    for demo_file in "$demo_dir"/*.html; do
      [ -f "$demo_file" ] || continue
      demo_basename=$(basename "$demo_file")
      demo_name=$(echo "$demo_basename" | sed 's/\.html$//' | sed 's/-/ /g' | awk '{for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) tolower(substr($i,2))}1')

      # Create as Hugo content page with the raw demo HTML
      {
        echo "---"
        echo "title: \"$demo_name Demo\""
        echo "---"
        echo ""
        # Strip script tags from demos (imports handled by baseof)
        sed '/<script[^>]*type="module"[^>]*>/,/<\/script>/d' "$demo_file" | \
        sed '/<script[^>]*>/,/<\/script>/d' | \
        sed '/<style[^>]*>/,/<\/style>/d'
      } > "$target/${demo_basename%.html}-demo.html"

    done
  fi

  echo "  Ported: $elem ($elem_name) - $(ls "$target"/*.md "$target"/*.html 2>/dev/null | wc -l) files"
done

echo "Done. Total element files: $(find "$HUGO_CONTENT"/rh-* -type f 2>/dev/null | wc -l)"
