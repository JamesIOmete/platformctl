#!/usr/bin/env bash
set -euo pipefail

# ---------- Resolve platformctl binary ----------
if [[ -n "${BIN_PATH:-}" ]]; then
  BIN="$BIN_PATH"
elif command -v platformctl >/dev/null 2>&1; then
  BIN="$(command -v platformctl)"
else
  # try common Go install paths
  GOBIN="$(go env GOBIN 2>/dev/null || true)"
  GOPATH="$(go env GOPATH 2>/dev/null || true)"
  for p in "$GOBIN/platformctl" "$GOPATH/bin/platformctl" "$HOME/go/bin/platformctl"; do
    [[ -x "$p" ]] && BIN="$p" && break
  done
  : "${BIN:=}" || true
fi

if [[ -z "${BIN:-}" || ! -x "$BIN" ]]; then
  echo "❌ platformctl not found. Set BIN_PATH or add it to PATH."
  echo "   Example quick setup:"
  echo "     export GOBIN=\$(mktemp -d) && go install github.com/JamesIOmete/platformctl/cmd/platformctl@latest"
  echo "     export PATH=\"\$GOBIN:\$PATH\""
  exit 1
fi

# ---------- Isolated HOME ----------
TEST_DIR="$(mktemp -d -t platformctl-e2e-XXXXXX)"
export HOME="$TEST_DIR"
export XDG_CONFIG_HOME="$HOME/.config"
trap 'rm -rf "$TEST_DIR"' EXIT

echo "Running E2E Tests for platformctl..."
echo "Binary: $BIN"
echo "Temp Home: $TEST_DIR"
echo "---------------------------------------------------"

# ---------- Helpers ----------
assert_contains() {
  local hay="$1" needle="$2"
  if [[ "$hay" != *"$needle"* ]]; then
    echo "❌ FAIL: expected to find '$needle' in:"
    printf "%s\n" "$hay"
    exit 1
  else
    echo "✅ PASS: found '$needle'"
  fi
}

# ---------- TEST 1: Onboarding ----------
echo "TEST 1: Onboarding (Init)"
# adjust inputs if your init prompts differ
output="$(printf 'testuser\n3\n' | "$BIN" init || true)"
assert_contains "$output" "Configuration saved"
assert_contains "$output" "Setup complete"

cfg="$XDG_CONFIG_HOME/platformctl/config.yaml"
[[ -f "$cfg" ]] && echo "✅ PASS: Config file created: $cfg" || { echo "❌ FAIL: Config file missing"; exit 1; }

echo "---------------------------------------------------"
# ---------- TEST 2: Secrets ----------
echo "TEST 2: Secrets Management"
output="$("$BIN" secrets set api_key super_secret_value)"
assert_contains "$output" "successfully"

output="$("$BIN" secrets get api_key)"
assert_contains "$output" "super_secret_value"

output="$("$BIN" secrets ls)"
assert_contains "$output" "api_key"
echo "✅ PASS: Secrets CRUD verified."

echo "---------------------------------------------------"
# ---------- TEST 3: Fleet ----------
echo "TEST 3: Fleet (Built-in)"
output="$("$BIN" fleet ls)"
assert_contains "$output" "robot-001"
assert_contains "$output" "online"

output="$("$BIN" fleet status robot-001)"
assert_contains "$output" "Battery"
echo "✅ PASS: Fleet commands verified."

echo "---------------------------------------------------"
# ---------- TEST 4: Simulation plugin ----------
echo "TEST 4: Simulation (Plugin)"
if command -v platformctl-sim >/dev/null 2>&1; then
  output="$("$BIN" help)"
  assert_contains "$output" "sim"

  SCEN="e2e_test_scenario_$RANDOM"
  output="$("$BIN" sim run --scenario="$SCEN")"
  assert_contains "$output" "Simulation submitted"

  output="$("$BIN" sim ls)"
  assert_contains "$output" "$SCEN"
  echo "✅ PASS: Sim plugin verified."
else
  echo "ℹ️ Skipping sim tests: platformctl-sim not found on PATH."
  echo "   Install with:"
  echo "     go install github.com/JamesIOmete/platformctl/cmd/platformctl-sim@latest"
fi

echo "---------------------------------------------------"
echo "🎉 ALL TESTS PASSED"
