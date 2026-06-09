# SalonFlow Track - Desktop E2E Test Script
# Launches the actual desktop .exe, then runs Playwright E2E tests against it.
# This validates the real desktop application stack: WebView2 + Go backend + SQLite + React frontend.

param(
    [switch]$Headed,
    [string]$TestFilter = "",
    [switch]$SkipBuild
)

$ErrorActionPreference = "Stop"
$RootDir = Split-Path -Parent $PSScriptRoot
$BackendDir = Join-Path $RootDir "backend"
$FrontendDir = Join-Path $RootDir "frontend"
$BuildDir = Join-Path $RootDir "build" "bin"
$ExePath = Join-Path $BuildDir "SalonFlow-Track.exe"

Write-Host "=== SalonFlow Track - Desktop E2E Testing ===" -ForegroundColor Cyan
Write-Host "    Tests the REAL desktop application binary" -ForegroundColor DarkCyan
Write-Host ""

# Step 0: Kill any existing instances
Get-Process | Where-Object { $_.Path -like "*SalonFlow-Track*" } | Stop-Process -Force -ErrorAction SilentlyContinue
Start-Sleep 1

# Step 1: Build if needed
if (-not $SkipBuild -and -not (Test-Path $ExePath)) {
    Write-Host "[1/4] Building desktop app via 'wails build'..." -ForegroundColor Yellow
    Push-Location $BackendDir
    $env:CGO_ENABLED = "1"
    $env:GOARCH = "amd64"
    if (Test-Path "C:\mingw64\bin\gcc.exe") {
        $env:PATH = "C:\mingw64\bin;$env:PATH"
        $env:CC = "C:\mingw64\bin\gcc.exe"
    }
    wails build
    if ($LASTEXITCODE -ne 0) { throw "Wails build failed" }
    Pop-Location
    Write-Host "    Build complete." -ForegroundColor Green
} else {
    Write-Host "[1/4] Desktop app binary ready." -ForegroundColor Green
}

# Step 2: Launch the desktop application
Write-Host "[2/4] Launching SalonFlow-Track.exe..." -ForegroundColor Yellow

$appProcess = Start-Process -FilePath $ExePath -WorkingDirectory $BackendDir -PassThru

# Wait for the API server to be ready
$maxRetries = 20
$retryCount = 0
$apiReady = $false
while ($retryCount -lt $maxRetries) {
    Start-Sleep 1
    try {
        $response = Invoke-WebRequest "http://localhost:8080/api/v1/health" -UseBasicParsing -TimeoutSec 2 -ErrorAction Stop
        if ($response.StatusCode -eq 200) {
            $apiReady = $true
            break
        }
    } catch {
        $retryCount++
    }
}

if (-not $apiReady) {
    Write-Host "    ERROR: Desktop app API did not start within 20 seconds." -ForegroundColor Red
    if ($appProcess -and -not $appProcess.HasExited) { $appProcess | Stop-Process -Force }
    exit 1
}

$healthData = (Invoke-WebRequest "http://localhost:8080/api/v1/health" -UseBasicParsing).Content | ConvertFrom-Json
Write-Host "    Desktop app running (PID: $($appProcess.Id))" -ForegroundColor Green
Write-Host "    API: healthy | DB: $($healthData.database_status) | Version: $($healthData.version)" -ForegroundColor DarkGreen

# Step 3: Run Playwright E2E tests
Write-Host ""
Write-Host "[3/4] Running Playwright E2E tests against desktop app..." -ForegroundColor Yellow
Write-Host "      (Frontend via Vite dev server, API served by desktop .exe)" -ForegroundColor DarkGray
Write-Host ""

Push-Location $FrontendDir
$playwrightArgs = @("playwright", "test", "--config=playwright.desktop.config.ts", "--reporter=line")
if ($TestFilter) { $playwrightArgs += "--grep"; $playwrightArgs += $TestFilter }
if ($Headed) { $playwrightArgs += "--headed" }

try {
    & npx @playwrightArgs
    $testExitCode = $LASTEXITCODE
} finally {
    Pop-Location
}

# Step 4: Cleanup
Write-Host ""
Write-Host "[4/4] Stopping desktop app (PID: $($appProcess.Id))..." -ForegroundColor Yellow
if ($appProcess -and -not $appProcess.HasExited) {
    $appProcess | Stop-Process -Force
}
Get-Process | Where-Object { $_.Path -like "*SalonFlow-Track*" } | Stop-Process -Force -ErrorAction SilentlyContinue

Write-Host ""
if ($testExitCode -eq 0) {
    Write-Host "=== ALL DESKTOP E2E TESTS PASSED ===" -ForegroundColor Green
    Write-Host "    Tested against: SalonFlow-Track.exe (real desktop binary)" -ForegroundColor DarkGreen
} else {
    Write-Host "=== SOME DESKTOP E2E TESTS FAILED (exit code: $testExitCode) ===" -ForegroundColor Red
}

exit $testExitCode
