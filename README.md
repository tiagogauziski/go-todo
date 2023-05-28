# TODO in Go
TODO RESTFul API in Go

## Development commands

- Run webserver
```bash
go run ./cmd/go-todo/go-todo.go
```

- Run tests
```bash
go test ./internal
```

- Run Dockerfile
```bash
# Build image
docker build -f ./build/go-todo/Dockerfile -t go-todo:latest .

# Run image
docker run -p 9000:9000 -e PORT=9000 -e "<MYSQL_DSN>" go-todo:latest
```

- Run docker-compose integrated tests
```bash
docker-compose up -f ./build/go-todo/docker-compose.test.yaml
```

## Environment Variables
```bash
# Database connection
DATABASE_URI=<USER>:<PASSWORD>@tcp(<HOSTNAME>:31835)/<DATABASE_NAME>?parseTime=true

# Database connection for running tests
DATABASE_TEST_URI=<USER>:<PASSWORD>@tcp(<HOSTNAME>:31835)/<DATABASE_NAME>?parseTime=true

# Gin port
PORT=9000
```