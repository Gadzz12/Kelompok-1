# 🚀 Flowspace — Modern To-Do List & Team Workspace

<p align="center">
  <img src="static/logo.png" alt="Flowspace Logo" width="96" height="96" style="border-radius: 18px; box-shadow: 0 10px 25px rgba(0,0,0,0.2);">
  <br>
  <b>Aplikasi manajemen tugas dan ruang kerja kolaborasi tim yang cepat, responsif, dan siap di-hosting ke Cloudflare.</b>
  <br>
  <i>Dibangun dengan Go (Golang), GORM, PostgreSQL (Neon Cloud), HTML5, CSS3, dan Vanilla JavaScript.</i>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Version">
  <img src="https://img.shields.io/badge/Database-PostgreSQL-336791?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL">
  <img src="https://img.shields.io/badge/Cloud-Neon.tech-00E599?style=for-the-badge&logo=neon&logoColor=black" alt="Neon">
  <img src="https://img.shields.io/badge/Proxy-Cloudflare-F38020?style=for-the-badge&logo=cloudflare&logoColor=white" alt="Cloudflare">
  <img src="https://img.shields.io/badge/Design-Responsive-blueviolet?style=for-the-badge" alt="Responsive">
</p>

---

## 📱 Dukungan Multi-Perangkat (Multi-Device Support)

Aplikasi ini dirancang dengan prinsip **Mobile-First & Fluid Responsive Design**, bekerja optimal di berbagai ukuran layar:

| Perangkat | Resolusi Layar | Karakteristik Tampilan |
| :--- | :--- | :--- |
| 📱 **Mobile Phone** | `< 520px` (iPhone, Android) | Sidebar tersembunyi dengan tombol drawer menu, navigasi touch-friendly, tabel bergeser halus (horizontal touch-scroll), tombol aksi penuh (full-width), modal fleksibel. |
| 📱 **Phablet / Tablet Portabel** | `520px - 768px` | Kalender nasional dengan touch-scroll mulus, layout formulir responsif, kartu anggota tim adaptif. |
| 💻 **Tablet / Small Laptop** | `769px - 1024px` | Drawer navigasi off-canvas dengan backdrop overlay (`scrim`), statistik grid 2 kolom, tata letak konten efisien. |
| 💻 **Laptop / Desktop** | `1025px - 1440px` | Sidebar navigasi tetap (pinned), overview statistik lengkap (4 kolom), tabel fixed dengan aksi cepat. |
| 🖥️ **Desktop / Ultra-Wide** | `> 1440px` | Ruang kerja luas, layout proporsional tanpa distorsi visual. |

---

## ✨ Fitur Utama

- 📋 **Manajemen Tugas Lengkap (CRUD To-Do)**:
  - Buat tugas baru dengan judul, deskripsi, batas waktu (deadline), tingkat prioritas (*Low*, *Medium*, *High*), dan status (*Belum Dimulai*, *Sedang Berjalan*, *Selesai*).
  - Kelola *subtasks* (tugas turunan) interaktif dengan checklist langsung.
- 👥 **Kolaborasi Tim & Workspace**:
  - Undang anggota tim (*Collaborators*) ke tugas tertentu melalui username atau email.
  - Kelola daftar pertemanan (*Team & Friends*) untuk kolaborasi cepat via chip picker.
  - Lampirkan dokumen proyek (*Documents Attachment*) seperti PDF, DOC, Sheet.
- 📬 **Pusat Notifikasi & Kotak Masuk (Inbox)**:
  - Terima undangan kolaborasi, setujui (*Accept*) atau tolak (*Reject*) secara instan.
  - Tandai seluruh notifikasi terbaca dengan satu klik.
- 📅 **Kalender Nasional Indonesia & Jam WIB**:
  - Tampilan kalender bulanan interaktif terintegrasi jadwal tugas.
  - Penandaan hari libur nasional Indonesia dan jam digital real-time Waktu Indonesia Barat (WIB).
- 🎨 **Avatar Builder & Profil Pengguna**:
  - Kustomisasi avatar SVG interaktif: pilihan gender, warna kulit (*skin tone*), warna baju, warna rambut, dan kacamata (*accessories*).
  - Tampilan live preview langsung tersinkronisasi ke seluruh bagian dashboard.
- 🌓 **Tema Gelap & Terang (Dark / Light Mode)**:
  - Deteksi tema preferensi sistem operasi secara otomatis atau pemilihan manual via toggle yang tersimpan di `localStorage`.
- 🌐 **Cloud & Cloudflare Ready**:
  - Aset HTML template dan file statis dikompilasi menggunakan `embed.FS` standar Go menjadi binary tunggal mandiri.
  - Skrip otomasi Cloudflare Quick Tunnel untuk akses publik HTTPS instan tanpa localhost.
  - Multi-stage `Dockerfile` ringan (~20MB) siap deploy ke platform cloud (Render, Koyeb, Railway, VPS).

---

## 📁 Struktur Direktori

