#!/bin/bash

# Test authentication flow for the CanIFly app
# This script simulates the authentication flow with Bearer tokens

BASE_URL="http://localhost:42423"

echo "Testing CanIFly Authentication Flow"
echo "==================================="

# Test 1: Check session without auth (should fail)
echo -e "\n1. Testing /api/session without authentication..."
curl -s -X GET "$BASE_URL/api/session" | jq '.' || echo "Failed"

# Test 2: Check session with invalid token (should fail)  
echo -e "\n2. Testing /api/session with invalid token..."
curl -s -X GET "$BASE_URL/api/session" \
  -H "Authorization: Bearer invalid-token" | jq '.' || echo "Failed"

# Test 3: Try to access protected endpoint without auth (should fail)
echo -e "\n3. Testing /api/accounts without authentication..."
curl -s -X GET "$BASE_URL/api/accounts" | jq '.' || echo "Failed"

# Test 4: Try to access protected endpoint with invalid token (should fail)
echo -e "\n4. Testing /api/accounts with invalid token..."
curl -s -X GET "$BASE_URL/api/accounts" \
  -H "Authorization: Bearer invalid-token" | jq '.' || echo "Failed"

# Test 5: Check health endpoint (should always work)
echo -e "\n5. Testing /health endpoint..."
curl -s -X GET "$BASE_URL/health" | jq '.' || echo "Failed"

echo -e "\n==================================="
echo "Authentication flow testing complete"
echo ""
echo "To test successful authentication:"
echo "1. Login through the app UI"
echo "2. Check browser DevTools for the session_token in localStorage"
echo "3. Use that token with: curl -H 'Authorization: Bearer <token>' $BASE_URL/api/accounts"