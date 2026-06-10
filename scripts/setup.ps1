# SalonFlow Track - Development Setup Script (Windows)
# Run this script from the project root directory

Write-Host "=== SalonFlow Track - Development Setup ===" -ForegroundColor Cyan

# Check prerequisites
Write-Host "`nChecking prerequisites..." -ForegroundColor Yellow

# Check Go
$goVersion = go version 2>$null
if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: Go is not installed. Install Go 1.24+ from https://go.dev/dl/" -ForegroundColor Red
    exit 1
}
Write-Host "  Go: $goVersion" -ForegroundColor Green

# Check Node.js
$nodeVersion = node --version 2>$null
if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: Node.js is not installed. Install Node.js 20+ from https://nodejs.org/" -ForegroundColor Red
    exit 1
}
Write-Host "  Node.js: $nodeVersion" -ForegroundColor Green

# Check npm
$npmVersion = npm --version 2>$null
Write-Host "  npm: $npmVersion" -ForegroundColor Green

# Create data directory
Write-Host "`nCreating directories..." -ForegroundColor Yellow
New-Item -ItemType Directory -Force -Path "data" | Out-Null
New-Item -ItemType Directory -Force -Path "logs" | Out-Null
Write-Host "  Created data/ and logs/ directories" -ForegroundColor Green

# Install backend dependencies
Write-Host "`nInstalling backend dependencies..." -ForegroundColor Yellow
Push-Location backend
go mod tidy
if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: Failed to install Go dependencies" -ForegroundColor Red
    Pop-Location
    exit 1
}
Pop-Location
Write-Host "  Backend dependencies installed" -ForegroundColor Green

# Install frontend dependencies
Write-Host "`nInstalling frontend dependencies..." -ForegroundColor Yellow
Push-Location frontend
npm install
if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: Failed to install frontend dependencies" -ForegroundColor Red
    Pop-Location
    exit 1
}
Pop-Location
Write-Host "  Frontend dependencies installed" -ForegroundColor Green

Write-Host "`n=== Setup Complete ===" -ForegroundColor Cyan
Write-Host "`nTo start development:" -ForegroundColor Yellow
Write-Host "  cd backend && wails dev" -ForegroundColor White
Write-Host "`nTo build the production exe:" -ForegroundColor Yellow
Write-Host "  .\scripts\build-desktop.ps1" -ForegroundColor White
Write-Host ""
