# SalonFlow Track - Desktop Build Script
# Uses `wails build` to compile frontend + Go into a single .exe

param(
    [switch]$Dev,
    [switch]$Clean,
    [switch]$Production
)

$ErrorActionPreference = "Stop"
$RootDir = Split-Path -Parent $PSScriptRoot
$BackendDir = Join-Path $RootDir "backend"
$BuildDir = Join-Path $RootDir "build"

Write-Host "=== SalonFlow Track Desktop Build ===" -ForegroundColor Cyan

# Clean build artifacts
if ($Clean) {
    Write-Host "Cleaning build artifacts..." -ForegroundColor Yellow
    if (Test-Path (Join-Path $BuildDir "bin")) { Remove-Item -Recurse -Force (Join-Path $BuildDir "bin") }
    Write-Host "Clean complete." -ForegroundColor Green
    exit 0
}

# Consolidate all migrations into backend/database/migrations/
Write-Host "Consolidating migrations..." -ForegroundColor Yellow
$RootMigrations = Join-Path $RootDir "database" "migrations"
$BackendMigrations = Join-Path $BackendDir "database" "migrations"
if (Test-Path $RootMigrations) {
    Copy-Item (Join-Path $RootMigrations "*") $BackendMigrations -Force
}
Write-Host "Migrations consolidated: $((Get-ChildItem $BackendMigrations -Filter '*.sql').Count) SQL files" -ForegroundColor Green

# Setup environment
$env:CGO_ENABLED = "1"
$env:GOARCH = "amd64"
if (Test-Path "C:\mingw64\bin\gcc.exe") {
    $env:PATH = "C:\mingw64\bin;$env:PATH"
    $env:CC = "C:\mingw64\bin\gcc.exe"
}

# Build with Wails CLI
Push-Location $BackendDir
try {
    if ($Dev) {
        Write-Host "`nBuilding in debug mode..." -ForegroundColor Cyan
        wails build -debug
    } else {
        Write-Host "`nBuilding production binary..." -ForegroundColor Cyan
        wails build
    }
    if ($LASTEXITCODE -ne 0) { throw "Wails build failed" }
} finally {
    Pop-Location
}

$OutputFile = Join-Path $BuildDir "bin" "SalonFlow-Track.exe"
Write-Host "`n=== Build Successful ===" -ForegroundColor Green
Write-Host "Output: $OutputFile" -ForegroundColor Cyan
Write-Host "Size: $([math]::Round((Get-Item $OutputFile).Length / 1MB, 2)) MB" -ForegroundColor Cyan
