# GoTicket API 🎟️

**GoTicket API** adalah layanan RESTful API backend modern untuk sistem pemesanan tiket konser musik yang dibangun menggunakan **Go (Golang)**, **Gin Web Framework**, **GORM ORM**, dan **PostgreSQL**.

Aplikasi ini dirancang dengan arsitektur bersih (*Clean Architecture pattern*), mendukung proteksi *race condition* pada pemesanan tiket berkecepatan tinggi, integrasi OTP Email, ekspor laporan transaksi dalam format PDF, serta dokumentasi interaktif **OpenAPI / Swagger UI**.

---

## 🛠️ Teknologi & Stack

- **Bahasa Pemrograman**: Go (Golang) 1.20+
- **Web Framework**: [Gin Gonic](https://github.com/gin-gonic/gin)
- **Database ORM**: [GORM](https://gorm.io/) (Driver PostgreSQL)
- **Database**: PostgreSQL 13+
- **Autentikasi & Keamanan**: JWT (JSON Web Tokens) & Bcrypt Password Hashing
- **Dokumentasi API**: [Swaggo / Gin-Swagger](https://github.com/swaggo/gin-swagger)
- **Laporan PDF**: [gofpdf](https://github.com/jung-kurt/gofpdf)
- **Live Reload (Dev)**: [Air](https://github.com/air-verse/air)

---

## ✨ Fitur Utama

1. **Autentikasi & Otorisasi Lengkap**:
   - Registrasi Pengguna & Login dengan Token JWT.
   - Verifikasi **OTP Email** 6-digit saat registrasi akun.
   - Fitur **Lupa Password** (`Forgot Password` & `Reset Password` via OTP Email).
   - Invalidasi Token JWT saat **Logout** (Token Blacklisting).
   - Otorigasi Berbasis Peran (*Role-Based Access Control*): `admin` dan `customer`.
2. **Pemesanan Tiket & Proteksi Race Condition**:
   - Transaksi pemesanan tiket dengan penguncian baris PostgreSQL (`FOR UPDATE`) untuk mencegah kuota minus / *overselling* saat pemesanan bersamaan.
3. **Ekspor Laporan Transaksi PDF**:
   - Fitur *Generate & Download* Laporan Pemesanan Tiket Konser format PDF Landscape yang dapat difilter berdasarkan rentang tanggal (`start_date` dan `end_date`).
4. **Manajemen Konser & Upload Berkas**:
   - CRUD Konser musik mendukung *multipart/form-data* untuk mengunggah Gambar Thumbnail dan Berkas Aturan Konser (PDF).
5. **Health Check & Monitoring Database**:
   - Endpoint `/health` yang secara otomatis melakukan *ping check* ke koneksi PostgreSQL.
6. **Swagger UI Interaktif**:
   - Uji coba API langsung melalui antarmuka Swagger UI di `/swagger/index.html` (otomatis mengembalikan `404 Not Found` pada mode produksi).

---

## 📋 Struktur Direktori Proyek

```text
goticket-api/
├── config/             # Konfigurasi aplikasi & koneksi database PostgreSQL
├── docs/               # Berkas spesifikasi OpenAPI Swagger (generated)
├── dto/                # Data Transfer Objects (Request/Response structs)
├── handler/            # Layer HTTP Handlers (Gin Handlers)
├── middleware/         # Middleware JWT Auth, Role Enforcement, Recovery
├── model/              # Entity & Struct Model Database (GORM)
├── repository/         # Layer Akses Data PostgreSQL
├── routes/             # Registrasi Rute & Routing Engine Gin
├── service/            # Layer Business Logic & PDF Generator
├── uploads/            # Direktori penyimpanan media/berkas lokal
├── utils/              # Helper JWT & Pengirim Email SMTP
├── .env.example        # Template variabel environment
├── .gitignore          # Berkas pengabaian Git
├── go.mod / go.sum     # Berkas dependensi Go Modules
├── main.go             # Entrypoint utama aplikasi
└── README.md           # Dokumentasi proyek
```

---

## 🚀 Panduan Memulai (Setup & Run)

### 1. Prasyarat
Pastikan komputer Anda telah terinstal:
- **Go** (Versi 1.20 atau lebih baru)
- **PostgreSQL** Server (Port default `5433` atau `5432`)

### 2. Kloning Repositori & Salin Environment
```bash
# Salin template variabel environment
cp .env.example .env
```

Sesuaikan nilai pada berkas `.env`:
```env
APP_PORT=18080
APP_MODE=development
BASE_URL=http://localhost:18080

DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=p4ssw0rd
DB_NAME=goticketdb
DB_SSLMODE=disable
DB_TIMEZONE=Asia/Jakarta

JWT_SECRET=supersecretjwtkey123_dev

SMTP_HOST=sandbox.smtp.mailtrap.io
SMTP_PORT=587
SENDER_EMAIL=noreply@goticket.com
AUTH_PASSWORD=your_smtp_app_password
```

### 3. Jalankan Aplikasi
```bash
# Unduh seluruh dependensi Go
go mod tidy

# Jalankan server
go run main.go
```

Atau menggunakan **Air** untuk *Live Reload*:
```bash
air
```

Aplikasi akan berjalan pada port **18080** (`http://localhost:18080`).

---

## 📖 Dokumentasi API (Endpoints)

### 🔑 Autentikasi (`/api/v1/auth`)
| Method | Endpoint | Akses | Deskripsi |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/auth/register` | Publik | Registrasi pengguna baru (Mengirimkan OTP Email) |
| `POST` | `/api/v1/auth/verify-otp` | Publik | Verifikasi kode OTP 6-digit (Aktivasi Akun / Reset Pass) |
| `POST` | `/api/v1/auth/resend-otp` | Publik | Mengirim ulang kode OTP ke email |
| `POST` | `/api/v1/auth/forgot-password` | Publik | Permintaan OTP reset password |
| `POST` | `/api/v1/auth/reset-password` | Publik | Mengubah password dengan OTP yang valid |
| `POST` | `/api/v1/auth/login` | Publik | Login & dapatkan Token JWT Access |
| `POST` | `/api/v1/auth/logout` | Protected | Invalidasi token JWT (Blacklist) |
| `GET`  | `/api/v1/auth/me` | Protected | Ambil profil pengguna yang sedang aktif |

### 🎵 Konser (`/api/v1/concerts`)
| Method | Endpoint | Akses | Deskripsi |
| :--- | :--- | :--- | :--- |
| `GET`    | `/api/v1/concerts` | Publik | Ambil daftar konser (Mendukung Search & Paginasi) |
| `GET`    | `/api/v1/concerts/:id` | Publik | Detail konser berdasarkan ID |
| `POST`   | `/api/v1/concerts` | Admin | Tambah konser baru & upload biner thumbnail/rules PDF |
| `PUT`    | `/api/v1/concerts/:id` | Admin | Perbarui data konser |
| `DELETE` | `/api/v1/concerts/:id` | Admin | Hapus data konser |

### 🎟️ Pemesanan Tiket (`/api/v1/bookings`)
| Method | Endpoint | Akses | Deskripsi |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/bookings` | Protected | Pemesanan tiket konser (Proteksi Race Condition) |
| `GET`  | `/api/v1/bookings/:id` | Protected | Detail transaksi pemesanan (Akses dibatasi sesuai pemilik/admin) |
| `GET`  | `/api/v1/bookings/report/pdf` | Admin | **Export PDF Report** transaksi booking dengan filter `start_date` & `end_date` |

### 📊 Health Check & Swagger
| Method | Endpoint | Akses | Deskripsi |
| :--- | :--- | :--- | :--- |
| `GET` | `/health` | Publik | Database Ping Health Check |
| `GET` | `/swagger/index.html` | Dev Only | Antarmuka Dokumentasi Swagger UI (404 di Prod) |

---

## 👤 Akun Bawaan Seeding (Development)

Saat pertama kali database diinisialisasi, sistem secara otomatis melakukan seeding data default:

- **Admin Account**:
  - Email: `admin@goticket.com`
  - Password: `password123`
  - Peran: `admin`

---

## 📄 Lisensi

Hak Cipta &copy; 2026 **GoTicket API Team**. Proyek ini berlisensi Apache 2.0.
