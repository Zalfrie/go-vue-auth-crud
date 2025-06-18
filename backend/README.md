#Backend Setup

cd backend
cp .env.example .env
# edit .env for MySQL, email creds, JWT secret
go mod download
# run migrations & seed
go run cmd/main.go migrate
go run cmd/main.go seed
# generate swagger docs
swag init -g cmd/main.go
# start server
go run cmd/main.go serve
# test
go test ./tests/...

Swagger UI: http://localhost:8080/docs/index.html