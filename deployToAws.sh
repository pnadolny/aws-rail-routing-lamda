GOOS=linux go build main.go
zip -j main.zip main
aws lambda update-function-code \
    --function-name  hackathon \
    --zip-file fileb://main.zip
