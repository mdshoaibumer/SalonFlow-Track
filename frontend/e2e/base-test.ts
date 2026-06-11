/**
 * Custom Playwright test with Wails mock injection.
 * 
 * Since Wails IPC (window.go) only works in WebView2, not regular Chromium,
 * we inject mock bindings before each page navigation for testing.
 */
import { test as base } from '@playwright/test'
import { injectWailsMock } from './fixtures/wails-mock'

export const test = base.extend({
  page: async ({ page }, use) => {
    await injectWailsMock(page)
    await use(page)
  },
})

export { expect } from '@playwright/test'
