# Core Project Structure

[![Go Version](https://img.shields.io/badge/Go-1.23.0-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

Đây là cấu trúc thư mục core cho các dự án Go microservice của tôi. Dự án này cung cấp một foundation hoàn chỉnh với các thành phần cơ bản cần thiết cho việc phát triển microservice.

## 🚀 Tính năng

### Core Infrastructure
- **gRPC Server**: Cấu hình server gRPC với Protocol Buffers
- **Database Layer**: PostgreSQL với go-pg ORM và migrations
- **Authentication**: JWT token management
- **Logging**: Structured logging với Zap
- **Configuration**: Viper configuration management
- **Validation**: Input validation với protobuf-validate
- **Caching**: Redis integration (tùy chọn)
- **ID Generation**: NanoID cho unique identifiers

### Project Structure
- **Clean Architecture**: Tổ chức theo domain-driven design
- **Dependency Injection**: Bootstrap pattern cho dependency management
- **Repository Pattern**: Abstraction layer cho data access
- **Service Layer**: Business logic separation
- **Use Case Pattern**: Application business rules

## 🛠️ Công nghệ sử dụng

- **Backend**: Go 1.23.0
- **Framework**: gRPC với Protocol Buffers
- **Database**: PostgreSQL với go-pg ORM
- **Authentication**: JWT (JSON Web Tokens)
- **Password Hashing**: Argon2id
- **Configuration**: Viper
- **Logging**: Zap
- **Validation**: protobuf-validate
- **ID Generation**: NanoID
- **Caching**: Redis (tùy chọn)

## 📋 Yêu cầu hệ thống

- Go 1.23.0 hoặc cao hơn
- PostgreSQL 12.0 hoặc cao hơn
- Redis (tùy chọn, cho caching)
- Docker (tùy chọn)

## 🚀 Cài đặt

### 1. Clone repository

```bash
git clone https://github.com/anhvanhoa/stuct-core.git
cd core
```

### 2. Cài đặt dependencies

```bash
go mod download
```

### 3. Cài đặt công cụ cần thiết

```bash
# Cài đặt migrate CLI
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Cài đặt buf CLI
go install github.com/bufbuild/buf/cmd/buf@latest

# Cài đặt protoc-gen-go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### 4. Cấu hình môi trường

Tạo file `dev.config.yaml` dựa trên template:

```yaml
# Database Configuration
database:
  host: localhost
  port: 5432
  user: postgres
  password: 123456
  name: core_db
  sslmode: disable

# Redis Configuration (optional)
redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

# JWT Configuration
jwt:
  secret_key: "your-secret-key-here"
  access_token_expiry: 15m
  refresh_token_expiry: 7d

# Server Configuration
server:
  grpc_port: 50051
  http_port: 8080

# Logging
logging:
  level: debug
  format: json
```

### 5. Khởi tạo cơ sở dữ liệu

```bash
# Tạo database
make dev-create-db

# Chạy migrations
make migrate-dev-up
```

### 6. Chạy ứng dụng

```bash
# Build ứng dụng
make build

# Chạy ứng dụng
make run
```

Hoặc chạy trực tiếp:

```bash
go run cmd/main.go
```

## 📖 Sử dụng

### API Endpoints (gRPC)

Dự án sử dụng gRPC với các service được định nghĩa trong `proto/exam/v1/exam.proto`:

#### AuthService

- `CheckToken` - Kiểm tra tính hợp lệ của token
- `Login` - Đăng nhập người dùng
- `Register` - Đăng ký tài khoản mới
- `RefreshToken` - Làm mới access token
- `Logout` - Đăng xuất
- `VerifyAccount` - Xác thực tài khoản
- `ForgotPassword` - Yêu cầu đặt lại mật khẩu
- `ResetPasswordByCode` - Đặt lại mật khẩu bằng mã
- `ResetPasswordByToken` - Đặt lại mật khẩu bằng token
- `CheckCode` - Kiểm tra mã xác thực

### Ví dụ sử dụng với gRPC client

```go
package main

import (
    "context"
    "log"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    pb "your-project/proto/gen/exam/v1"
)

func main() {
    conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()

    client := pb.NewAuthServiceClient(conn)

    // Đăng nhập
    loginResp, err := client.Login(context.Background(), &pb.LoginRequest{
        EmailOrPhone: "user@example.com",
        Password:     "password123",
        Os:           "web",
    })
    if err != nil {
        log.Fatalf("Login failed: %v", err)
    }

    log.Printf("Login successful: %s", loginResp.Message)
}
```

## 🏗️ Cấu trúc dự án

```
core/
├── bootstrap/          # Khởi tạo ứng dụng
│   ├── app.go         # Cấu hình ứng dụng chính
│   ├── database.go    # Kết nối database
│   ├── env.go         # Cấu hình môi trường
│   └── redis.go       # Kết nối Redis
├── cmd/               # Entry point
│   └── main.go        # Main function
├── domain/            # Domain layer
│   ├── common/        # Common utilities
│   │   ├── pagination.go  # Pagination helpers
│   │   └── status.go      # Status constants
│   ├── entity/        # Domain entities
│   ├── repository/    # Repository interfaces
│   │   └── transaction.go # Transaction management
│   ├── service/       # Domain services
│   │   ├── argon/     # Password hashing
│   │   ├── cache/     # Caching utilities
│   │   ├── error/     # Error handling
│   │   ├── goid/      # ID generation
│   │   ├── jwt/       # JWT management
│   │   └── logger/    # Logging utilities
│   └── usecase/       # Use cases
├── infrastructure/    # Infrastructure layer
│   ├── grpc_service/  # gRPC server implementation
│   │   ├── base.go    # Base gRPC service
│   │   ├── logging.go # gRPC logging middleware
│   │   ├── sayhi.go   # Example service
│   │   ├── server.go  # gRPC server setup
│   │   └── validation_helpers.go # Validation utilities
│   ├── repo/          # Repository implementations
│   │   ├── base.go    # Base repository
│   │   └── transaction.go # Transaction implementation
│   └── service/       # External service implementations
│       ├── argon/     # Argon2 implementation
│       ├── database/  # Database hooks
│       ├── goid/      # ID generation service
│       ├── jwt/       # JWT service implementation
│       └── logger/    # Logger service
├── migrations/        # Database migrations
├── proto/             # Protocol Buffers definitions
│   ├── buf/           # Buf configuration
│   │   └── validate/  # Validation rules
│   ├── exam/          # API definitions
│   │   └── v1/        # Version 1 API
│   └── gen/           # Generated Go code
├── logs/              # Log files
├── go.mod             # Go modules
├── go.sum             # Go modules checksum
├── MakeFile           # Build scripts
├── buf.gen.yaml       # Buf generation config
├── buf.yaml           # Buf configuration
├── dev.config.yaml    # Development configuration
└── README.md          # This file
```

## 🔧 Development

### Tạo migration mới

```bash
make migrate-dev-create name=migration_name
```

### Chạy migrations

```bash
# Chạy tất cả migrations
make migrate-dev-up

# Rollback migration cuối
make migrate-dev-down

# Reset database
make migrate-dev-reset

# Drop database
make migrate-dev-drop
```

### Generate Protocol Buffers

```bash
# Generate Go code từ proto files
buf generate
```

### Chạy tests

```bash
go test ./...
```

### Linting và formatting

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run
```

## 🐳 Docker

### Build Docker image

```bash
docker build -t core-service .
```

### Chạy với Docker Compose

Tạo file `docker-compose.yml`:

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "50051:50051"
    environment:
      - CONFIG_FILE=dev.config.yaml
    depends_on:
      - postgres
      - redis

  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: core_db
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: 123456
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

volumes:
  postgres_data:
```

Chạy:

```bash
docker-compose up -d
```

## 📝 Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `CONFIG_FILE` | Path to configuration file | `dev.config.yaml` |
| `GRPC_PORT` | gRPC server port | `50051` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | `123456` |
| `DB_NAME` | Database name | `core_db` |
| `JWT_SECRET` | JWT secret key | `your-secret-key` |

## 🤝 Đóng góp

Chúng tôi rất hoan nghênh mọi đóng góp! Vui lòng đọc [CONTRIBUTING.md](CONTRIBUTING.md) để biết thêm chi tiết.

### Quy trình đóng góp

1. Fork dự án
2. Tạo feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit thay đổi (`git commit -m 'Add some AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Mở Pull Request

## 📄 License

Dự án này được phân phối dưới giấy phép MIT. Xem file [LICENSE](LICENSE) để biết thêm chi tiết.

## 🆘 Hỗ trợ

Nếu bạn gặp vấn đề hoặc có câu hỏi:

- Tạo [Issue](https://github.com/yourusername/core/issues) trên GitHub
- Liên hệ qua email: your-email@example.com
- Tham gia [Discussions](https://github.com/yourusername/core/discussions)

## 🙏 Acknowledgments

- [gRPC](https://grpc.io/) - High-performance RPC framework
- [go-pg](https://github.com/go-pg/pg) - PostgreSQL ORM for Go
- [Zap](https://github.com/uber-go/zap) - Fast, structured logging
- [Viper](https://github.com/spf13/viper) - Configuration solution for Go applications
- [Buf](https://buf.build/) - Modern Protocol Buffers toolkit

---

⭐ Nếu dự án này hữu ích, hãy cho chúng tôi một star trên GitHub!
