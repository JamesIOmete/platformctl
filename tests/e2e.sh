#!/bin/bash
set -e

# Setup isolated environment
TEST_DIR=$(mktemp -d)
export HOME="$TEST_DIR"
BIN_PATH="/usr/local/bin/platformctl"

echo "Running E2E Tests for platformctl..."
echo "Temp Home: $TEST_DIR"

# Helper for assertions
assert_contains() {
    if [[ "$1" != *"$2"* ]]; then
        echo "❌ FAIL: Expected output to contain '$2', got:"
        echo "$1"
        exit 1
    else
        echo "✅ PASS: Output contains '$2'"
    fi
}

echo "---------------------------------------------------"
echo "TEST 1: Onboarding (Init)"
# Simulate inputs: Username="testuser", Choice="3" (Full Access)
# We pipe these inputs into the interactive command
output=$(printf "testuser\n3\n" | $BIN_PATH init)
assert_contains "$output" "Configuration saved"
assert_contains "$output" "Setup complete"

# Verify config file exists
if [ -f "$TEST_DIR/.config/platformctl/config.yaml" ]; then
    echo "✅ PASS: Config file created."
else
    echo "❌ FAIL: Config file missing."
    exit 1
fi

echo "---------------------------------------------------"
echo "TEST 2: Secrets Management"
# Set
output=$($BIN_PATH secrets set "api_key" "super_secret_value")
assert_contains "$output" "successfully"

# Get
output=$($BIN_PATH secrets get "api_key")
assert_contains "$output" "super_secret_value"

# List
output=$($BIN_PATH secrets ls)
assert_contains "$output" "api_key"
echo "✅ PASS: Secrets CRUD verified."

echo "---------------------------------------------------"
echo "TEST 3: Fleet (Built-in)"
# List
output=$($BIN_PATH fleet ls)
assert_contains "$output" "robot-001"
assert_contains "$output" "online"

# Status
output=$($BIN_PATH fleet status robot-001)
assert_contains "$output" "Battery"
echo "✅ PASS: Fleet commands verified."

echo "---------------------------------------------------"
echo "TEST 4: Simulation (Plugin)"
# Verify sim plugin discovery
output=$($BIN_PATH help)
assert_contains "$output" "sim"

# Run Simulation
# Use a random scenario name to track it
output=$($BIN_PATH sim run --scenario=e2e_test_scenario)
assert_contains "$output" "Simulation submitted"

# List and verify the scenario exists
output=$($BIN_PATH sim ls)
assert_contains "$output" "e2e_test_scenario"
echo "✅ PASS: Sim plugin verified."

echo "---------------------------------------------------"
echo "🎉 ALL TESTS PASSED"

# Cleanup
rm -rf "$TEST_DIR"
