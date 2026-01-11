#!/usr/bin/env bash
set -euo pipefail

# ======== CONFIG ========
TYPE_SPEED=${TYPE_SPEED:-40}   # chars/sec for "typing"
BLOCK_PAUSE=${BLOCK_PAUSE:-2}  # seconds between command blocks
SECTION_PAUSE=${SECTION_PAUSE:-2}
PROMPT="${PROMPT:-$ }"
# If you want to hide your real HOME in outputs, uncomment:
# export DEMO_HOME=$(mktemp -d -t platformctl-XXXXXX)
# export HOME="$DEMO_HOME"; export XDG_CONFIG_HOME="$HOME/.config"
# export GOBIN="$HOME/bin"; export PATH="$GOBIN:$PATH"; export PS1="$ "

# ======== COLORS (respect NO_COLOR) ========
if [[ -t 1 && -z "${NO_COLOR:-}" ]]; then
  BOLD='\033[1m'; DIM='\033[2m'; RESET='\033[0m'
  CYAN='\033[36m'; GREEN='\033[32m'; GREY='\033[90m'
else
  BOLD=''; DIM=''; RESET=''; CYAN=''; GREEN=''; GREY=''
fi

# ======== HELPERS ========
hr(){ printf "${GREY}%s${RESET}\n" "────────────────────────────────────────────────────────"; }
newline(){ printf "\n"; }

say(){  # section header
  newline
  hr
  printf "${BOLD}${CYAN}%s${RESET}\n" "$1"
  hr
  sleep "$SECTION_PAUSE"
}

slowprint(){  # prints one line at TYPE_SPEED cps
  local s="$1" delay
  delay=$(awk "BEGIN{print 1/$TYPE_SPEED}")
  for ((i=0; i<${#s}; i++)); do
    printf "%s" "${s:$i:1}"
    sleep "$delay"
  done
  printf "\n"
}

type_and_run(){ # shows a green typed prompt, then executes
  local cmd="$*"
  printf "${GREEN}%s${RESET}" "$PROMPT"
  slowprint "$cmd"
  eval "$cmd"
  newline
  sleep "$BLOCK_PAUSE"
}

# ======== DEMO ========
clear

say "0) Install (no sudo)"
type_and_run 'go install github.com/JamesIOmete/platformctl/cmd/platformctl@latest'
type_and_run 'go install github.com/JamesIOmete/platformctl/cmd/platformctl-sim@latest'
sleep 3

say "1) Onboard + sanity"
type_and_run 'platformctl init'
type_and_run 'platformctl doctor'

say "2) Fleet basics"
type_and_run 'platformctl fleet ls'
type_and_run 'platformctl fleet status robot-001'

# follow logs for ~2s, then stop cleanly
printf "${GREEN}%s${RESET}" "$PROMPT"; slowprint 'platformctl fleet logs robot-001 --follow'
platformctl fleet logs robot-001 --follow & LOGS_PID=$!
sleep 2
kill "$LOGS_PID" 2>/dev/null || true
newline
sleep "$BLOCK_PAUSE"

say "3) Sim plugin (decoupled team)"
type_and_run 'platformctl sim run --scenario=warehouse_v2'
type_and_run 'platformctl sim ls'

say "4) Secrets story"
type_and_run 'platformctl secrets set api_key "s3cr3t"'
type_and_run 'platformctl secrets get api_key'

# Tidy the sandbox directory if you enabled it
[[ -n "${DEMO_HOME:-}" ]] && rm -rf -- "$DEMO_HOME"

say "5) End-to-end test (isolated HOME)"
# ensure tests can find the binaries you just installed in this session:
export BIN_PATH="$(command -v platformctl || true)"
type_and_run 'bash tests/e2e.sh'

# --- finish cleanly ---
say "✅ Demo complete"
printf "${DIM}(closing in 10s)${RESET}\n"
sleep 10
exit 0
