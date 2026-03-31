# Postpartum_BackEnd

Backend service untuk aplikasi digital postpartum care yang membantu pemantauan kesehatan ibu dan bayi pada masa nifas. Sistem ini mendukung autentikasi pengguna, pelacakan tidur bayi, pencatatan gejala harian ibu, alert kesehatan, serta booking konsultasi psikolog melalui WhatsApp.

## Fitur Utama

- autentikasi pengguna: register, login, refresh token, logout
- profil pengguna: lihat dan ubah nama pengguna
- pelacakan tidur bayi: start, end, input manual, input backdate, riwayat, insight, dan prediksi
- pencatatan gejala harian ibu: gejala fisik, pendarahan, kondisi emosional, dan alert harian
- konsultasi psikolog: daftar psikolog, jadwal, foto profil, dan booking ke WhatsApp

## Tech Stack

- Go
- Gin
- GORM
- MySQL
- Redis
- golang-migrate
- Docker Compose
- JWT Authentication

## Struktur Project

```text
cmd/api                 entry point aplikasi
config/                 konfigurasi aplikasi
internal/controller/    layer HTTP handler
internal/domain/        business rules dan domain logic
internal/dto/           request dan response DTO
internal/entity/        entity database
internal/repository/    akses data
internal/usecase/       application use case
pkg/                    package pendukung
```

## Prasyarat

Pastikan environment berikut tersedia:

- Go 1.24 atau versi yang kompatibel dengan `go.mod`
- MySQL 8
- Docker dan Docker Compose

## Konfigurasi Environment

Project menggunakan file `.env` untuk konfigurasi. Beberapa variabel utama yang perlu diperhatikan:

```env
APP_PORT=8080
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=postpartum_app
MYSQL_ROOT_PASSWORD=your_mysql_root_password
MYSQL_DATABASE=postpartum_app
MYSQL_HOST_PORT=3307
JWT_SECRET=your_secret
ADMIN_WA_NUMBER=628xxxxxxxxxx
REDIS_ADDR=127.0.0.1:6379
REDIS_PASSWORD=
REDIS_DB=0
```

Gunakan `.env.example` sebagai template awal, lalu salin menjadi `.env` dan isi dengan nilai environment Anda sendiri.

Keterangan:

- `DB_HOST` dan `DB_PORT` dipakai saat menjalankan aplikasi secara lokal
- untuk Docker Compose, service API akan otomatis menggunakan koneksi internal ke service `db`
- `ADMIN_WA_NUMBER` digunakan untuk membentuk link booking WhatsApp psikolog
- Redis bersifat opsional; jika `REDIS_ADDR` tidak diisi, aplikasi tetap berjalan tanpa cache

## Menjalankan Project

### Menjalankan Secara Lokal

Jalankan migration terlebih dahulu:

```bash
migrate -path migrations -database "mysql://root:your_password@tcp(127.0.0.1:3306)/postpartum_app" up
```

Jalankan seed data awal:

```bash
go run ./cmd/seed
```

Kemudian jalankan aplikasi:

```bash
go run ./cmd/api
```

Server akan berjalan di:

```text
http://localhost:8080
```

### Menjalankan dengan Docker

Jalankan database dan API:

```bash
docker compose up --build
```

Port default:

- API: `http://localhost:8080`
- MySQL dari host: `localhost:3307`
- Redis dari host: `localhost:6379`

Jika ingin menjalankan migration ke database Docker dari host:

```bash
migrate -path migrations -database "mysql://root:your_mysql_root_password@tcp(127.0.0.1:3307)/postpartum_app" up
```

## Database Migration

Project ini menggunakan `golang-migrate` untuk mengelola perubahan schema database secara versioned.

Lokasi file migration:

```text
migrations/
```

Contoh command yang umum dipakai:

