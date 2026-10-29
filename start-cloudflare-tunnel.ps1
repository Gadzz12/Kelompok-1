# PowerShell Script: Jalankan Cloudflare Tunnel untuk Aplikasi To-Do List
# =======================================================================

$ErrorActionPreference = "Stop"

Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host "   CLOUDFLARE TUNNEL RUNNER - TO-DO LIST WEB APPLICATION   " -ForegroundColor Cyan
Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host ""

$exePath = ".\cloudflared.exe"

# 1. Cek apakah cloudflared sudah ada
if (-not (Test-Path $exePath)) {
    # Cek apakah terinstall secara global di sistem
    $systemCloudflared = Get-Command cloudflared -ErrorAction SilentlyContinue
    if ($systemCloudflared) {
        $exePath = "cloudflared"
        Write-Host "[OK] Ditemukan cloudflared di sistem: $($systemCloudflared.Source)" -ForegroundColor Green
    } else {
        Write-Host "[INFO] cloudflared.exe belum ditemukan di komputer Anda." -ForegroundColor Yellow
        Write-Host "[INFO] Mengunduh cloudflared.exe resmi dari Cloudflare (GitHub Releases)..." -ForegroundColor Yellow
        
        $downloadUrl = "https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-windows-amd64.exe"
        try {
            [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
            Invoke-WebRequest -Uri $downloadUrl -OutFile $exePath -UseBasicParsing
            Write-Host "[OK] Berhasil mengunduh cloudflared.exe!" -ForegroundColor Green
        } catch {
            Write-Host "[GAGAL] Gagal mengunduh otomatis: $_" -ForegroundColor Red
            Write-Host "Silakan unduh manual dari: https://github.com/cloudflare/cloudflared/releases" -ForegroundColor Red
            Write-Host "Simpan file tersebut sebagai 'cloudflared.exe' di folder ini." -ForegroundColor Red
            pause
            exit 1
        }
    }
} else {
    Write-Host "[OK] cloudflared.exe ditemukan di direktori proyek." -ForegroundColor Green
}

# 2. Periksa apakah server Go (port 8080) sudah berjalan
$port = 8080
$portActive = $false
try {
    $tcp = New-Object System.Net.Sockets.TcpClient
    $async = $tcp.BeginConnect("127.0.0.1", $port, $null, $null)
    $wait = $async.AsyncWaitHandle.WaitOne(1000, $false)
    if ($wait -and $tcp.Connected) {
        $tcp.EndConnect($async)
        $portActive = $true
    }
    $tcp.Close()
} catch {
    $portActive = $false
}

if (-not $portActive) {
    Write-Host ""
    Write-Host "[PERINGATAN] Aplikasi Go To-Do belum aktif di port 8080!" -ForegroundColor Yellow
    Write-Host "Silakan jalankan aplikasi di terminal lain terlebih dahulu:" -ForegroundColor Yellow
    Write-Host "    go run main.go" -ForegroundColor White
    Write-Host ""
    $answer = Read-Host "Apakah Anda ingin saya jalankan 'go run main.go' di background sekarang? (y/n)"
    if ($answer -eq 'y' -or $answer -eq 'Y') {
        Write-Host "Menjalankan 'go run main.go'..." -ForegroundColor Cyan
        Start-Process powershell -ArgumentList "-NoExit", "-Command", "go run main.go"
        Start-Sleep -Seconds 3
    }
}

Write-Host ""
Write-Host "==========================================================" -ForegroundColor Green
Write-Host " Membuka Cloudflare Tunnel ke http://localhost:8080...    " -ForegroundColor Green
Write-Host " Tunggu beberapa detik, URL publik HTTPS akan muncul di bawah: " -ForegroundColor Green
Write-Host "==========================================================" -ForegroundColor Green
Write-Host ""

& $exePath tunnel --url http://localhost:8080
