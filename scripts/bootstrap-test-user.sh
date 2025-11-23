#!/bin/bash
# Bootstrap script to create a test user and API key for testing

set -e

API_URL="${PEEKAPING_ENDPOINT:-http://localhost:8034}"
TEST_EMAIL="${TEST_USER_EMAIL:-test@example.com}"
TEST_PASSWORD="${TEST_USER_PASSWORD:-TestPassword123!}"

echo "==> Waiting for API to be ready..."
timeout=60
while [ $timeout -gt 0 ]; do
    if curl -s -f "${API_URL}/api/v1/health" > /dev/null 2>&1; then
        echo "API is ready!"
        break
    fi
    echo "Waiting for API... (${timeout} seconds remaining)"
    sleep 2
    timeout=$((timeout - 2))
done

if [ $timeout -le 0 ]; then
    echo "Error: API did not become ready in time"
    exit 1
fi

echo "==> Attempting to register test user (will fail if user already exists)..."
CREATE_USER_RESPONSE=$(curl -s -X POST "${API_URL}/api/v1/auth/register" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"${TEST_EMAIL}\",\"password\":\"${TEST_PASSWORD}\"}" 2>/dev/null || echo "")

if echo "$CREATE_USER_RESPONSE" | grep -q "admin already exists\|already exists"; then
    echo "User already exists, skipping registration"
elif echo "$CREATE_USER_RESPONSE" | grep -q "success\|accessToken"; then
    echo "User created successfully"
else
    echo "Warning: User registration response: $CREATE_USER_RESPONSE"
fi

echo "==> Logging in to create API key..."
LOGIN_RESPONSE=$(curl -s -X POST "${API_URL}/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"${TEST_EMAIL}\",\"password\":\"${TEST_PASSWORD}\"}" 2>/dev/null || echo "")

if [ -z "$LOGIN_RESPONSE" ]; then
    echo "Error: Failed to login"
    exit 1
fi

# Extract access token from API response (handles both direct and wrapped responses)
# Try to extract from data.accessToken first (wrapped response), then accessToken (direct)
ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"data":{[^}]*"accessToken":"[^"]*' | grep -o '"accessToken":"[^"]*' | cut -d'"' -f4 || echo "")
if [ -z "$ACCESS_TOKEN" ]; then
    ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"accessToken":"[^"]*' | cut -d'"' -f4 || echo "")
fi

if [ -z "$ACCESS_TOKEN" ]; then
    echo "Error: Failed to extract access token from login response"
    echo "Response: $LOGIN_RESPONSE"
    exit 1
fi

echo "==> Creating API key..."
API_KEY_RESPONSE=$(curl -s -X POST "${API_URL}/api/v1/api-keys" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}" \
    -d "{\"name\":\"test-key\"}" 2>/dev/null || echo "")

if [ -z "$API_KEY_RESPONSE" ]; then
    echo "Error: Failed to create API key"
    exit 1
fi

# Extract API key token from API response (handles both direct and wrapped responses)
API_KEY=$(echo "$API_KEY_RESPONSE" | grep -o '"data":{[^}]*"token":"[^"]*' | grep -o '"token":"[^"]*' | cut -d'"' -f4 || echo "")
if [ -z "$API_KEY" ]; then
    API_KEY=$(echo "$API_KEY_RESPONSE" | grep -o '"token":"[^"]*' | cut -d'"' -f4 || echo "")
fi

if [ -z "$API_KEY" ]; then
    echo "Error: Failed to extract API key from response"
    echo "Response: $API_KEY_RESPONSE"
    exit 1
fi

echo "==> Test user and API key created successfully!"
echo "API Key: ${API_KEY:0:20}..."
echo ""
echo "Export these for tests:"
echo "export PEEKAPING_ENDPOINT=${API_URL}"
echo "export PEEKAPING_API_KEY=${API_KEY}"

# Write to a file for the Makefile to source
cat > .test-env <<EOF
export PEEKAPING_ENDPOINT=${API_URL}
export PEEKAPING_API_KEY=${API_KEY}
EOF

echo ".test-env file created with credentials"

