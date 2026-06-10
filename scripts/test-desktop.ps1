# SalonFlow Track - Desktop E2E Test Script
# Uses 'wails dev' to run the REAL Go backend (with SQLite) and frontend,
# then runs Playwright E2E tests against the live application.
# This tests the exact same code as the production .exe - same Go handlers,
# same SQLite database, same React frontend - just served via HTTP for testability.

param(
    [switch]$Headed,
    [string]$TestFilter = ""
)

$ErrorActionPreference = "Stop"
$RootDir = Split-Path -Parent $PSScriptRoot
$BackendDir = Join-Path $RootDir "backend"
$FrontendDir = Join-Path $RootDir "frontend"

Write-Host "=== SalonFlow Track - Desktop E2E Testing ===" -ForegroundColor Cyan
Write-Host "    Tests the REAL Go backend + SQLite + React frontend via wails dev" -ForegroundColor DarkCyan
Write-Host ""

# Step 0: Kill any existing wails dev processes
Get-Process | Where-Object { $_.Path -like "*SalonFlow-Track*" } | Stop-Process -Force -ErrorAction SilentlyContinue
Get-Process | Where-Object { $_.Name -eq "msedgewebview2" } | Stop-Process -Force -ErrorAction SilentlyContinue
# Kill any node/vite processes on ports we'll use
$portProcs = netstat -ano | Select-String ":34115\s|:5173\s" | ForEach-Object { ($_ -split '\s+')[-1] } | Sort-Object -Unique
foreach ($pid in $portProcs) { if ($pid -match '^\d+$' -and $pid -ne '0') { Stop-Process -Id $pid -Force -ErrorAction SilentlyContinue } }
Start-Sleep 2

# Use a clean test data directory
$testDataDir = Join-Path $env:TEMP "salonflow-e2e-test"
if (Test-Path $testDataDir) { Remove-Item -Recurse -Force $testDataDir }
New-Item -ItemType Directory -Path $testDataDir -Force | Out-Null

# Step 1: Start wails dev (launches Go backend + Vite frontend)
Write-Host "[1/3] Starting 'wails dev' (Go backend + Vite frontend)..." -ForegroundColor Yellow
Push-Location $BackendDir
$env:CGO_ENABLED = "1"
$env:GOARCH = "amd64"
if (Test-Path "C:\mingw64\bin\gcc.exe") {
    $env:PATH = "C:\mingw64\bin;$env:PATH"
    $env:CC = "C:\mingw64\bin\gcc.exe"
}
# Set test data path so DB is created in temp directory
$env:SALONFLOW_DATA_DIR = $testDataDir

# Start wails dev in background (it opens a window + serves frontend on Vite port)
$wailsJob = Start-Job -ScriptBlock {
    param($dir, $dataDir)
    Set-Location $dir
    $env:CGO_ENABLED = "1"
    $env:GOARCH = "amd64"
    $env:PATH = "C:\mingw64\bin;$env:PATH"
    $env:CC = "C:\mingw64\bin\gcc.exe"
    $env:SALONFLOW_DATA_DIR = $dataDir
    wails dev -browser 2>&1
} -ArgumentList $BackendDir, $testDataDir
Pop-Location

# Wait for Vite dev server to be ready
$maxRetries = 60
$retryCount = 0
$serverReady = $false
Write-Host "    Waiting for dev server..." -ForegroundColor DarkGray
while ($retryCount -lt $maxRetries) {
    Start-Sleep 2
    try {
        $response = Invoke-WebRequest "http://localhost:34115" -UseBasicParsing -TimeoutSec 3 -ErrorAction Stop
        if ($response.StatusCode -eq 200) {
            $serverReady = $true
            break
        }
    } catch {
        # Also try Vite default port
        try {
            $response = Invoke-WebRequest "http://localhost:5173" -UseBasicParsing -TimeoutSec 2 -ErrorAction Stop
            if ($response.StatusCode -eq 200) {
                $env:PLAYWRIGHT_BASE_URL = "http://localhost:5173"
                $serverReady = $true
                break
            }
        } catch {}
        $retryCount++
    }
}

if (-not $serverReady) {
    Write-Host "    ERROR: Dev server did not start within 120 seconds." -ForegroundColor Red
    Write-Host "    Job output:" -ForegroundColor DarkGray
    Receive-Job $wailsJob -ErrorAction SilentlyContinue | Select-Object -Last 20
    Stop-Job $wailsJob -ErrorAction SilentlyContinue
    Remove-Job $wailsJob -ErrorAction SilentlyContinue
    exit 1
}

if (-not $env:PLAYWRIGHT_BASE_URL) { $env:PLAYWRIGHT_BASE_URL = "http://localhost:34115" }
Write-Host "    Dev server ready at $($env:PLAYWRIGHT_BASE_URL)" -ForegroundColor Green

# Step 2: Run Playwright E2E tests
Write-Host ""
Write-Host "[2/3] Running Playwright E2E tests against live backend..." -ForegroundColor Yellow
Write-Host "      (URL: $($env:PLAYWRIGHT_BASE_URL))" -ForegroundColor DarkGray
Write-Host ""

Push-Location $FrontendDir
$playwrightArgs = @("playwright", "test", "--config=playwright.desktop.config.ts", "--reporter=list")
if ($TestFilter) { $playwrightArgs += "--grep"; $playwrightArgs += $TestFilter }
if ($Headed) { $playwrightArgs += "--headed" }

try {
    & npx @playwrightArgs
    $testExitCode = $LASTEXITCODE
} finally {
    Pop-Location
}

# Step 3: Cleanup
Write-Host ""
Write-Host "[3/3] Stopping dev server..." -ForegroundColor Yellow
Stop-Job $wailsJob -ErrorAction SilentlyContinue
Remove-Job $wailsJob -ErrorAction SilentlyContinue
Get-Process | Where-Object { $_.Path -like "*SalonFlow-Track*" } | Stop-Process -Force -ErrorAction SilentlyContinue
Get-Process | Where-Object { $_.Name -eq "msedgewebview2" } | Stop-Process -Force -ErrorAction SilentlyContinue

Write-Host ""
if ($testExitCode -eq 0) {
    Write-Host "=== ALL DESKTOP E2E TESTS PASSED ===" -ForegroundColor Green
    Write-Host "    Tested: Real Go backend + SQLite + React frontend" -ForegroundColor DarkGreen
} else {
    Write-Host "=== SOME DESKTOP E2E TESTS FAILED (exit code: $testExitCode) ===" -ForegroundColor Red
    Write-Host "    Run with -Headed to see the browser during tests" -ForegroundColor Yellow
}

exit $testExitCode
