#!/bin/bash
set -e

echo "=== Envy Storage Layer Test ==="
echo

# Clean up
rm -f ~/.envy.json ~/.envy.json.backup

echo "Test 1: Import with first-run (creates vault)"
echo "Expected: Should prompt for new password, create vault, import 4 keys"
echo "Action: Running './envy --import test.env'"
echo "Please enter password: testpass123"
echo "Confirm: testpass123"
echo "Project name: myproject"
echo "Environment: dev"
echo
# Note: This requires manual input
# ./envy --import test.env

echo "Test 2: Verify vault file exists and check structure"
if [ -f ~/.envy.json ]; then
    echo "✅ Vault file created at ~/.envy.json"
    echo
    echo "Vault contents:"
    cat ~/.envy.json | jq .
    echo
    
    echo "Checking vault structure..."
    VERSION=$(cat ~/.envy.json | jq -r '.version')
    SALT=$(cat ~/.envy.json | jq -r '.salt')
    AUTH_HASH=$(cat ~/.envy.json | jq -r '.auth_hash')
    PROJECTS=$(cat ~/.envy.json | jq '.projects | length')
    
    echo "  Version: $VERSION"
    echo "  Salt (first 20 chars): ${SALT:0:20}..."
    echo "  Auth Hash (first 20 chars): ${AUTH_HASH:0:20}..."
    echo "  Number of projects: $PROJECTS"
    
    if [ "$VERSION" == "1" ] && [ ! -z "$SALT" ] && [ ! -z "$AUTH_HASH" ]; then
        echo "✅ Vault structure is correct"
    else
        echo "❌ Vault structure is incorrect"
        exit 1
    fi
else
    echo "❌ Vault file not found"
    exit 1
fi

echo
echo "Test 3: Verify secrets are encrypted"
FIRST_VALUE=$(cat ~/.envy.json | jq -r '.projects[0].keys[0].current.value')
echo "First secret value (should be encrypted base64):"
echo "  ${FIRST_VALUE:0:60}..."

if [[ "$FIRST_VALUE" =~ ^[A-Za-z0-9+/]+=*$ ]]; then
    echo "✅ Value appears to be base64 encoded (likely encrypted)"
else
    echo "❌ Value does not appear to be encrypted"
    exit 1
fi

echo
echo "Test 4: Verify plaintext metadata"
PROJECT_NAME=$(cat ~/.envy.json | jq -r '.projects[0].name')
KEY_NAME=$(cat ~/.envy.json | jq -r '.projects[0].keys[0].key')

echo "Project name: $PROJECT_NAME"
echo "First key name: $KEY_NAME"

if [ "$PROJECT_NAME" != "null" ] && [ "$KEY_NAME" != "null" ]; then
    echo "✅ Metadata is stored in plaintext (correct)"
else
    echo "❌ Metadata is missing"
    exit 1
fi

echo
echo "=== Manual Tests Required ==="
echo
echo "Test 5: Login with correct password"
echo "  Run: ./envy"
echo "  Enter: testpass123"
echo "  Expected: TUI should launch with the imported project"
echo
echo "Test 6: Login with wrong password"
echo "  Run: ./envy"
echo "  Enter: wrongpassword"
echo "  Expected: Error message 'authentication failed: incorrect password'"
echo
echo "Test 7: Export and verify decryption"
echo "  Run: ./envy --export myproject"
echo "  Enter: testpass123"
echo "  Expected: Creates .env with plaintext secrets"
echo "  Verify: cat .env should show API_KEY=sk_test_12345abcdefg"
echo
echo "Test 8: TUI functionality"
echo "  Run: ./envy"
echo "  Test: Navigate, copy to clipboard (30s auto-clear)"
echo

echo "=== Automated Tests Complete ==="
