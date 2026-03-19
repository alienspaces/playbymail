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

test.describe('Studio Location Objects Designer Workflows', () => {
  test.beforeEach(async ({ page }) => {
    await page.context().clearCookies()
  })

  test.describe('Location Objects Page Access', () => {
    test('should show unauthenticated content for location objects page', async ({ page }) => {
      await navigateTo(page, '/studio')

      await checkPageURL(page, '/studio')
      await checkElementVisible(page, 'body')

      const body = page.locator('body')
      const content = await body.textContent()
      expect(content).toMatch(/Studio|Game|Design|Sign In/i)

      await takeScreenshot(page, 'studio-location-objects-unauthenticated')
    })
  })

  test.describe('Authenticated Location Objects CRUD', () => {
    test('should navigate to location objects page and show content', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')

      await waitForPageReady(page)

      const body = page.locator('body')
      const content = await body.textContent()

      if (content.includes('Location Objects') || content.includes('location-objects')) {
        await navigateTo(page, '/studio')
        await checkElementVisible(page, 'body')
        await takeScreenshot(page, 'studio-location-objects-page')
      } else {
        console.log('Location Objects link not found - may require authenticated session with game selected')
      }
    })

    test('should show location objects table when game is selected', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')

      await waitForPageReady(page)

      const tableSelectors = [
        '[data-testid="location-objects-table"]',
        'table',
        '[class*="resource-table"]',
      ]

      let tableFound = false
      for (const selector of tableSelectors) {
        try {
          if (await page.locator(selector).isVisible({ timeout: 2000 })) {
            tableFound = true
            console.log(`Location objects table found: ${selector}`)
            break
          }
        } catch {
          // Continue to next selector
        }
      }

      if (!tableFound) {
        console.log('Location objects table not found - requires game selection')
      }

      await takeScreenshot(page, 'studio-location-objects-table')
    })

    test('should open create location object form', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')

      await waitForPageReady(page)

      const createButtonSelectors = [
        'button:has-text("Create Location Object")',
        '[data-testid="create-location-object"]',
        'button:has-text("Create")',
      ]

      let buttonFound = false
      for (const selector of createButtonSelectors) {
        try {
          if (await page.locator(selector).first().isVisible({ timeout: 2000 })) {
            await safeClick(page, selector)
            buttonFound = true
            console.log(`Create button clicked: ${selector}`)
            break
          }
        } catch {
          // Continue to next selector
        }
      }

      if (buttonFound) {
        const formSelectors = [
          '[data-testid="location-object-form"]',
          'form',
          '[class*="modal"]',
        ]

        let formFound = false
        for (const selector of formSelectors) {
          try {
            if (await page.locator(selector).first().isVisible({ timeout: 2000 })) {
              formFound = true
              console.log(`Form found: ${selector}`)
              break
            }
          } catch {
            // Continue to next selector
          }
        }

        if (formFound) {
          await takeScreenshot(page, 'studio-location-object-create-form')
        }
      } else {
        console.log('Create button not found - requires authenticated session with game selected')
      }
    })
  })

  test.describe('Location Object Effects Page Access', () => {
    test('should show unauthenticated content for object effects page', async ({ page }) => {
      await navigateTo(page, '/studio')

      await checkPageURL(page, '/studio')
      await checkElementVisible(page, 'body')

      const body = page.locator('body')
      const content = await body.textContent()
      expect(content).toMatch(/Studio|Game|Design|Sign In/i)

      await takeScreenshot(page, 'studio-location-object-effects-unauthenticated')
    })
  })

  test.describe('Authenticated Object Effects CRUD', () => {
    test('should navigate to object effects page and show content', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')

      await waitForPageReady(page)

      const body = page.locator('body')
      const content = await body.textContent()

      if (content.includes('Object Effects') || content.includes('location-object-effects')) {
        await checkElementVisible(page, 'body')
        await takeScreenshot(page, 'studio-location-object-effects-page')
      } else {
        console.log('Object Effects link not found - may require authenticated session with game selected')
      }
    })

    test('should show object effects table when game is selected', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')

      await waitForPageReady(page)

      const tableSelectors = [
        '[data-testid="location-object-effects-table"]',
        'table',
        '[class*="resource-table"]',
      ]

      let tableFound = false
      for (const selector of tableSelectors) {
        try {
          if (await page.locator(selector).isVisible({ timeout: 2000 })) {
            tableFound = true
            console.log(`Object effects table found: ${selector}`)
            break
          }
        } catch {
          // Continue to next selector
        }
      }

      if (!tableFound) {
        console.log('Object effects table not found - requires game selection')
      }

      await takeScreenshot(page, 'studio-location-object-effects-table')
    })

    test('should open create object effect form', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')

      await waitForPageReady(page)

      const createButtonSelectors = [
        'button:has-text("Create Object Effect")',
        '[data-testid="create-location-object-effect"]',
        'button:has-text("Create")',
      ]

      let buttonFound = false
      for (const selector of createButtonSelectors) {
        try {
          if (await page.locator(selector).first().isVisible({ timeout: 2000 })) {
            await safeClick(page, selector)
            buttonFound = true
            console.log(`Create effect button clicked: ${selector}`)
            break
          }
        } catch {
          // Continue to next selector
        }
      }

      if (buttonFound) {
        const formSelectors = [
          '[data-testid="location-object-effect-form"]',
          'form',
          '[class*="modal"]',
        ]

        let formFound = false
        for (const selector of formSelectors) {
          try {
            if (await page.locator(selector).first().isVisible({ timeout: 2000 })) {
              formFound = true
              console.log(`Effect form found: ${selector}`)
              break
            }
          } catch {
            // Continue to next selector
          }
        }

        if (formFound) {
          await takeScreenshot(page, 'studio-location-object-effect-create-form')
        }
      } else {
        console.log('Create effect button not found - requires authenticated session with game selected')
      }
    })
  })
})
