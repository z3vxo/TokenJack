TOKEN=$(curl -s -X POST "http://localhost:8080/rest/user/login" -d '{"username":"admin", "password":"admin"}' |jq -r '.access')

TOKEN=$(curl -s -X POST "http://localhost:8080/rest/user/login" -d '{"username":"admin", "password":"admin"}' |jq -r '.access')

ID=$(curl -X POST -H "Authorization: Bearer $TOKEN" "http://localhost:8080/rest/campaign/new" \
  -d '{"name":"test","tenant_id":"9ec43e1d-d482-4c28-a82c-5c6a2771ec2a","resources":{"name":"graph","scopes":[".default"]}}' | jq -r '.id')

echo "$ID"

jq -n --arg id "$ID" '{id:$id}' | curl -s -X POST -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/rest/code/generate -d @-
