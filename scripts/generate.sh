TOKEN=$(curl -s -X POST "http://localhost:8080/rest/user/login" -d '{"username":"admin", "password":"admin"}' |jq -r '.access')


curl -s -X POST -H "Authorization: Bearer $TOKEN" http://localhost:8080/rest/code/generate -d '{"id":"KEDURmvLdQpRFY"}' | jq
