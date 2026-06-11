import { defineConfig } from '@playwright/test'

/**
 * Playwright config for E2E testing against the real Wails backend + SQLite.
 * 
 * Strategy:
 *   1. 'wails dev' starts the real Go backend (with SQLite) and Vite frontend dev server
 *   2. Playwright connects to the dev server URL (same as what WebView2 renders)
 *   3. Tests interact with the REAL backend — same Go handlers, same database
 *
 * What's tested:
 *   - Real Go backend processing (same code as production .exe)
 *   - Real SQLite database operations via Wails IPC bindings
 *   - Real React frontend rendering (same components)
 *   - Full end-to-end user flows (navigation, CRUD, billing, etc.)
 *
 * Usage:
 *   ..\scripts\test-desktop.ps1              # headless
 *   ..\scripts\test-desktop.ps1 -Headed      # watch browser during tests
 *   ..\scripts\test-desktop.ps1 -TestFilter "Staff"   # run specific tests
 */
export default defineConfig({
  testDir: './e2e',
  fullyParallel: false,
  workers: 1,
  retries: 1,
  reporter: [['list'], ['html', { open: 'never' }]],
  timeout: 30000,
  expect: {
    timeout: 15000,
  },
  use: {
    baseURL: process.env.PLAYWRIGHT_BASE_URL || 'http://localhost:5173',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    actionTimeout: 15000,
  },
  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:5173',
    reuseExistingServer: true,
  },
  projects: [
    {
      name: 'desktop-chromium',
      use: {
        browserName: 'chromium',
        viewport: { width: 1400, height: 900 },
      },
    },
  ],
})