```text
todo-list/
├── cmd/
│   ├── drop-table/              # Utilitas drop tabel untuk pengujian
│   ├── seed/                    # Utilitas seeding data awal
│   └── server/                  # Alternatif entry point server
├── internal/
│   ├── app/                     # Lifecycle aplikasi & inisialisasi server
│   ├── config/                  # Pengelolaan environment variable & .env
│   ├── database/                # Koneksi GORM PostgreSQL & automigrasi
│   ├── handler/                 # HTTP controller/handlers & routing
│   ├── model/                   # Definisi struct & model database
│   └── repository/              # Data access layer (GORM & JSON fallback)
├── static/
│   ├── app.css                  # CSS utama antarmuka workspace & responsive media queries
│   ├── app.js                   # Interaktivitas JavaScript (Avatar, filter, drawer, tema)
│   ├── landing.css              # CSS khusus halaman landing page
│   └── logo.png                 # Logo aset Flowspace
├── templates/                   # Template HTML (Go html/template)
│   ├── create_todo.html         # Halaman form buat tugas baru
│   ├── dashboard.html           # Dashboard utama (Overview, Inbox, Team, Docs, Kalender)
│   ├── delete_todo.html         # Konfirmasi penghapusan tugas
│   ├── detail_todo.html         # Detail proyek, subtask, dan kolaborator
│   ├── edit_todo.html           # Formulir edit tugas
│   ├── landing.html             # Landing page publik
│   ├── login.html               # Halaman masuk (Sign In)
│   ├── register.html            # Halaman pendaftaran akun (Sign Up)
│   └── settings.html            # Pengaturan profil, avatar builder, & preferensi
├── .dockerignore
├── .env.example                 # Contoh konfigurasi environment
├── Dockerfile                   # Multi-stage Docker build untuk deployment cloud
├── go.mod                       # Modul Go & dependencies
├── go.sum                       # Checksum dependencies
├── main.go                      # Entry point utama dengan go:embed
├── main_test.go                 # Pengujian integritas embedded assets
├── PANDUAN_HOSTING_CLOUDFLARE.md# Panduan lengkap hosting via Cloudflare
├── README.md                    # Dokumentasi proyek
├── start-cloudflare-tunnel.bat  # 1-klik Cloudflare Tunnel untuk Windows
└── start-cloudflare-tunnel.ps1  # PowerShell runner untuk Cloudflare Tunnel
```

---

## 🛠️ Persyaratan Sistem

- **Go**: Versi `1.21` atau lebih baru ([Unduh Go](https://go.dev/dl/))
- **Database**: PostgreSQL (direkomendasikan layanan cloud gratis seperti [Neon.tech](https://neon.tech) atau [Supabase](https://supabase.com))
- **Koneksi Internet**: Untuk mengakses database cloud Neon dan mengunduh dependencies

---

## 🚀 Cara Menjalankan secara Lokal

### 1. Salin Pengaturan Lingkungan (.env)
Salin file `.env.example` menjadi `.env`:
```powershell
cp .env.example .env
```
Buka file `.env` dan masukkan connection string database PostgreSQL Anda:
```env
PORT=8080
DATABASE_URL=postgresql://USER:PASSWORD@HOST:5432/DATABASE?sslmode=require
```

### 2. Unduh Dependency
```powershell
go mod download
```

### 3. Jalankan Aplikasi
```powershell
go run main.go
```
Buka peramban (browser) di alamat: [http://localhost:8080](http://localhost:8080)

---

## ☁️ Menjalankan lewat Cloudflare (Tanpa Localhost)

Agar aplikasi dapat diakses dari luar jaringan (handphone, tablet, teman, klien) tanpa mengetik `localhost:8080`:

### Cara Instan: Cloudflare Quick Tunnel (1 Menit, 100% Gratis)
1. Jalankan aplikasi Go:
   ```powershell
   go run main.go
   ```
2. Klik ganda file:
   ```text
   start-cloudflare-tunnel.bat
   ```
   *(atau jalankan `.\start-cloudflare-tunnel.ps1` di PowerShell)*
3. Skrip akan otomatis mengunduh `cloudflared.exe` dan menampilkan URL HTTPS publik resmi dari Cloudflare:
   ```text
   +--------------------------------------------------------------------------------------------+
   |  Your quick Tunnel has been created! Visit it at:                                          |
   |  https://bebas-acak-nama.trycloudflare.com                                                |
   +--------------------------------------------------------------------------------------------+
   ```
4. Buka URL tersebut di perangkat mana saja di seluruh dunia!

> 💡 **Panduan Lengkap Cloudflare**: Untuk petunjuk menghubungkan ke nama domain sendiri (`todo.domainanda.com`) atau deploy 24 jam nonstop via Render/Koyeb, silakan baca [**PANDUAN_HOSTING_CLOUDFLARE.md**](PANDUAN_HOSTING_CLOUDFLARE.md).

---

## 🐳 Menjalankan dengan Docker

Aplikasi telah dilengkapi `Dockerfile` multi-stage:

```powershell
# Build image
docker build -t flowspace-app .

# Jalankan container
docker run -p 8080:8080 --env-file .env flowspace-app
```

---

## 🧪 Menjalankan Pengujian (Unit Tests)

Seluruh komponen memiliki pengujian unit otomatis yang mencakup autentikasi, handler, repositori, database, dan embedding aset:

```powershell
go test -v ./...
```

---

## 📄 Lisensi & Kredit

- Dikembangkan oleh **Flowspace Team**
- Ditenagai oleh **Go**, **GORM**, **Neon PostgreSQL**, dan **Cloudflare**
