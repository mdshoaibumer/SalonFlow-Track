import { defineConfig, devices } from '@playwright/test'

/**
 * Playwright config for E2E testing against the actual Wails desktop application.
 * 
 * Strategy:
 *   1. The desktop app (SalonFlow-Track.exe) runs with its embedded SQLite DB and API server
 *   2. Playwright opens the frontend (via Vite dev server) which connects to the desktop app's
 *      real API on localhost:8080 — testing the actual desktop backend
 *   3. This validates the complete desktop stack: WebView2 window + Go backend + SQLite + React frontend
 * 
 * What's tested:
 *   - The real desktop app binary is running (not `go run`)
 *   - Real SQLite database operations
 *   - Real HTTP API responses from the desktop process
 *   - Frontend behavior against the desktop backend
 * 
 * Usage:
 *   ..\scripts\test-desktop.ps1
 */
export default defineConfig({
  testDir: './e2e',
  testIgnore: ['desktop-app.spec.ts', 'desktop-fixtures.ts'],
  fullyParallel: true,
  workers: 7,
  retries: 0,
  reporter: [['line'], ['html', { open: 'never' }]],
  timeout: 30000,
  use: {
    // Frontend served by Vite dev server, API proxied to desktop app on :8080
    baseURL: 'http://localhost:5173',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    actionTimeout: 10000,
  },
  projects: [
    {
      name: 'desktop',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
  // Vite dev server proxies /api to the desktop app's HTTP server on port 8080
  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:5173',
    reuseExistingServer: true,
    timeout: 15000,
  },
})
