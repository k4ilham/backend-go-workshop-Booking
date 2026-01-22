# Online Booking API (Go, Fiber, Hexagonal Architecture)

Proyek backend untuk aplikasi Online Booking UMKM dengan arsitektur Hexagonal. HTTP adapter menggunakan Fiber, database utama PostgreSQL, logging aktivitas menggunakan Turso (SQLite HTTP API), dan notifikasi menggunakan webhook n8n.

## Fitur
- Admin Authentication (login email/password, bcrypt, JWT)
- Booking System (create + list terbaru dulu, status default "pending")
- Admin Dashboard (total booking hari ini + latest 10 bookings)
- Service Management (create, delete, list aktif)
- Notification (webhook POST ke n8n saat booking dibuat)
- Activity Logging (kirim log ke Turso saat login/booking dibuat)

## Arsitektur
- Domain: model entitas (User, Service, Booking)
- Ports: interface untuk repositori, logger, notifier
- Usecases: logika aplikasi (auth, booking, dashboard, services)
- Adapters:
  - HTTP (Fiber): routing, middleware JWT
  - PostgreSQL: repositori pengguna, layanan, booking
  - Turso: logger HTTP API
  - n8n: webhook notifier

Referensi kode:
- Entrypoint: [main.go](file:///c:/inercorp/project/workshop%203%20Inercorp/be-golang/cmd/server/main.go)
- Wiring aplikasi: [app.go](file:///c:/inercorp/project/workshop%203%20Inercorp/be-golang/internal/app/app.go)
- Konfigurasi: [config.go](file:///c:/inercorp/project/workshop%203%20Inercorp/be-golang/internal/config/config.go)
- Router & JWT middleware: [router.go](file:///c:/inercorp/project/workshop%203%20Inercorp/be-golang/internal/adapter/http/fiber/router.go)
- Repos PostgreSQL: [postgres.go](file:///c:/inercorp/project/workshop%203%20Inercorp/be-golang/internal/adapter/repository/postgres/postgres.go)
- Logger Turso: [turso.go](file:///c:/inercorp/project/workshop%203%20Inercorp/be-golang/internal/adapter/logger/turso/turso.go)
- Notifier n8n: [webhook.go](file:///c:/inercorp/project/workshop%203%20Inercorp/be-golang/internal/adapter/notification/n8n/webhook.go)
- Usecases: [auth_login.go](file:///c:/inercorp/project/workshop%203%20Inercorp/be-golang/internal/usecase/auth_login.go), [booking_create.go](file:///c:/inercorp/project/workshop%203%20Inercorp/be-golang/internal/usecase/booking_create.go), [booking_list.go](file:///c:/inercorp/project/workshop%203%20Inercorp/be-golang/internal/usecase/booking_list.go), [dashboard_stats.go](file:///c:/inercorp/project/workshop%203%20Inercorp/be-golang/internal/usecase/dashboard_stats.go), [service_manage.go](file:///c:/inercorp/project/workshop%203%20Inercorp/be-golang/internal/usecase/service_manage.go)

## Prasyarat
- Go 1.21+
- PostgreSQL siap jalan
- Turso HTTP API endpoint + token
- n8n workflow dengan HTTP Trigger

## Konfigurasi Lingkungan
Isi file `.env` di root proyek:

```
SERVER_ADDR=:8080
TOKEN_TTL=86400
POSTGRES_DSN=postgres://postgres:postgres@localhost:5432/bookingdb?sslmode=disable
JWT_SECRET=changeme-supersecret-jwt
TURSO_URL=https://your-turso-host/v2/execute
TURSO_TOKEN=changeme-turso-token
N8N_WEBHOOK_URL=http://localhost:5678/webhook/booking
```

Server otomatis membaca `.env` saat start. File `.gitignore` sudah mengabaikan `.env`.

## Migrasi Database
PostgreSQL:

```sql
CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS services (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  price INT NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS bookings (
  id SERIAL PRIMARY KEY,
  customer_name TEXT NOT NULL,
  customer_phone TEXT NOT NULL,
  service_id INT NOT NULL REFERENCES services(id),
  booking_date DATE NOT NULL,
  booking_time TEXT NOT NULL,
  status TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

Buat admin user:

```go
package main
import ("fmt"; "golang.org/x/crypto/bcrypt")
func main(){h,_:=bcrypt.GenerateFromPassword([]byte("secret"),bcrypt.DefaultCost);fmt.Println(string(h))}
```

Jalankan, ambil hasilnya sebagai `password_hash`, lalu:

```sql
INSERT INTO users(email, password_hash) VALUES ('admin@example.com', '$2a$10$zFgXR6CwvA.i17khfIhj6u.XV9xh.dncpo74hhRbPZeeztHGeO8Nu');
```

Turso (SQLite via HTTP API):

```json
{
  "statements": [
    { "sql": "CREATE TABLE IF NOT EXISTS activity_logs (id INTEGER PRIMARY KEY, action TEXT, detail TEXT, created_at TEXT)" }
  ]
}
```

Kirim payload di atas ke `TURSO_URL` dengan header `Authorization: Bearer <TURSO_TOKEN>`.

## Menjalankan

```bash
go mod tidy
go run ./cmd/server
```

Server mendengarkan di `SERVER_ADDR` (default `:8080`).

## Endpoint
- POST /admin/login
- POST /bookings
- GET /bookings
- GET /admin/dashboard (JWT)
- POST /services (JWT)
- DELETE /services/:id (JWT)
- GET /services (JWT)

## Contoh Request
Login:

```bash
curl -X POST http://localhost:8080/admin/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"secret"}'
```

Buat booking:

```bash
curl -X POST http://localhost:8080/bookings \
  -H "Content-Type: application/json" \
  -d '{"customer_name":"Jane","customer_phone":"08123456789","service_id":1,"booking_date":"2026-01-22","booking_time":"10:30"}'
```

Dashboard:

```bash
curl -H "Authorization: Bearer <JWT>" http://localhost:8080/admin/dashboard
```

Service create:

```bash
curl -X POST http://localhost:8080/services \
  -H "Authorization: Bearer <JWT>" -H "Content-Type: application/json" \
  -d '{"name":"Haircut","price":150000,"is_active":true}'
```

## Catatan
- N8N_WEBHOOK_URL harus mengarah ke workflow HTTP Trigger.
- Turso endpoint `TURSO_URL` mengikuti API execute; token diperlukan jika disetup.
- JWT menggunakan HS256 dengan `JWT_SECRET`. Simpan rahasia di environment, jangan commit.

