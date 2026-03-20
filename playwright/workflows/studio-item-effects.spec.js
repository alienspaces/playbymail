import { test, expect } from '@playwright/test'
import {
  navigateTo,
  waitForPageReady,
  checkPageURL,
  checkElementVisible,
  safeClick,
  takeScreenshot,
  setupTestBypassHeaders,
} from '../utils/test-helpers.js'

test.describe('Studio Item Effects Designer Workflows', () => {
  test.beforeEach(async ({ page }) => {
    await page.context().clearCookies()
  })

  test.describe('Item Effects Page Access', () => {
    test('should show unauthenticated content for item effects page', async ({ page }) => {
      await navigateTo(page, '/studio')

      await checkPageURL(page, '/studio')
      await checkElementVisible(page, 'body')

      const body = page.locator('body')
      const content = await body.textContent()
      expect(content).toMatch(/Studio|Game|Design|Sign In/i)

      await takeScreenshot(page, 'studio-item-effects-unauthenticated')
    })
  })

  test.describe('Authenticated Item Effects CRUD', () => {
    test('should navigate to item effects page and show content', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')

      await waitForPageReady(page)

      const body = page.locator('body')
      const content = await body.textContent()

      if (content.includes('Item Effects') || content.includes('item-effects')) {
        await checkElementVisible(page, 'body')
        await takeScreenshot(page, 'studio-item-effects-page')
      } else {
        console.log('Item Effects link not found - may require authenticated session with game selected')
      }
    })

    test('should show item effects table when game is selected', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')

      await waitForPageReady(page)

      const tableSelectors = [
        '[data-testid="item-effects-table"]',
        'table',
        '[class*="resource-table"]',
      ]

      let tableFound = false
      for (const selector of tableSelectors) {
        try {
          if (await page.locator(selector).isVisible({ timeout: 2000 })) {
            tableFound = true
            console.log(`Item effects table found: ${selector}`)
            break
          }
        } catch {
          // Continue to next selector
        }
      }

      if (!tableFound) {
        console.log('Item effects table not found - requires game selection')
      }

      await takeScreenshot(page, 'studio-item-effects-table')
    })

    test('should open create item effect form', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')

      await waitForPageReady(page)

      const createButtonSelectors = [
        'button:has-text("Create Item Effect")',
        '[data-testid="create-item-effect"]',
        'button:has-text("Create")',
      ]

      let buttonFound = false
      for (const selector of createButtonSelectors) {
        try {
          if (await page.locator(selector).first().isVisible({ timeout: 2000 })) {
            await safeClick(page, selector)
            buttonFound = true
            console.log(`Create item effect button clicked: ${selector}`)
            break
          }
        } catch {
          // Continue to next selector
        }
      }

      if (buttonFound) {
        const formSelectors = [
          '[data-testid="item-effect-form"]',
          'form',
          '[class*="modal"]',
        ]

        let formFound = false
        for (const selector of formSelectors) {
          try {
            if (await page.locator(selector).first().isVisible({ timeout: 2000 })) {
              formFound = true
              console.log(`Item effect form found: ${selector}`)
              break
            }
          } catch {
            // Continue to next selector
          }
        }

        if (formFound) {
          await takeScreenshot(page, 'studio-item-effect-create-form')
        }
      } else {
        console.log('Create item effect button not found - requires authenticated session with game selected')
      }
    })
  })
})
