# Golang

Simple Go project scaffold.

## Project Structure

```text
.
├── assets/
├── cmd/
│   └── app/
│       └── main.go
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal/
│   └── user/
│       ├── dto.go
│       ├── error.go
│       ├── handler.go
│       ├── repository.go
│       ├── service.go
│       └── user.go
├── pkg/
│   ├── config/
│   │   └── config.go
│   ├── constants/
│   │   └── constants.go
│   └── response/
│       └── response.go
├── scripts/
├── .env
├── .env.example
├── go.mod
├── go.sum
└── README.md
```

## API Endpoints

- `POST /api/users/register` - Register


## Development Commands

```bash
# Install dependencies
go mod tidy

# Generate API documentation
swag init -g cmd/app/main.go

# Start the server
go run cmd/app/main.go
```