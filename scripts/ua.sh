#!/usr/bin/env bash
set -euo pipefail

fail() { echo "FAIL: $1"; exit 1; }
pass() { echo "PASS: $1"; }
reset_config() {
  rm -f "$HOME/.config/platformctl/config.yaml" "$HOME/.config/platformctl/config.yml" "$HOME/.config/platformctl/config.json"
}

# 1) help
out=$(./bin/platformctl help)
grep -q "Built-in commands" <<<"$out" || fail "help missing built-ins"
grep -q "platformctl hello" <<<"$out" || fail "help missing plugin note"
pass "help"

# 2) doctor (expect WARN for config when missing)
reset_config
out=$(./bin/platformctl doctor)
grep -q "\[OK\] Go:" <<<"$out" || fail "doctor missing Go OK"
grep -q "\[OK\] Git:" <<<"$out" || fail "doctor missing Git OK"
grep -q "\[WARN\] Config" <<<"$out" || fail "doctor missing config WARN"
pass "doctor (no config)"

# 3) auth unauthenticated
reset_config
out=$(./bin/platformctl auth status)
grep -q "Authenticated: false" <<<"$out" || fail "auth unauthenticated"
pass "auth unauthenticated"

# 4) auth with scopes
mkdir -p ~/.config/platformctl
cat > ~/.config/platformctl/config.yaml <<'EOF'
principal: demo-user
scopes:
  - fleet:read
  - infra:write
EOF
out=$(./bin/platformctl auth status)
grep -q "Authenticated: true" <<<"$out" || fail "auth authenticated"
grep -q "demo-user" <<<"$out" || fail "auth principal"
grep -q "fleet:read" <<<"$out" || fail "auth scopes"
pass "auth authenticated"

# 5) fleet requires scope
cat > ~/.config/platformctl/config.yaml <<'EOF'
principal: demo-user
scopes:
  - infra:write
EOF
if ./bin/platformctl fleet ls; then fail "fleet should have failed (missing scope)"; else pass "fleet denied without scope"; fi

# 6) fleet success
cat > ~/.config/platformctl/config.yaml <<'EOF'
principal: demo-user
scopes:
  - fleet:read
EOF
out=$(./bin/platformctl fleet ls)
grep -q "robot-001" <<<"$out" || fail "fleet listing missing device"
pass "fleet success"

# 7) env bootstrap scope check
cat > ~/.config/platformctl/config.yaml <<'EOF'
principal: demo-user
scopes:
  - fleet:read
EOF
if ./bin/platformctl env bootstrap dev; then fail "env should have failed (missing infra:write)"; else pass "env denied without scope"; fi

# 8) env bootstrap success
cat > ~/.config/platformctl/config.yaml <<'EOF'
principal: demo-user
scopes:
  - fleet:read
  - infra:write
EOF
out=$(./bin/platformctl env bootstrap dev)
grep -q "Bootstrap plan for dev" <<<"$out" || fail "env plan header"
grep -q "dry-run" <<<"$out" || fail "env plan dry-run note"
pass "env bootstrap success"

# 9) plugin fallback demo
echo -e '#!/usr/bin/env bash\necho plugin works' > /tmp/platformctl-hello
chmod +x /tmp/platformctl-hello
PATH="/tmp:$PATH" out=$(./bin/platformctl hello)
grep -q "plugin works" <<<"$out" || fail "plugin execution"
pass "plugin fallback"

echo "ALL PASS"