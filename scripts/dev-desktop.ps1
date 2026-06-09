# SalonFlow Track - Desktop Development Script
# Starts the Wails desktop app in development mode (hot-reload for frontend)

param(
    [switch]$NoBrowser
)

$ErrorActionPreference = "Stop"
$RootDir = Split-Path -Parent $PSScriptRoot
$BackendDir = Join-Path $RootDir "backend"
$FrontendDir = Join-Path $RootDir "frontend"

Write-Host "=== SalonFlow Track - Desktop Dev Mode ===" -ForegroundColor Cyan
Write-Host ""

# Ensure 64-bit GCC is available
if (Test-Path "C:\mingw64\bin\gcc.exe") {
    $env:PATH = "C:\mingw64\bin;$env:PATH"
    $env:CC = "C:\mingw64\bin\gcc.exe"
} elseif (Test-Path "C:\msys64\mingw64\bin\gcc.exe") {
    $env:PATH = "C:\msys64\mingw64\bin;$env:PATH"
    $env:CC = "C:\msys64\mingw64\bin\gcc.exe"
}

$env:CGO_ENABLED = "1"

# Start the backend (HTTP API server)
Write-Host "Starting Go backend server..." -ForegroundColor Yellow
$backendJob = Start-Job -ScriptBlock {
    Set-Location $using:BackendDir
    $env:CGO_ENABLED = "1"
    if (Test-Path "C:\mingw64\bin\gcc.exe") {
        $env:PATH = "C:\mingw64\bin;$env:PATH"
        $env:CC = "C:\mingw64\bin\gcc.exe"
    }
    go run ./cmd/desktop/ 2>&1
}

Start-Sleep -Seconds 3

# Start the frontend dev server (Vite with hot-reload)
Write-Host "Starting frontend dev server (Vite HMR)..." -ForegroundColor Yellow
$frontendJob = Start-Job -ScriptBlock {
    Set-Location $using:FrontendDir
    npm run dev 2>&1
}

Start-Sleep -Seconds 2

Write-Host ""
Write-Host "Development servers running:" -ForegroundColor Green
Write-Host "  Backend API:  http://localhost:8080/api/v1" -ForegroundColor White
Write-Host "  Frontend:     http://localhost:5173" -ForegroundColor White
Write-Host ""
Write-Host "Open http://localhost:5173 in your browser for development." -ForegroundColor Cyan
Write-Host "Press Ctrl+C to stop all servers." -ForegroundColor Yellow
Write-Host ""

try {
    while ($true) {
        # Check if jobs are still running
        $bState = (Get-Job -Id $backendJob.Id).State
        $fState = (Get-Job -Id $frontendJob.Id).State

        if ($bState -eq "Failed") {
            Write-Host "Backend failed:" -ForegroundColor Red
            Receive-Job $backendJob
            break
        }
        if ($fState -eq "Failed") {
            Write-Host "Frontend failed:" -ForegroundColor Red
            Receive-Job $frontendJob
            break
        }

        Start-Sleep -Seconds 1
    }
} finally {
    Write-Host "`nStopping servers..." -ForegroundColor Yellow
    Stop-Job $backendJob -ErrorAction SilentlyContinue
    Stop-Job $frontendJob -ErrorAction SilentlyContinue
    Remove-Job $backendJob -ErrorAction SilentlyContinue
    Remove-Job $frontendJob -ErrorAction SilentlyContinue
    Write-Host "All servers stopped." -ForegroundColor Green
}
