TOKEN=$(curl -s -X POST "http://localhost:8080/rest/user/login" -d '{"username":"admin", "password":"admin"}' |jq -r '.access')


curl -v -X POST "http://localhost:8080/rest/campaign/meta/t2QJpeozeF2UnP" -H "Authorization: Bearer $TOKEN"
