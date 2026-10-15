TOKEN=$(curl -s -X POST "http://localhost:8080/rest/user/login" -d '{"username":"admin", "password":"admin"}' |jq -r '.access')
echo $TOKEN
curl -s http://localhost:8080/rest/user/me -H "Authorization: bearer $TOKEN"
