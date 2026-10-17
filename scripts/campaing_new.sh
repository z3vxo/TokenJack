TOKEN=$(curl -s -X POST "http://localhost:8080/rest/user/login" -d '{"username":"admin", "password":"admin"}' |jq -r '.access')

curl -X POST -H "Authorization: Bearer $TOKEN" "http://localhost:8080/rest/campaign/new" \
  -d '{"name":"test","tenant_id":"common","resources":{"name":"graph","scopes":[".default"]}}'
