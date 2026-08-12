# 🏢 KaryawanApp – Aplikasi Manajemen Karyawan

**KaryawanApp** adalah aplikasi web modern untuk mengelola data karyawan secara CRUD (Create, Read, Update, Delete) dengan arsitektur Clean Architecture dan **JWT Authentication**.

## 🔐 Default Login Credentials

Setelah setup, Anda dapat login dengan:

```
Username: admin
Password: admin123
```

**⚠️ PENTING**: Segera ubah password default setelah instalasi!

## 🏗️ Arsitektur Aplikasi

### Backend (Go)
- **Framework**: Gorilla Mux untuk routing
- **Arsitektur**: Clean Architecture (Handler → Service → Repository)
- **Database**: MySQL/MariaDB dengan raw SQL
- **Keamanan**:
  - JWT Authentication
  - Input sanitization (XSS prevention)
  - Rate limiting (100 requests/menit per IP)
  - CORS middleware
  - Input validation

### Frontend (React + TypeScript)
- **Framework**: React 19 dengan TypeScript
- **UI Library**: Material-UI (MUI) v7
- **Data Grid**: MUI X Data Grid
- **Routing**: React Router v6
- **HTTP Client**: Axios

## 🛠️ Persyaratan Sistem

- Go 1.24+
- Node.js 16+
- MySQL 8.0+ atau MariaDB 10.5+
- npm atau yarn

## 🚀 Panduan Instalasi

### **Quick Setup (Recommended)**

Jalankan script setup otomatis:

```bash
./setup.sh
```

Script ini akan:
1. Membuat database `karyawan_db`
2. Update file `.env` dengan konfigurasi yang benar
3. Generate JWT secret yang aman
4. Build aplikasi
5. Create default admin user (username: `admin`, password: `admin123`)

### **Manual Setup**

#### 1. Setup Database

```bash
# Login ke MySQL/MariaDB
mysql -u root -p

# Buat database
CREATE DATABASE karyawan_db;
EXIT;
```

#### 2. Konfigurasi Environment

File `.env` sudah dibuat dengan konfigurasi default:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=karyawan_db
PORT=8083
HOST=127.0.0.1
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60
JWT_SECRET=change-this-secret-key-in-production
JWT_EXPIRATION_HOURS=24
```

**⚠️ PENTING**: 
- Ubah `DB_PASSWORD` sesuai password MySQL Anda
- Ubah `JWT_SECRET` sebelum production!

#### 3. Jalankan Seed Command (Create Admin User)

```bash
go run cmd/seed/seed.go
```

Ini akan membuat default admin user:
- **Username**: `admin`
- **Password**: `admin123`

#### 4. Menjalankan Aplikasi

#### **Opsi A: Development Mode (Recommended)**

```bash
# Terminal 1 - Backend
make run
# atau
go run cmd/server/main.go

# Terminal 2 - Frontend (opsional, untuk development)
cd frontend
npm start
```

Akses:
- **Backend API**: http://127.0.0.1:8083
- **Frontend Dev**: http://localhost:3000

#### **Opsi B: Production Build**

```bash
# Build frontend
cd frontend
npm install
npm run build

# Build backend
cd ..
make build

