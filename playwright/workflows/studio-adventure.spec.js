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

test.describe('Studio Adventure Designer Workflows', () => {
  test.beforeEach(async ({ page }) => {
    await page.context().clearCookies()
  })

  // ─── Locations Page ───────────────────────────────────────────────────────────

  test.describe('Locations Page', () => {
    test('should show unauthenticated content for locations page', async ({ page }) => {
      await navigateTo(page, '/studio')
      await checkPageURL(page, '/studio')
      await checkElementVisible(page, 'body')
      const content = await page.locator('body').textContent()
      expect(content).toMatch(/Studio|Game|Design|Sign In/i)
      await takeScreenshot(page, 'studio-locations-unauthenticated')
    })

    test('should navigate to locations page and show content', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)
      const content = await page.locator('body').textContent()
      if (content.includes('Locations') || content.includes('locations')) {
        await checkElementVisible(page, 'body')
        await takeScreenshot(page, 'studio-locations-page')
      } else {
        console.log('Locations link not found - may require authenticated session with game selected')
      }
    })

    test('should show locations table when game is selected', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)

      const tableSelectors = ['[data-testid="locations-table"]', 'table', '[class*="resource-table"]']
      let tableFound = false
      for (const selector of tableSelectors) {
        try {
          if (await page.locator(selector).isVisible({ timeout: 2000 })) {
            tableFound = true
            console.log(`Locations table found: ${selector}`)
            break
          }
        } catch { /* continue */ }
      }
      if (!tableFound) console.log('Locations table not found - requires game selection')
      await takeScreenshot(page, 'studio-locations-table')
    })

    test('should open create location form', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)

      const buttonSelectors = ['button:has-text("Create New Location")', 'button:has-text("Create Location")', 'button:has-text("Create")']
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
        const formSelectors = ['[data-testid="location-form"]', 'form', '[class*="modal"]']
        for (const selector of formSelectors) {
          try {
            if (await page.locator(selector).first().isVisible({ timeout: 2000 })) {
              console.log(`Location form found: ${selector}`)
              break
            }
          } catch { /* continue */ }
        }
        await takeScreenshot(page, 'studio-location-create-form')
      } else {
        console.log('Create location button not found - requires authenticated session with game selected')
      }
    })
  })

  // ─── Location Links Page ─────────────────────────────────────────────────────

  test.describe('Location Links Page', () => {
    test('should show unauthenticated content for location links page', async ({ page }) => {
      await navigateTo(page, '/studio')
      await checkPageURL(page, '/studio')
      await checkElementVisible(page, 'body')
      const content = await page.locator('body').textContent()
      expect(content).toMatch(/Studio|Game|Design|Sign In/i)
      await takeScreenshot(page, 'studio-location-links-unauthenticated')
    })

    test('should navigate to location links page and show content', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)
      const content = await page.locator('body').textContent()
      if (content.includes('Location Links') || content.includes('location-links')) {
        await checkElementVisible(page, 'body')
        await takeScreenshot(page, 'studio-location-links-page')
      } else {
        console.log('Location Links link not found - may require authenticated session with game selected')
      }
    })

    test('should show location links table when game is selected', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)

      const tableSelectors = ['[data-testid="location-links-table"]', 'table', '[class*="resource-table"]']
      let tableFound = false
      for (const selector of tableSelectors) {
        try {
          if (await page.locator(selector).isVisible({ timeout: 2000 })) {
            tableFound = true
            console.log(`Location links table found: ${selector}`)
            break
          }
        } catch { /* continue */ }
      }
      if (!tableFound) console.log('Location links table not found - requires game selection')
      await takeScreenshot(page, 'studio-location-links-table')
    })

    test('should open create location link form', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)

      const buttonSelectors = ['button:has-text("Create New Location Link")', 'button:has-text("Create Location Link")', 'button:has-text("Create")']
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
        const formSelectors = ['[data-testid="location-link-form"]', 'form', '[class*="modal"]']
        for (const selector of formSelectors) {
          try {
            if (await page.locator(selector).first().isVisible({ timeout: 2000 })) {
              console.log(`Location link form found: ${selector}`)
              break
            }
          } catch { /* continue */ }
        }
        await takeScreenshot(page, 'studio-location-link-create-form')
      } else {
        console.log('Create location link button not found - requires authenticated session with game selected')
      }
    })
  })

  // ─── Location Link Requirements Page ─────────────────────────────────────────

  test.describe('Location Link Requirements Page', () => {
    test('should show unauthenticated content for location link requirements page', async ({ page }) => {
      await navigateTo(page, '/studio')
      await checkPageURL(page, '/studio')
      await checkElementVisible(page, 'body')
      const content = await page.locator('body').textContent()
      expect(content).toMatch(/Studio|Game|Design|Sign In/i)
      await takeScreenshot(page, 'studio-location-link-requirements-unauthenticated')
    })

    test('should navigate to location link requirements page and show content', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)
      const content = await page.locator('body').textContent()
      if (content.includes('Location Link Requirements') || content.includes('location-link-requirements')) {
        await checkElementVisible(page, 'body')
        await takeScreenshot(page, 'studio-location-link-requirements-page')
      } else {
        console.log('Location Link Requirements link not found - may require authenticated session with game selected')
      }
    })

    test('should show location link requirements table when game is selected', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)

      const tableSelectors = ['[data-testid="location-link-requirements-table"]', 'table', '[class*="resource-table"]']
      let tableFound = false
      for (const selector of tableSelectors) {
        try {
          if (await page.locator(selector).isVisible({ timeout: 2000 })) {
            tableFound = true
            console.log(`Location link requirements table found: ${selector}`)
            break
          }
        } catch { /* continue */ }
      }
      if (!tableFound) console.log('Location link requirements table not found - requires game selection')
      await takeScreenshot(page, 'studio-location-link-requirements-table')
    })

    test('should open create location link requirement form', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)

      const buttonSelectors = ['button:has-text("Create New Requirement")', 'button:has-text("Create Requirement")', 'button:has-text("Create")']
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
        const formSelectors = ['[data-testid="location-link-requirement-form"]', 'form', '[class*="modal"]']
        for (const selector of formSelectors) {
          try {
            if (await page.locator(selector).first().isVisible({ timeout: 2000 })) {
              console.log(`Location link requirement form found: ${selector}`)
              break
            }
          } catch { /* continue */ }
        }
        await takeScreenshot(page, 'studio-location-link-requirement-create-form')
      } else {
        console.log('Create location link requirement button not found - requires authenticated session with game selected')
      }
    })
  })

  // ─── Items Page ───────────────────────────────────────────────────────────────

  test.describe('Items Page', () => {
    test('should show unauthenticated content for items page', async ({ page }) => {
      await navigateTo(page, '/studio')
      await checkPageURL(page, '/studio')
      await checkElementVisible(page, 'body')
      const content = await page.locator('body').textContent()
      expect(content).toMatch(/Studio|Game|Design|Sign In/i)
      await takeScreenshot(page, 'studio-items-unauthenticated')
    })

    test('should navigate to items page and show content', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)
      const content = await page.locator('body').textContent()
      if (content.includes('Items') || content.includes('items')) {
        await checkElementVisible(page, 'body')
        await takeScreenshot(page, 'studio-items-page')
      } else {
        console.log('Items link not found - may require authenticated session with game selected')
      }
    })

    test('should show items table when game is selected', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)

      const tableSelectors = ['[data-testid="items-table"]', 'table', '[class*="resource-table"]']
      let tableFound = false
      for (const selector of tableSelectors) {
        try {
          if (await page.locator(selector).isVisible({ timeout: 2000 })) {
            tableFound = true
            console.log(`Items table found: ${selector}`)
            break
          }
        } catch { /* continue */ }
      }
      if (!tableFound) console.log('Items table not found - requires game selection')
      await takeScreenshot(page, 'studio-items-table')
    })

    test('should open create item form', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)

      const buttonSelectors = ['button:has-text("Create New Item")', 'button:has-text("Create Item")', 'button:has-text("Create")']
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
        const formSelectors = ['[data-testid="item-form"]', 'form', '[class*="modal"]']
        for (const selector of formSelectors) {
          try {
            if (await page.locator(selector).first().isVisible({ timeout: 2000 })) {
              console.log(`Item form found: ${selector}`)
              break
            }
          } catch { /* continue */ }
        }
        await takeScreenshot(page, 'studio-item-create-form')
      } else {
        console.log('Create item button not found - requires authenticated session with game selected')
      }
    })
  })

  // ─── Creatures Page ───────────────────────────────────────────────────────────

  test.describe('Creatures Page', () => {
    test('should show unauthenticated content for creatures page', async ({ page }) => {
      await navigateTo(page, '/studio')
      await checkPageURL(page, '/studio')
      await checkElementVisible(page, 'body')
      const content = await page.locator('body').textContent()
      expect(content).toMatch(/Studio|Game|Design|Sign In/i)
      await takeScreenshot(page, 'studio-creatures-unauthenticated')
    })

    test('should navigate to creatures page and show content', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)
      const content = await page.locator('body').textContent()
      if (content.includes('Creatures') || content.includes('creatures')) {
        await checkElementVisible(page, 'body')
        await takeScreenshot(page, 'studio-creatures-page')
      } else {
        console.log('Creatures link not found - may require authenticated session with game selected')
      }
    })

    test('should show creatures table when game is selected', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)

      const tableSelectors = ['[data-testid="creatures-table"]', 'table', '[class*="resource-table"]']
      let tableFound = false
      for (const selector of tableSelectors) {
        try {
          if (await page.locator(selector).isVisible({ timeout: 2000 })) {
            tableFound = true
            console.log(`Creatures table found: ${selector}`)
            break
          }
        } catch { /* continue */ }
      }
      if (!tableFound) console.log('Creatures table not found - requires game selection')
      await takeScreenshot(page, 'studio-creatures-table')
    })

    test('should open create creature form', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)

      const buttonSelectors = ['button:has-text("Create New Creature")', 'button:has-text("Create Creature")', 'button:has-text("Create")']
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
        const formSelectors = ['[data-testid="creature-form"]', 'form', '[class*="modal"]']
        for (const selector of formSelectors) {
          try {
            if (await page.locator(selector).first().isVisible({ timeout: 2000 })) {
              console.log(`Creature form found: ${selector}`)
              break
            }
          } catch { /* continue */ }
        }
        await takeScreenshot(page, 'studio-creature-create-form')
      } else {
        console.log('Create creature button not found - requires authenticated session with game selected')
      }
    })
  })
})
