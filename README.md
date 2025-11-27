# 🎓 Student Achievement Reporting System – Backend API  
Backend untuk Sistem Pelaporan Prestasi Mahasiswa menggunakan **Golang Fiber**, **PostgreSQL**, **MongoDB**, dan **JWT RBAC**, serta menerapkan **Clean Architecture**.  
Project ini dibuat sebagai bagian dari **Ujian Akhir Semester (UAS)** Mata Kuliah *Pemrograman Backend Lanjut*.

## 🚀 Tech Stack
- Golang (Fiber Framework)
- Clean Architecture
- PostgreSQL + GORM
- MongoDB
- JWT Authentication
- RBAC
- Swagger Docs

## 📁 Project Structure
```
project/
│
├── app/
│   ├── models/
│   ├── mongo_models/
│   ├── repositories/
│   ├── services/
│   ├── controllers/
│
├── config/
├── routes/
├── middlewares/
├── main.go
```

## 📌 Fitur Utama
### Mahasiswa
- Input & edit prestasi
- Kirim prestasi untuk verifikasi

### Dosen Wali
- Lihat prestasi mahasiswa
- Verifikasi / tolak prestasi

### Admin
- Kelola user, role, permission
- Kelola referensi prestasi

## ⚙️ Cara Menjalankan
1. Clone repo  
2. `go mod tidy`  
3. Buat file `.env`  
4. `go run main.go`

## 👨‍💻 Author
Kenzie  
Project UAS Backend.