```bash
migrate create -ext sql -dir migrations -seq migration_name
migrate -path migrations -database "mysql://root:your_password@tcp(127.0.0.1:3306)/postpartum_app" up
migrate -path migrations -database "mysql://root:your_password@tcp(127.0.0.1:3306)/postpartum_app" down 1
migrate -path migrations -database "mysql://root:your_password@tcp(127.0.0.1:3306)/postpartum_app" version
```

Contoh command untuk Windows jika binary `migrate` belum ada di `PATH`:

```powershell
C:\Users\USER\go\bin\migrate.exe -path migrations -database "mysql://root:your_password@tcp(127.0.0.1:3306)/postpartum_app" up
```

Catatan:

- aplikasi tidak lagi menjalankan `AutoMigrate` saat start
- schema database harus dibuat melalui migration
- seed data psikolog dan alert rule dijalankan melalui command terpisah

Urutan yang disarankan untuk Docker:

```powershell
docker compose up -d db
C:\Users\USER\go\bin\migrate.exe -path migrations -database "mysql://root:your_mysql_root_password@tcp(127.0.0.1:3307)/postpartum_app" up
go run ./cmd/seed
docker compose up --build -d api
```

## Endpoint Utama

### Auth

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/logout`

### User

- `GET /api/v1/user/profile`
- `PATCH /api/v1/user/profile`

### Sleep

- `POST /api/v1/sleep/start`
- `POST /api/v1/sleep/end`
- `POST /api/v1/sleep/manual`
- `POST /api/v1/sleep/bulk`
- `GET /api/v1/sleep/status`
- `GET /api/v1/sleep/daily`
- `GET /api/v1/sleep/history`
- `GET /api/v1/sleep/predict`
- `GET /api/v1/sleep/insight`

### Symptom

- `POST /api/v1/symptom`
- `GET /api/v1/symptom/history`
- `GET /api/v1/symptom/:date`

### Psychologist

- `GET /api/v1/psychologists`
- `GET /api/v1/psychologists/:id`
- `PATCH /api/v1/psychologists/:id/photo`
- `POST /api/v1/psychologists/:id/booking`

### Health Check

- `GET /health`

## Alur Singkat Fitur

### Tidur Bayi

- pengguna dapat mencatat maksimal 8 sesi tidur per hari
- input tanggal lampau hanya untuk tracking dan hanya dapat disimpan satu kali
- prediksi tidur hanya menggunakan riwayat valid sebelumnya, bukan data backdate
- insight hari ini menampilkan kondisi tidur bayi berdasarkan data hari ini dan prediksi yang tersedia

### Gejala Harian dan Alert

- gejala harian dapat diperbarui berkali-kali pada hari yang sama
- data tanggal lampau hanya dapat disimpan satu kali
- alert yang disimpan pada riwayat adalah hasil keputusan terakhir pada hari tersebut
- sistem menampilkan satu alert utama berdasarkan prioritas level dan gejala yang paling relevan

### Booking Psikolog

- pengguna memilih psikolog dan jadwal yang tersedia
- sistem membentuk template pesan booking
- pengguna diarahkan ke WhatsApp admin menggunakan link `wa.me`

### Cache

- daftar psikolog dicache dengan Redis
- detail psikolog dicache dengan Redis
- alert rules aktif dicache dengan Redis

## Testing

Pengujian endpoint dapat dilakukan menggunakan Postman. Disarankan untuk menguji alur berikut:

- register dan login
- get profile dan update profile
- sleep tracking, history, insight, dan predict
- symptom save, update same day, backdate, history, dan detail
- psychologist list, detail, dan booking
- logout

## Catatan

- field `birth_date` ditampilkan di profile, tetapi tidak diubah melalui endpoint profile update
- hasil prediksi tidur akan tersedia setelah data historis yang valid mencukupi
- endpoint yang membutuhkan autentikasi harus menggunakan Bearer token

## Status

Project ini sudah siap untuk integrasi ke frontend dan testing fitur utama, dengan fokus pada backend API dan business logic sesuai kebutuhan aplikasi postpartum care.
