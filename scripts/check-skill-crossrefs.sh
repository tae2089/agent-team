#!/usr/bin/env bash
set -euo pipefail

status=0

while IFS= read -r skill; do
  if [[ ! -f "skills/${skill}/SKILL.md" ]]; then
    echo "missing skill reference: ${skill}"
    status=1
  fi
done < <(
  rg -o '`([a-z0-9]+-)*[a-z0-9]+`' skills README.md \
    | sed -E 's/.*`([^`]+)`.*/\1/' \
    | rg '^(agent-team|persona-agent-team|recipe-agent-team)-' \
    | sort -u
)

while IFS=: read -r file line href; do
  [[ -z "${href}" ]] && continue
  [[ "${href}" =~ ^https?:// ]] && continue
  [[ "${href}" =~ ^# ]] && continue
  path="${href%%#*}"
  [[ -z "${path}" ]] && continue
  target="$(dirname "${file}")/${path}"
  if [[ ! -e "${target}" ]]; then
    echo "missing markdown link: ${file}:${line}: ${href}"
    status=1
  fi
done < <(
  rg -n -o '\[[^]]+\]\(([^)]+)\)' skills README.md \
    | sed -E 's/^([^:]+):([0-9]+):.*\(([^)]+)\).*$/\1:\2:\3/'
)

exit "${status}"