# Jalankan aplikasi
./karyawan-app
```

Akses: **http://127.0.0.1:8083**

**Catatan**: Migrasi database dijalankan otomatis saat aplikasi start pertama kali!

## 📁 Struktur Project

```
Managemen_Karyawan_GO/
├── cmd/
│   └── server/          # Main entry point
├── config/              # Database configuration
├── internal/
│   ├── auth/           # JWT authentication
│   ├── domain/         # Interfaces & types
│   ├── handler/        # HTTP handlers
│   ├── middleware/     # Auth & role middleware
│   ├── models/         # Data models
│   ├── repository/     # Database layer
│   └── service/        # Business logic
├── frontend/           # React TypeScript app
├── assets/             # Static assets
├── .env                # Environment variables
├── .env.example        # Environment template
├── Makefile            # Build commands
└── go.mod              # Go dependencies
```

## 🧪 Testing

### Backend Tests
```bash
make test
# atau
go test -v ./...
```

### Frontend Tests
```bash
cd frontend
npm test
```

### Run dengan Coverage
```bash
make test-coverage
# Lihat hasil: coverage.html
```

## 🚀 Fitur Utama

### ✅ **Manajemen Karyawan**
- Tambah karyawan baru
- Lihat daftar karyawan dengan pagination
- Edit data karyawan
- Hapus karyawan (soft delete)
- Pencarian dan filter
- Export data (CSV, Print)

### ✅ **Keamanan & Validasi**
- **JWT Authentication**: Token-based auth
- **Input Sanitization**: Mencegah XSS attacks
- **Rate Limiting**: 100 requests per minute per IP
- **CORS Headers**: Cross-origin resource sharing
- **Form Validation**: Client-side dan server-side validation
- **Password Hashing**: bcrypt untuk password

### ✅ **User Experience**
- **Loading States**: Visual feedback saat loading
- **Error Handling**: Pesan error yang informatif
- **Responsive Design**: Mobile-friendly interface
- **Modern UI**: Material-UI untuk tampilan konsisten

## 📝 Dokumentasi API

### Authentication API

#### Login
```bash
POST /api/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "role": "admin",
    "is_active": true,
    "created_at": "2026-02-17T00:00:00Z"
  }
}
```

#### Get Profile (Protected)
```bash
GET /api/auth/profile
Authorization: Bearer <token>
```

### Health Check
```bash
GET /api/health
```

### Employees API (Protected - Requires Authentication)

Semua endpoint employees sekarang **memerlukan JWT token** untuk akses.

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/employees` | Get all employees |
| GET | `/api/employees/:id` | Get employee by ID |
| POST | `/api/employees` | Create new employee |
| PUT | `/api/employees/:id` | Update employee |
| DELETE | `/api/employees/:id` | Delete employee (soft) |

### Contoh Request dengan curl

```bash
# 1. Login untuk mendapatkan token
curl -X POST http://localhost:8083/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# Simpan token dari response

# 2. Gunakan token untuk akses protected API
TOKEN="your-jwt-token-here"

curl -X GET http://localhost:8083/api/employees \
  -H "Authorization: Bearer $TOKEN"

# 3. Create employee
curl -X POST http://localhost:8083/api/employees \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "position": "Developer",
    "role": "Engineer",
    "phone": "08123456789",
    "alamat": "Jl. Contoh No. 123, Jakarta"
  }'
```

## 🔧 Commands (Makefile)

```bash
make help          # Show all commands
make build         # Build backend
make run           # Run backend
make test          # Run tests
make test-coverage # Run tests with coverage
make clean         # Clean build artifacts
make deps          # Install dependencies
make lint          # Run linter
make bench         # Run benchmarks
```

## 🐛 Troubleshooting

### Database Connection Failed
```bash
# Pastikan MySQL/MariaDB running
sudo systemctl start mariadb
# atau
sudo systemctl start mysql

# Cek koneksi
mysql -u root -p -e "SHOW DATABASES;"
```

### Port Already in Use
```bash
# Ubah port di .env
PORT=8084
HOST=127.0.0.1
```

### Frontend Build Error
```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
npm run build
```

## 🤝 Berkontribusi

1. Fork repository ini
2. Buat branch fitur baru (`git checkout -b fitur/namafitur`)
3. Commit perubahan Anda (`git commit -am 'Menambahkan fitur baru'`)
4. Push ke branch (`git push origin fitur/namafitur`)
5. Buat Pull Request

## 📄 Lisensi

Proyek ini dilisensikan di bawah MIT License.

## 📞 Support

Untuk pertanyaan atau masalah, silakan buat issue di repository ini.
