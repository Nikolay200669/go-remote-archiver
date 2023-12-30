## for remote manager archive

### build for windows
```bash
GOARCH=amd64 GOOS=windows go build -o sysarch.exe -ldflags="-s -w"
```

```bash
curl -X POST http://localhost:8088/arch \
   -H 'Content-Type: application/json' \
   -d '{"password": "123456","save_to": "/Users/nik/Downloads/","catalog": "test"}' \
   --output /Users/nik/Downloads/test.zip
```