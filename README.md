
## ⚡️ Quick start
### before project(if not working)
```bash
#zsh에서 cmd가 안깔린 경우에 실행이 안되는 경우 존재
export PATH=$(go env GOPATH)/bin:$PATH
go get -u github.com/swaggo/swag/cmd/swag
```

### start project
```bash
make docker.run
```
Go to API Docs page (Swagger): [127.0.0.1:3000/swagger/index.html](http://127.0.0.1:3000/swagger/index.html)

### format project
```bash
gofumpt -l -w .
```

## Configuration

```ini
# .env

# Stage status to start server:
#   - "dev", for start server without graceful shutdown
#   - "prod", for start server with graceful shutdown
STAGE_STATUS="dev"

# Server settings:
SERVER_HOST="0.0.0.0"
SERVER_PORT=8000
SERVER_READ_TIMEOUT=60

# JWT settings:
JWT_SECRET_KEY="secret"
JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT=15
JWT_REFRESH_KEY="refresh"
JWT_REFRESH_KEY_EXPIRE_HOURS_COUNT=720

# OCI SDK settings:
tenancyID=tenancy
userID=user
fingerprint=fingerprint
privateKeyFile=filePath
region=us-ashburn-1
compartmentID=compartmentID
privateKeyPass=

```

