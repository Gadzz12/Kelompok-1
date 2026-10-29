# 🚀 Panduan Lengkap Hosting Aplikasi Flowspace ke Vercel

Panduan ini menjelaskan langkah demi langkah cara mempublikasikan aplikasi **Flowspace (Go To-Do List & Workspace)** ke **Vercel** agar dapat diakses publik 24/7 secara gratis dengan domain HTTPS otomatis (`https://nama-proyek-anda.vercel.app`).

---

## 📋 Ringkasan Arsitektur Vercel & Go

Aplikasi Flowspace dibangun dengan bahasa **Go (Golang)**. Vercel mendukung Go secara native melalui Go Runtime:
- Semua template HTML (`templates/`) dan file statis CSS/JS/Gambar (`static/`) sudah di-embed langsung ke dalam program menggunakan fitur bawaan Go `//go:embed`.
- File konfigurasi [`vercel.json`](file:///D:/project%20copy%20team%201/todo-list/vercel.json) telah disediakan untuk mengarahkan seluruh rute trafik web ke entry point `main.go`.
- Aplikasi berjalan otomatis di infrastruktur Vercel Serverless tanpa perlu mengelola server manual (serverless).

---

## 🛠️ Persiapan Sebelum Deploy

Sebelum deploy ke Vercel, pastikan Anda memiliki:
1. **Akun Vercel**: Daftar gratis di [vercel.com](https://vercel.com/signup) (bisa login menggunakan akun GitHub).
2. **Database PostgreSQL Cloud**: Database cloud gratis yang aktif (seperti [Neon.tech](https://neon.tech) atau [Supabase](https://supabase.com)).
   - Pastikan Anda sudah memiliki connection string `DATABASE_URL` dengan format:
     ```text
     postgresql://USER:PASSWORD@HOST:5432/DATABASE?sslmode=require
     ```
3. **Akun GitHub** (untuk Metode 1 - Otomatis).

---

## 🌟 Metode 1: Deploy Otomatis via GitHub (Sangat Direkomendasikan)

Metode ini adalah cara paling mudah, praktis, dan standar industri. Setiap kali Anda melakukan `git push` ke GitHub, Vercel akan otomatis melakukan build dan deploy versi terbaru.

### Langkah-langkah:

### 1. Push Kode ke Repository GitHub
Pastikan seluruh file proyek sudah di-commit dan di-push ke repository GitHub Anda (bisa public atau private):
```powershell
git add .
git commit -m "Setup hosting Vercel dan konfigurasi vercel.json"
git push origin main
```

### 2. Hubungkan Repository ke Vercel
1. Buka dashboard Vercel di [https://vercel.com](https://vercel.com).
2. Klik tombol **"Add New..."** di kanan atas, lalu pilih **"Project"**.
3. Pada bagian **"Import Git Repository"**, cari dan pilih repository GitHub aplikasi Anda, lalu klik **"Import"**.

### 3. Konfigurasi Pengaturan Project
Pada halaman konfigurasi import:
- **Project Name**: Masukkan nama proyek yang Anda inginkan (misal `flowspace-todo`). Nama ini akan menjadi subdomain Vercel Anda: `https://flowspace-todo.vercel.app`.
- **Framework Preset**: Pilih **`Go`** (atau biarkan Vercel mendeteksi otomatis).
- **Root Directory**: `./` (biarkan default).
- **Build and Output Settings**: Biarkan default (Vercel otomatis membaca `go.mod` dan `vercel.json`).

### 4. Tambahkan Environment Variable (`DATABASE_URL`)
> ⚠️ **Sangat Penting**: Aplikasi memerlukan database PostgreSQL cloud agar dapat berjalan.
1. Buka bagian **"Environment Variables"**.
2. Masukkan:
   - **Key**: `DATABASE_URL`
   - **Value**: Masukkan connection string database Neon/Supabase Anda (contoh: `postgresql://flowspace_owner:xxxx@ep-cool-sample.ap-southeast-1.aws.neon.tech/flowspace?sslmode=require`).
3. Klik tombol **"Add"**.

### 5. Klik Deploy!
1. Klik tombol **"Deploy"**.
2. Vercel akan mengunduh Go module, mengompilasi aplikasi, dan melakukan deployment.
3. Setelah proses selesai (biasanya ~1 menit), layar perayaan akan muncul dengan tombol **"Continue to Dashboard"** atau **"Visit"**.
4. Klik tautan URL yang diberikan (contoh: `https://flowspace-todo.vercel.app`) untuk membuka aplikasi Anda secara publik!

---

## ⚡ Metode 2: Deploy Cepat via Vercel CLI (Terminal)

Jika Anda tidak ingin menghubungkan ke GitHub dan ingin deploy langsung dari terminal komputer:

### 1. Pasang Vercel CLI
Buka PowerShell atau Command Prompt, lalu instal Vercel CLI:
```powershell
npm install -g vercel
```

### 2. Login ke Vercel
```powershell
vercel login
```
Ikuti instruksi di terminal untuk verifikasi email atau login browser.

### 3. Lakukan Deploy Pertama (Preview)
Di direktori proyek `todo-list`, jalankan:
```powershell
vercel
```
Jawab pertanyaan interaktif dari Vercel CLI:
- `Set up and deploy ...?` -> ketik `y`
- `Which scope do you want to deploy to?` -> pilih akun Anda
- `Link to existing project?` -> ketik `n` (jika proyek baru)
- `What's your project's name?` -> tekan Enter (atau beri nama, misal `flowspace`)
- `In which directory is your code located?` -> tekan Enter (`./`)
- `Want to modify these settings?` -> ketik `n`

### 4. Masukkan Environment Variable DATABASE_URL
Tambahkan koneksi database Anda ke Vercel:
```powershell
vercel env add DATABASE_URL production
```
Lalu tempel (paste) connection string Neon/PostgreSQL Anda ketika diminta.

### 5. Deploy ke Production
Jalankan perintah berikut untuk meluncurkan ke URL produksi resmi:
```powershell
vercel --prod
```
Terminal akan langsung menampilkan URL publik aktif aplikasi Anda!

---

## 🌐 Menghubungkan Custom Domain (Domain Sendiri)

Jika Anda memiliki domain sendiri (misal `todo.domainanda.com`):

1. Masuk ke dashboard proyek di Vercel.
2. Buka menu **Settings** > **Domains**.
3. Masukkan nama domain atau subdomain Anda, lalu klik **Add**.
4. Vercel akan menampilkan DNS record (CNAME atau A record).
5. Tambahkan record DNS tersebut di penyedia domain Anda (seperti Namecheap, Niagahoster, Domainesia, dsb).
6. Tunggu beberapa menit hingga status DNS terverifikasi (Centang Hijau). SSL/HTTPS otomatis terbit gratis dari Vercel!

---

## ⚙️ Penjelasan File `vercel.json`

File [`vercel.json`](file:///D:/project%20copy%20team%201/todo-list/vercel.json) di root proyek telah dikonfigurasi secara optimal:

```json
{
  "$schema": "https://openapi.vercel.sh/vercel.json",
  "framework": "go",
  "build": {
    "env": {
      "GO_BUILD_FLAGS": "-ldflags '-s -w'"
    }
  }
}
```

- **`framework: "go"`**: Memberi tahu Vercel untuk menggunakan runtime resmi Go yang secara otomatis merutekan semua request HTTP ke server `main.go` dengan mempertahankan URL path asli.
- **`GO_BUILD_FLAGS: "-ldflags '-s -w'"`**: Memangkas simbol debug dari binary hasil kompilasi sehingga ukuran file lebih ringkas dan waktu cold-start jauh lebih kencang.
- **Catatan Routing**: Blok legacy `rewrites` ke `/main.go` dihilangkan agar sesuai dengan aturan routing Vercel terbaru dan tidak menimpa request path internal aplikasi.

---

## 🔍 Troubleshooting & Pertanyaan Umum (FAQ)

### 1. Muncul error `DATABASE_URL belum diatur` di Vercel Logs
- **Penyebab**: Variabel lingkungan `DATABASE_URL` belum ditambahkan di Vercel.
- **Solusi**: Buka dashboard Vercel > Proyek Anda > **Settings** > **Environment Variables**. Tambahkan key `DATABASE_URL` dengan value connection string database Neon Anda, lalu lakukan **Redeploy** (Deployments > klik titik tiga pada deployment terbaru > Redeploy).

### 2. Error koneksi database: `no pg_hba.conf entry for host` atau TLS Error
- **Penyebab**: Connection string PostgreSQL tidak menyertakan parameter SSL.
- **Solusi**: Pastikan di ujung `DATABASE_URL` terdapat `?sslmode=require`. Contoh:
  `postgresql://user:pass@host:5432/db?sslmode=require`

### 3. Apakah data saya aman jika Vercel serverless mati/reboot?
- **Ya, 100% Aman**: Semua data tugas, akun pengguna, profil, dan kolaborasi disimpan di database cloud eksternal (PostgreSQL Neon). Vercel hanya menjalankan kode backend tanpa menyimpan data di disk lokal.

---

🎉 **Selamat! Aplikasi Flowspace Anda kini telah online di Vercel.**
