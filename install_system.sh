#!/bin/bash
set -e

echo "Requested system installation for platformctl..."
echo "This script requires sudo privileges to write to /usr/local/bin."

echo ""
echo "Step 1: Removing stale 'platformctl-fleet' from /usr/local/bin..."
sudo rm -f /usr/local/bin/platformctl-fleet

echo "Step 2: Installing 'platformctl' to /usr/local/bin..."
sudo cp ./bin/platformctl /usr/local/bin/

echo "Step 3: Installing 'platformctl-sim' to /usr/local/bin..."
sudo cp ./bin/platformctl-sim /usr/local/bin/

echo ""
echo "Done! Verifying 'platformctl help' (Look for 'sim' in plugins)..."
/usr/local/bin/platformctl help
