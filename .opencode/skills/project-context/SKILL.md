---
name: project-context
description: Ngữ cảnh dự án SpendSnap (Go Clean Architecture, Locket-style chi tiêu + MXH). Use when working on SpendSnap backend/frontend tasks — database schema, API design, folder structure, roadmap/WBS, Redis cache, WebSocket chat, Cloudflare R2 upload, widget, Docker/CI-CD — để đảm bảo code tuân theo kiến trúc và thiết kế đã chốt.
---

# Project Context: SpendSnap

**Tên dự án:** SpendSnap - Chi tiêu cá nhân & Mạng xã hội (Locket Style)
**Kiến trúc:** Monolith (Go Clean Architecture), sẵn sàng tách Microservices.

## 1. Roadmap

| Giai đoạn | Tuần | Trọng tâm                                          | Deliverables                                            |
| --------- | ---- | -------------------------------------------------- | ------------------------------------------------------- |
| 1         | 1    | Thiết kế DB nền tảng (User, Expense, Friend, Chat) | Schema Postgres Cloud + bộ khung folder                 |
| 2         | 2    | Locket Media Pipeline xử lý/nén ảnh                | API upload, Goroutine nén .webp → Cloudflare R2         |
| 3         | 3    | Phân tích tài chính & danh mục chi tiêu            | CRUD chi tiêu, cache báo cáo qua Redis                  |
| 4         | 4    | Real-time Feed cho Widget màn hình khóa            | API Widget đọc từ Redis, đồng bộ ảnh/tiền               |
| 5         | 5    | MXH: Kết bạn & Chat/Bình luận real-time            | API kết bạn, WebSocket chat dưới bài đăng               |
| 6         | 6    | App mobile Frontend + Widget HĐH                   | React Native: Camera, Feed bạn bè, biểu đồ, widget khóa |
| 7         | 7    | Docker, CI/CD, deploy                              | Dockerfile, deploy Koyeb/Render, GitHub Actions         |

## 2. Database Schema

### `users`

| Trường         | Kiểu         | Ràng buộc                        |
| -------------- | ------------ | -------------------------------- |
| `id`           | UUID         | PK                               |
| `username`     | VARCHAR(50)  | NOT NULL                         |
| `email`        | VARCHAR(100) | UNIQUE, NOT NULL                 |
| `password`     | VARCHAR(255) | NOT NULL                         |
| `avatar_url`   | TEXT         |                                  |
| `phone_number` | VARCHAR(20)  | UNIQUE                           |
| `status`       | VARCHAR(20)  | NOT NULL — 'active' / 'inactive' |
| `created_at`   | TIMESTAMP    | NOT NULL DEFAULT NOW()           |

### `friends`

| Trường       | Kiểu        | Ràng buộc                         |
| ------------ | ----------- | --------------------------------- |
| `user_id_1`  | UUID        | FK(users.id), PK                  |
| `user_id_2`  | UUID        | FK(users.id), PK                  |
| `status`     | VARCHAR(20) | NOT NULL — 'pending' / 'accepted' |
| `updated_at` | TIMESTAMP   | NOT NULL                          |

### `personal_expenses`

| Trường        | Kiểu      | Ràng buộc                   |
| ------------- | --------- | --------------------------- |
| `id`          | UUID      | PK                          |
| `user_id`     | UUID      | FK(users.id), NOT NULL      |
| `category_id` | INT       | FK(categories.id), NOT NULL |
| `amount`      | BIGINT    | NOT NULL                    |
| `image_url`   | TEXT      | NOT NULL — ảnh locket       |
| `created_at`  | TIMESTAMP | NOT NULL DEFAULT NOW()      |

### `comments`

| Trường       | Kiểu      | Ràng buộc                          |
| ------------ | --------- | ---------------------------------- |
| `id`         | BIGSERIAL | PK                                 |
| `expense_id` | UUID      | FK(personal_expenses.id), NOT NULL |
| `user_id`    | UUID      | FK(users.id), NOT NULL             |
| `content`    | TEXT      | NOT NULL                           |
| `created_at` | TIMESTAMP | NOT NULL DEFAULT NOW()             |

## 3. Cấu trúc thư mục (Go Clean Architecture)

```text
spendsnap-backend/
├── cmd/server/main.go          # Entry point
├── configs/config.yaml         # Cấu hình app
├── deployments/                # Dockerfile, docker-compose
├── internal/
│   ├── config/                 # Load cấu hình
│   ├── delivery/               # Controller (HTTP REST, WebSocket)
│   ├── domain/                 # Models, Interfaces
│   ├── repository/             # Giao tiếp DB (Postgres, Redis)
│   └── usecase/                # Business logic
├── migrations/                 # SQL migrations
├── pkg/                        # Tiện ích dùng chung (database, cache, logger...)
├── scripts/                    # Script CI/CD
├── .env.example
├── Makefile
└── README.md
```

## 4. Quy ước khi làm việc với codebase này

- Luôn đặt code vào đúng layer: handler trong `delivery/`, business logic trong `usecase/`, truy cập DB/cache trong `repository/`, models/interfaces trong `domain/`.
- Tiện ích dùng chung (kết nối DB, Redis client, logger) đặt trong `pkg/`.
- Migration SQL đặt trong `migrations/`.
- Stack chính: Go backend, PostgreSQL (Postgres Cloud), Redis (Upstash) làm cache, Cloudflare R2 lưu ảnh (.webp, crop vuông, nén bằng Goroutine), WebSocket cho chat real-time, React Native cho mobile app.
