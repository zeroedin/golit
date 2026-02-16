#!/bin/bash
# Generate one Hugo page per RHDS element.
# Each page contains ALL demo HTML files for that element inlined as sections.
# Each page gets per-page imports based on what custom elements appear in the demos.

RHDS_ROOT="$HOME/Sites/_rhds/red-hat-design-system"
HUGO_CONTENT="content/elements"
IMPORTABLE="/tmp/rhds-importable.txt"
LIGHTDOM_DIR="static/rhds/css"

# Build list of importable element names
ls -d node_modules/@rhds/elements/elements/rh-*/ | sed 's|.*/elements/||;s|/||' | sort > "$IMPORTABLE"

# Build lightdom CSS lookup file: "element-name filename.css" per line
LIGHTDOM_LOOKUP="/tmp/rhds-lightdom-lookup.txt"
: > "$LIGHTDOM_LOOKUP"
for css in "$LIGHTDOM_DIR"/*.css; do
  bn=$(basename "$css")
  el=$(echo "$bn" | sed 's/-lightdom.*\.css//')
  echo "$el $bn" >> "$LIGHTDOM_LOOKUP"
done

for elem_dir in "$RHDS_ROOT"/elements/rh-*/; do
  elem=$(basename "$elem_dir")
  demo_dir="$elem_dir/demo"

  # Skip elements with no demos
  [ -d "$demo_dir" ] || continue
  demos=("$demo_dir"/*.html)
  [ ${#demos[@]} -eq 0 ] && continue

  # Pretty name from data.yaml or slug
  elem_name=""
  if [ -f "$elem_dir/docs/data.yaml" ]; then
    elem_name=$(grep '^name:' "$elem_dir/docs/data.yaml" | head -1 | sed 's/^name: *//' | tr -d '"')
  fi
  [ -z "$elem_name" ] && elem_name=$(echo "$elem" | sed 's/^rh-//' | sed 's/-/ /g' | awk '{for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) tolower(substr($i,2))}1')

  # Collect all unique rh-* tags used across demos, filter to importable only
  all_tags=""
  for demo_file in "${demos[@]}"; do
    all_tags="$all_tags $(grep -oE '<rh-[a-z-]+' "$demo_file" | sed 's/<//' | sort -u)"
  done
  import_tags=$(echo "$all_tags" | tr ' ' '\n' | sort -u | grep -Fxf "$IMPORTABLE" | grep -v '^$')

  # Use .md so Hugo's goldmark renderer passes raw HTML through without
  # template-parsing (with markup.goldmark.renderer.unsafe = true in hugo.toml).
  # This avoids issues with {{ }} and < > in code block content scripts.
  outfile="$HUGO_CONTENT/${elem}.md"

  {
    # Frontmatter with per-page imports and lightdom CSS
    echo "---"
    echo "title: \"$elem_name\""
    echo "imports:"
    for tag in $import_tags; do
      echo "  - $tag"
    done

    # Collect lightdom CSS files for all imported elements + primary element (deduplicated)
    lightdom_files=""
    for tag in $import_tags $elem; do
      css=$(grep "^${tag} " "$LIGHTDOM_LOOKUP" | awk '{print $2}')
      if [ -n "$css" ]; then
        lightdom_files="$lightdom_files $css"
      fi
    done
    lightdom_unique=$(echo "$lightdom_files" | tr ' ' '\n' | sort -u | grep -v '^$')
    if [ -n "$lightdom_unique" ]; then
      echo "lightdom:"
      for lf in $lightdom_unique; do
        echo "  - $lf"
      done
    fi

    echo "---"
    echo ""
    echo "<p>${#demos[@]} demos for <code>&lt;${elem}&gt;</code></p>"
    echo ""

    # Inline each demo
    for demo_file in "${demos[@]}"; do
      demo_basename=$(basename "$demo_file" .html)
      demo_title=$(echo "$demo_basename" | sed 's/-/ /g')

      # Use Markdown heading (not HTML <div>) so Goldmark processes
      # fenced code blocks inside the demo section.
      echo ""
      echo "### ${demo_title}"
      echo ""

      # Process demo file: strip top-level scripts/styles/links,
      # convert rh-code-block content scripts to Markdown fenced code blocks
      # (Hugo's render-codeblock hook converts these back to <rh-code-block>).
      python3 -c "
import sys, re

html = open('$demo_file').read()

# Extract rh-code-block elements and convert to Markdown fenced code blocks
def code_block_to_md(m):
    attrs = m.group(1)  # e.g. ' actions=\"copy wrap\"'
    inner = m.group(2)
    # Find script content
    script_m = re.search(r'<script\s+type=\"([^\"]+)\">(.*?)</script>', inner, re.DOTALL)
    if not script_m:
        return ''
    mime = script_m.group(1)
    code = script_m.group(2)
    # Map mime type to language
    lang_map = {
        'text/html': 'html', 'text/css': 'css',
        'application/javascript': 'javascript', 'text/javascript': 'javascript',
        'text/yaml': 'yaml', 'text/bash': 'bash',
        'text/python': 'python', 'application/json': 'json',
        'text/plain': 'text', 'text/go': 'go',
        'text/typescript': 'typescript', 'text/ruby': 'ruby',
    }
    lang = lang_map.get(mime, 'text')
    # Also extract any non-script slot content (labels etc)
    slots = re.findall(r'<span\s+slot=[^>]+>[^<]*</span>', inner)
    prefix = '\n'.join(slots) + '\n' if slots else ''
    fence = chr(96)*3  # triple backtick
    return prefix + '\n' + fence + lang + '\n' + code.strip() + '\n' + fence + '\n'

html = re.sub(r'<rh-code-block([^>]*)>(.*?)</rh-code-block>', code_block_to_md, html, flags=re.DOTALL)

# Strip remaining top-level scripts, styles, and links
html = re.sub(r'<script[^>]*>.*?</script>', '', html, flags=re.DOTALL)
html = re.sub(r'<style[^>]*>.*?</style>', '', html, flags=re.DOTALL)
html = re.sub(r'<link[^>]*/?>', '', html)

# Remove blank lines
html = '\n'.join(l for l in html.splitlines() if l.strip())
print(html)
"

      echo ""
    done
  } > "$outfile"

  n_imports=$(echo "$import_tags" | wc -w | tr -d ' ')
  echo "  $elem ($elem_name): ${#demos[@]} demos, $n_imports imports"
done

echo ""
total_pages=$(ls "$HUGO_CONTENT"/rh-*.md 2>/dev/null | wc -l | tr -d ' ')
total_demos=$(grep -c '^### ' "$HUGO_CONTENT"/rh-*.md 2>/dev/null | awk -F: '{s+=$2} END{print s}')
echo "Total: $total_pages pages, $total_demos demos"
