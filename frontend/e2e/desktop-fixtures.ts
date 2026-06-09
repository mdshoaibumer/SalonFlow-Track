import { test as base, expect, type BrowserContext, type Page } from '@playwright/test'

/**
 * Custom Playwright fixtures for desktop E2E testing.
 * Connects to the running Wails desktop app's WebView2 via CDP.
 */

export const test = base.extend<{ desktopPage: Page }>({
  desktopPage: async ({ playwright }, use) => {
    // Connect to WebView2 via Chrome DevTools Protocol
    const browser = await playwright.chromium.connectOverCDP('http://localhost:9222')
    const contexts = browser.contexts()
    
    let context: BrowserContext
    let page: Page

    if (contexts.length > 0) {
      // Use existing context (the WebView2 window)
      context = contexts[0]
      const pages = context.pages()
      page = pages.length > 0 ? pages[0] : await context.newPage()
    } else {
      // Create new context if none exists
      context = await browser.newContext()
      page = await context.newPage()
    }

    await use(page)
    
    // Don't close the browser - it's the desktop app
    await browser.close()
  },
})

export { expect }
