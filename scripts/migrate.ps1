# SalonFlow Track - Migration Script
# Usage: .\scripts\migrate.ps1 [up|status]

param(
    [string]$Command = "up"
)

Write-Host "=== SalonFlow Track - Database Migrations ===" -ForegroundColor Cyan

switch ($Command) {
    "up" {
        Write-Host "`nRunning migrations..." -ForegroundColor Yellow
        Push-Location backend
        go run ./cmd/server/
        Pop-Location
        # Note: Migrations run automatically on server start
        Write-Host "Migrations are applied automatically when the server starts." -ForegroundColor Green
    }
    "status" {
        Write-Host "`nMigration files:" -ForegroundColor Yellow
        Get-ChildItem -Path "database/migrations" -Filter "*.up.sql" | 
            Sort-Object Name | 
            ForEach-Object { Write-Host "  $($_.Name)" -ForegroundColor White }
    }
    default {
        Write-Host "Usage: .\scripts\migrate.ps1 [up|status]" -ForegroundColor Yellow
    }
}
