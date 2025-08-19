# 🏢 KaryawanApp – Aplikasi Manajemen Karyawan (Refactored)

**KaryawanApp** adalah aplikasi web modern dan aman untuk mengelola data karyawan secara CRUD (Create, Read, Update, Delete). Aplikasi ini telah direfactor dengan arsitektur yang lebih baik dan menerapkan best practices pengembangan perangkat lunak.

## 🏗️ Arsitektur Aplikasi

### Backend (Go)
- **Framework**: Native HTTP server dengan Gorilla Mux untuk routing
- **Arsitektur**: Clean Architecture dengan pemisahan layer (Handler, Service, Repository)
- **Database**: MySQL dengan GORM sebagai ORM
- **Keamanan**: 
  - Input sanitization
  - Rate limiting (100 request/menut per IP)
  - CORS middleware
  - Validasi input

### Frontend (React + TypeScript)
- **Framework**: React dengan TypeScript
- **UI Library**: Material-UI (MUI)
- **State Management**: React Hooks (useState, useEffect, useContext)
- **Routing**: React Router v6
- **HTTP Client**: Axios untuk API calls

## 🛠️ Persyaratan Sistem

- Go 1.24+
- Node.js 16+
- MySQL 8.0+
- npm atau yarn

## 🚀 Panduan Instalasi

### 1. Setup Database

1. Buat database MySQL baru:
   ```sql
   CREATE DATABASE karyawan_db;
   ```

2. Buat file `.env` di root project dengan konfigurasi:
   ```env
   DB_HOST=localhost
   DB_PORT=3306
   DB_USER=root
   DB_PASSWORD=your_password
   DB_NAME=karyawan_db
   PORT=8080
   ```

### 2. Menjalankan Backend

1. Install dependencies:
   ```bash
   go mod tidy
   ```

2. Jalankan server:
   ```bash
   cd cmd/server
   go run main.go
   ```
   Server akan berjalan di `http://localhost:8080`

### 3. Menjalankan Frontend

1. Masuk ke direktori frontend:
   ```bash
   cd frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Jalankan development server:
   ```bash
   npm start
   ```
   Aplikasi akan terbuka di `http://localhost:3000`

## 🧪 Testing

### Backend Tests
```bash
go test ./...
```

### Frontend Tests
```bash
cd frontend
npm test
```

## 🚀 Fitur Utama

### ✅ **Manajemen Karyawan**
- Tambah karyawan baru
- Lihat daftar karyawan
- Edit data karyawan
- Hapus karyawan
- Pencarian dan filter

### ✅ **Keamanan & Validasi**
- **Input Sanitization**: Mencegah XSS attacks
- **Rate Limiting**: 100 requests per minute per IP
- **CORS Headers**: Cross-origin resource sharing
- **Form Validation**: Client-side dan server-side validation

### ✅ **User Experience**
- **Loading States**: Visual feedback saat loading
- **Error Handling**: Pesan error yang informatif
- **Responsive Design**: Mobile-friendly interface
- **Modern UI**: Menggunakan Material-UI untuk tampilan yang konsisten

## 📝 Dokumentasi API

### Daftar Karyawan
- **GET** `/api/employees` - Mendapatkan daftar semua karyawan
- **GET** `/api/employees/:id` - Mendapatkan detail karyawan
- **POST** `/api/employees` - Menambahkan karyawan baru
- **PUT** `/api/employees/:id` - Memperbarui data karyawan
- **DELETE** `/api/employees/:id` - Menghapus karyawan

## 🤝 Berkontribusi

1. Fork repository ini
2. Buat branch fitur baru (`git checkout -b fitur/namafitur`)
3. Commit perubahan Anda (`git commit -am 'Menambahkan fitur baru'`)
4. Push ke branch (`git push origin fitur/namafitur`)
5. Buat Pull Request

## 📄 Lisensi

Proyek ini dilisensikan di bawah MIT License - lihat file [LICENSE](LICENSE) untuk detailnya.
