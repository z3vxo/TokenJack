#!/usr/bin/env bash
set -euo pipefail

BASE=http://localhost:8080

TOKEN=$(curl -s -X POST "$BASE/rest/user/login" \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"admin"}' | jq -r '.access')

echo "token: ${TOKEN:0:24}..."

curl -s "$BASE/rest/campaigns/list" \
  -H "Authorization: Bearer $TOKEN" | jq
