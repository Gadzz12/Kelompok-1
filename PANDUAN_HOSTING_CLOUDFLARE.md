# Panduan Hosting Aplikasi To-Do List melalui Cloudflare

Panduan ini menjelaskan cara mempublikasikan aplikasi To-Do List Anda ke internet melalui **Cloudflare**, sehingga aplikasi dapat diakses melalui URL publik HTTPS tanpa perlu lagi menggunakan `localhost`.

---

## 1. Pemahaman Arsitektur: Go Backend & Cloudflare

Banyak pengguna pemula mengira program Go bisa langsung diunggah ke menu **Cloudflare Pages**. Namun secara teknis:
- **Cloudflare Pages / Workers** dirancang khusus untuk file statis (HTML/CSS/JS) atau serverless Edge JavaScript/TypeScript (V8 isolate).
- Aplikasi ini adalah **aplikasi web backend Go lengkap** (menggunakan HTTP router, GORM, HTML template rendering, dan koneksi langsung ke database PostgreSQL Neon).

Untuk menghubungkannya dengan Cloudflare, terdapat **2 metode resmi & terbaik**:

| Metode | Deskripsi | Kelebihan |
| :--- | :--- | :--- |
| **Metode A: Cloudflare Tunnel (Quick Tunnel)** | Menghubungkan aplikasi dari komputer Anda ke jaringan Cloudflare secara instan. | 100% Gratis, instan (1 menit), tanpa perlu beli domain, tanpa akun. |
| **Metode B: Cloudflare Zero Trust (Custom Domain)** | Menghubungkan tunnel ke nama domain pribadi Anda (misal `todo.domainanda.com`). | Alamat URL permanen dan profesional, gratis dari Cloudflare. |
| **Metode C: Cloud 24/7 (Render / Koyeb) + Cloudflare** | Menjalankan server Go di server cloud gratis, lalu diarahkan ke domain Cloudflare. | Aplikasi online 24 jam nonstop tanpa perlu laptop menyala. |

---

## 2. Metode A: Cloudflare Quick Tunnel (Paling Cepat & Instan)

Metode ini membuat tunnel publik aman berprotokol HTTPS langsung ke server lokal Anda.

### Langkah-langkah:
1. Pastikan server Go berjalan:
   ```powershell
   go run main.go
   ```
2. Buka terminal baru atau klik ganda (double-click) file:
   ```text
   start-cloudflare-tunnel.bat
   ```
   *(atau jalankan melalui PowerShell: `.\start-cloudflare-tunnel.ps1`)*

3. Skrip akan otomatis mengunduh `cloudflared.exe` resmi jika belum ada, kemudian menjalankan:
   ```powershell
   cloudflared tunnel --url http://localhost:8080
   ```

4. Pada jendela terminal akan muncul kotak informasi dari Cloudflare berisi link HTTPS:
   ```text
   +--------------------------------------------------------------------------------------------+
   |  Your quick Tunnel has been created! Visit it at (it may take some time to be reachable):  |
   |  https://bebas-nama-acak.trycloudflare.com                                                 |
   +--------------------------------------------------------------------------------------------+
   ```

5. **Selesai!** Buka link `https://bebas-nama-acak.trycloudflare.com` di HP, browser, atau bagikan ke teman. Anda tidak perlu lagi mengakses via `localhost:8080`.

---

## 3. Metode B: Cloudflare Tunnel dengan Domain Sendiri (Zero Trust)

Jika Anda sudah memiliki nama domain yang aktif di akun Cloudflare (misal: `domainanda.com`):

1. Buka [Cloudflare Zero Trust Dashboard](https://one.dash.cloudflare.com/).
2. Masuk ke menu **Networks** > **Tunnels** > klik **Add a tunnel**.
3. Pilih **Cloudflared**, beri nama tunnel (contoh: `todo-app`), klik **Save tunnel**.
4. Salin token koneksi yang diberikan oleh Cloudflare.
5. Jalankan tunnel dengan token tersebut di terminal:
   ```powershell
   .\cloudflared.exe service install <TOKEN_DARI_CLOUDFLARE>
   ```
6. Pada tab **Public Hostname** di dashboard Cloudflare:
   - **Subdomain**: `todo` (atau biarkan kosong jika root)
   - **Domain**: pilih domain Anda (`domainanda.com`)
   - **Type**: `HTTP`
   - **URL**: `localhost:8080`
7. Klik **Save hostname**. Aplikasi Anda sekarang dapat diakses secara permanen di `https://todo.domainanda.com`.

---

## 4. Metode C: Deploy Cloud 24/7 Gratis (Render.com)

Jika Anda ingin aplikasi tetap aktif 24 jam nonstop tanpa harus menyalakan komputer:

1. Push kode proyek ini ke repositori **GitHub**.
2. Buka [Render.com](https://render.com) dan buat akun (gratis).
3. Klik **New** > **Web Service** > hubungkan ke repositori GitHub Anda.
4. Pengaturan:
   - **Runtime**: `Go`
   - **Build Command**: `go build -o server .`
   - **Start Command**: `./server`
5. Tambahkan **Environment Variables** sesuai file `.env` Anda (`DB_HOST`, `DB_USER`, `DB_PASSWORD`, dll).
6. Render akan memberikan URL gratis `https://todo-app.onrender.com`.
7. (Opsional) Di dashboard Cloudflare, tambahkan **CNAME Record** yang mengarahkan subdomain Anda ke `todo-app.onrender.com` dengan proxy Cloudflare aktif (ikon awan oranye) untuk proteksi DDoS dan caching.
