# SpendSnap Backend

Backend cho ứng dụng SpendSnap — Chi tiêu cá nhân & Mạng xã hội (Locket Style).

Kiến trúc: Go Clean Architecture (Monolith), sẵn sàng tách Microservices.

## Cấu trúc

```text
spendsnap-backend/
├── cmd/server/        # Entry point
├── configs/           # config.yaml
├── deployments/       # Dockerfile, docker-compose
├── internal/
│   ├── config/        # Load cấu hình
│   ├── delivery/      # Controller (http/, websocket/)
│   ├── domain/        # Models, Interfaces
│   ├── repository/    # Postgres, Redis
│   └── usecase/       # Business logic
├── migrations/        # SQL migrations
├── pkg/               # database, cache, logger...
└── scripts/           # CI/CD scripts
```

## Chạy nhanh

```bash
cp .env.example .env
make run
```

Kiểm tra sức khỏe: `GET http://localhost:8080/health`
