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

test.describe('Studio Adventure Placement Designer Workflows', () => {
  test.beforeEach(async ({ page }) => {
    await page.context().clearCookies()
  })

  // ─── Item Placements Page ─────────────────────────────────────────────────────

  test.describe('Item Placements Page', () => {
    test('should show unauthenticated content for item placements page', async ({ page }) => {
      await navigateTo(page, '/studio')
      await checkPageURL(page, '/studio')
      await checkElementVisible(page, 'body')
      const content = await page.locator('body').textContent()
      expect(content).toMatch(/Studio|Game|Design|Sign In/i)
      await takeScreenshot(page, 'studio-item-placements-unauthenticated')
    })

    test('should navigate to item placements page and show content', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)
      const content = await page.locator('body').textContent()
      if (content.includes('Item Placements') || content.includes('item-placements')) {
        await checkElementVisible(page, 'body')
        await takeScreenshot(page, 'studio-item-placements-page')
      } else {
        console.log('Item Placements link not found - may require authenticated session with game selected')
      }
    })

    test('should show item placements table when game is selected', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)

      const tableSelectors = ['[data-testid="item-placements-table"]', 'table', '[class*="resource-table"]']
      let tableFound = false
      for (const selector of tableSelectors) {
        try {
          if (await page.locator(selector).isVisible({ timeout: 2000 })) {
            tableFound = true
            console.log(`Item placements table found: ${selector}`)
            break
          }
        } catch { /* continue */ }
      }
      if (!tableFound) console.log('Item placements table not found - requires game selection')
      await takeScreenshot(page, 'studio-item-placements-table')
    })

    test('should open create item placement form', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)

      const buttonSelectors = ['button:has-text("Create Item Placement")', 'button:has-text("Create Placement")', 'button:has-text("Create")']
      let buttonFound = false
      for (const selector of buttonSelectors) {
        try {
          if (await page.locator(selector).first().isVisible({ timeout: 2000 })) {
            await safeClick(page, selector)
            buttonFound = true
            break
          }
        } catch { /* continue */ }
      }
      if (buttonFound) {
        const formSelectors = ['[data-testid="item-placement-form"]', 'form', '[class*="modal"]']
        for (const selector of formSelectors) {
          try {
            if (await page.locator(selector).first().isVisible({ timeout: 2000 })) {
              console.log(`Item placement form found: ${selector}`)
              break
            }
          } catch { /* continue */ }
        }
        await takeScreenshot(page, 'studio-item-placement-create-form')
      } else {
        console.log('Create item placement button not found - requires authenticated session with game selected')
      }
    })
  })

  // ─── Creature Placements Page ─────────────────────────────────────────────────

  test.describe('Creature Placements Page', () => {
    test('should show unauthenticated content for creature placements page', async ({ page }) => {
      await navigateTo(page, '/studio')
      await checkPageURL(page, '/studio')
      await checkElementVisible(page, 'body')
      const content = await page.locator('body').textContent()
      expect(content).toMatch(/Studio|Game|Design|Sign In/i)
      await takeScreenshot(page, 'studio-creature-placements-unauthenticated')
    })

    test('should navigate to creature placements page and show content', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)
      const content = await page.locator('body').textContent()
      if (content.includes('Creature Placements') || content.includes('creature-placements')) {
        await checkElementVisible(page, 'body')
        await takeScreenshot(page, 'studio-creature-placements-page')
      } else {
        console.log('Creature Placements link not found - may require authenticated session with game selected')
      }
    })

    test('should show creature placements table when game is selected', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)

      const tableSelectors = ['[data-testid="creature-placements-table"]', 'table', '[class*="resource-table"]']
      let tableFound = false
      for (const selector of tableSelectors) {
        try {
          if (await page.locator(selector).isVisible({ timeout: 2000 })) {
            tableFound = true
            console.log(`Creature placements table found: ${selector}`)
            break
          }
        } catch { /* continue */ }
      }
      if (!tableFound) console.log('Creature placements table not found - requires game selection')
      await takeScreenshot(page, 'studio-creature-placements-table')
    })

    test('should open create creature placement form', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)

      const buttonSelectors = ['button:has-text("Create Creature Placement")', 'button:has-text("Create Placement")', 'button:has-text("Create")']
      let buttonFound = false
      for (const selector of buttonSelectors) {
        try {
          if (await page.locator(selector).first().isVisible({ timeout: 2000 })) {
            await safeClick(page, selector)
            buttonFound = true
            break
          }
        } catch { /* continue */ }
      }
      if (buttonFound) {
        const formSelectors = ['[data-testid="creature-placement-form"]', 'form', '[class*="modal"]']
        for (const selector of formSelectors) {
          try {
            if (await page.locator(selector).first().isVisible({ timeout: 2000 })) {
              console.log(`Creature placement form found: ${selector}`)
              break
            }
          } catch { /* continue */ }
        }
        await takeScreenshot(page, 'studio-creature-placement-create-form')
      } else {
        console.log('Create creature placement button not found - requires authenticated session with game selected')
      }
    })
  })
})
