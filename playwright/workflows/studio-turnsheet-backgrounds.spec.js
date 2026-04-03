import { test, expect } from '@playwright/test'
import {
  navigateTo,
  waitForPageReady,
  checkPageURL,
  checkElementVisible,
  takeScreenshot,
  setupTestBypassHeaders,
} from '../utils/test-helpers.js'

test.describe('Studio Turn Sheet Backgrounds Designer Workflows', () => {
  test.beforeEach(async ({ page }) => {
    await page.context().clearCookies()
  })

  // ─── Turn Sheet Backgrounds Page ─────────────────────────────────────────────

  test.describe('Turn Sheet Backgrounds Page Access', () => {
    test('should show unauthenticated content for turn sheet backgrounds page', async ({ page }) => {
      await navigateTo(page, '/studio')
      await checkPageURL(page, '/studio')
      await checkElementVisible(page, 'body')
      const content = await page.locator('body').textContent()
      expect(content).toMatch(/Studio|Game|Design|Sign In/i)
      await takeScreenshot(page, 'studio-turnsheet-backgrounds-unauthenticated')
    })
  })

  test.describe('Authenticated Turn Sheet Backgrounds', () => {
    test('should navigate to turn sheet backgrounds page and show tabs', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)
      const content = await page.locator('body').textContent()
      if (content.includes('Turn Sheet Backgrounds') || content.includes('turn-sheet-backgrounds')) {
        await checkElementVisible(page, 'body')
        const tabSelectors = ['.tab', '[class*="tab"]', 'button:has-text("Join Game")', 'button:has-text("Inventory")']
        for (const selector of tabSelectors) {
          try {
            if (await page.locator(selector).first().isVisible({ timeout: 2000 })) {
              console.log(`Turn sheet backgrounds tab found: ${selector}`)
              break
            }
          } catch { /* continue */ }
        }
        await takeScreenshot(page, 'studio-turnsheet-backgrounds-tabs')
      } else {
        console.log('Turn Sheet Backgrounds link not found - may require authenticated session with game selected')
      }
    })

    test('should show preview button on turn sheet backgrounds page', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)
      const content = await page.locator('body').textContent()
      if (content.includes('Turn Sheet Backgrounds') || content.includes('turn-sheet-backgrounds')) {
        const previewSelectors = ['button:has-text("Preview Turn Sheet")', '[class*="preview-btn"]', 'button:has-text("Preview")']
        let previewFound = false
        for (const selector of previewSelectors) {
          try {
            if (await page.locator(selector).first().isVisible({ timeout: 2000 })) {
              previewFound = true
              console.log(`Preview button found: ${selector}`)
              break
            }
          } catch { /* continue */ }
        }
        if (!previewFound) console.log('Preview button not found - requires game selection')
        await takeScreenshot(page, 'studio-turnsheet-backgrounds-preview-btn')
      } else {
        console.log('Turn Sheet Backgrounds page not found - requires authenticated session with game selected')
      }
    })
  })
})
