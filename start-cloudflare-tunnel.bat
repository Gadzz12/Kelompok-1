@echo off
title Cloudflare Tunnel - Todo List App
echo ==========================================================
echo    CLOUDFLARE TUNNEL RUNNER - TO-DO LIST WEB APPLICATION
echo ==========================================================
echo.

powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0start-cloudflare-tunnel.ps1"

pause
