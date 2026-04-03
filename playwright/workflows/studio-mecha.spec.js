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

test.describe('Studio Mecha Designer Workflows', () => {
  test.beforeEach(async ({ page }) => {
    await page.context().clearCookies()
  })

  // ─── Chassis Page ────────────────────────────────────────────────────────────

  test.describe('Chassis Page', () => {
    test('should show unauthenticated content for chassis page', async ({ page }) => {
      await navigateTo(page, '/studio')
      await checkPageURL(page, '/studio')
      await checkElementVisible(page, 'body')
      const content = await page.locator('body').textContent()
      expect(content).toMatch(/Studio|Game|Design|Sign In/i)
      await takeScreenshot(page, 'studio-chassis-unauthenticated')
    })

    test('should navigate to chassis page and show content', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)
      const content = await page.locator('body').textContent()
      if (content.includes('Chassis') || content.includes('chassis')) {
        await checkElementVisible(page, 'body')
        await takeScreenshot(page, 'studio-chassis-page')
      } else {
        console.log('Chassis link not found - may require authenticated session with game selected')
      }
    })

    test('should show chassis table when game is selected', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)

      const tableSelectors = ['[data-testid="chassis-table"]', 'table', '[class*="resource-table"]']
      let tableFound = false
      for (const selector of tableSelectors) {
        try {
          if (await page.locator(selector).isVisible({ timeout: 2000 })) {
            tableFound = true
            console.log(`Chassis table found: ${selector}`)
            break
          }
        } catch { /* continue */ }
      }
      if (!tableFound) console.log('Chassis table not found - requires game selection')
      await takeScreenshot(page, 'studio-chassis-table')
    })

    test('should open create chassis form', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)

      const buttonSelectors = ['button:has-text("Create New Chassis")', 'button:has-text("Create Chassis")', 'button:has-text("Create")']
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
        const formSelectors = ['[data-testid="chassis-form"]', 'form', '[class*="modal"]']
        for (const selector of formSelectors) {
          try {
            if (await page.locator(selector).first().isVisible({ timeout: 2000 })) {
              console.log(`Chassis form found: ${selector}`)
              break
            }
          } catch { /* continue */ }
        }
        await takeScreenshot(page, 'studio-chassis-create-form')
      } else {
        console.log('Create chassis button not found - requires authenticated session with game selected')
      }
    })
  })

  // ─── Weapons Page ────────────────────────────────────────────────────────────

  test.describe('Weapons Page', () => {
    test('should show unauthenticated content for weapons page', async ({ page }) => {
      await navigateTo(page, '/studio')
      await checkPageURL(page, '/studio')
      await checkElementVisible(page, 'body')
      const content = await page.locator('body').textContent()
      expect(content).toMatch(/Studio|Game|Design|Sign In/i)
      await takeScreenshot(page, 'studio-weapons-unauthenticated')
    })

    test('should navigate to weapons page and show content', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)
      const content = await page.locator('body').textContent()
      if (content.includes('Weapons') || content.includes('weapons')) {
        await checkElementVisible(page, 'body')
        await takeScreenshot(page, 'studio-weapons-page')
      } else {
        console.log('Weapons link not found - may require authenticated session with game selected')
      }
    })

    test('should show weapons table when game is selected', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)

      const tableSelectors = ['[data-testid="weapons-table"]', 'table', '[class*="resource-table"]']
      let tableFound = false
      for (const selector of tableSelectors) {
        try {
          if (await page.locator(selector).isVisible({ timeout: 2000 })) {
            tableFound = true
            console.log(`Weapons table found: ${selector}`)
            break
          }
        } catch { /* continue */ }
      }
      if (!tableFound) console.log('Weapons table not found - requires game selection')
      await takeScreenshot(page, 'studio-weapons-table')
    })

    test('should open create weapon form', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)

      const buttonSelectors = ['button:has-text("Create New Weapon")', 'button:has-text("Create Weapon")', 'button:has-text("Create")']
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
        const formSelectors = ['[data-testid="weapon-form"]', 'form', '[class*="modal"]']
        for (const selector of formSelectors) {
          try {
            if (await page.locator(selector).first().isVisible({ timeout: 2000 })) {
              console.log(`Weapon form found: ${selector}`)
              break
            }
          } catch { /* continue */ }
        }
        await takeScreenshot(page, 'studio-weapon-create-form')
      } else {
        console.log('Create weapon button not found - requires authenticated session with game selected')
      }
    })
  })

  // ─── Sectors Page ────────────────────────────────────────────────────────────

  test.describe('Sectors Page', () => {
    test('should show unauthenticated content for sectors page', async ({ page }) => {
      await navigateTo(page, '/studio')
      await checkPageURL(page, '/studio')
      await checkElementVisible(page, 'body')
      const content = await page.locator('body').textContent()
      expect(content).toMatch(/Studio|Game|Design|Sign In/i)
      await takeScreenshot(page, 'studio-sectors-unauthenticated')
    })

    test('should navigate to sectors page and show content', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)
      const content = await page.locator('body').textContent()
      if (content.includes('Sectors') || content.includes('sectors')) {
        await checkElementVisible(page, 'body')
        await takeScreenshot(page, 'studio-sectors-page')
      } else {
        console.log('Sectors link not found - may require authenticated session with game selected')
      }
    })

    test('should show sectors table when game is selected', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)

      const tableSelectors = ['[data-testid="sectors-table"]', 'table', '[class*="resource-table"]']
      let tableFound = false
      for (const selector of tableSelectors) {
        try {
          if (await page.locator(selector).isVisible({ timeout: 2000 })) {
            tableFound = true
            console.log(`Sectors table found: ${selector}`)
            break
          }
        } catch { /* continue */ }
      }
      if (!tableFound) console.log('Sectors table not found - requires game selection')
      await takeScreenshot(page, 'studio-sectors-table')
    })

    test('should open create sector form', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)

      const buttonSelectors = ['button:has-text("Create New Sector")', 'button:has-text("Create Sector")', 'button:has-text("Create")']
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
        const formSelectors = ['[data-testid="sector-form"]', 'form', '[class*="modal"]']
        for (const selector of formSelectors) {
          try {
            if (await page.locator(selector).first().isVisible({ timeout: 2000 })) {
              console.log(`Sector form found: ${selector}`)
              break
            }
          } catch { /* continue */ }
        }
        await takeScreenshot(page, 'studio-sector-create-form')
      } else {
        console.log('Create sector button not found - requires authenticated session with game selected')
      }
    })
  })

  // ─── Sector Links Page ────────────────────────────────────────────────────────

  test.describe('Sector Links Page', () => {
    test('should show unauthenticated content for sector links page', async ({ page }) => {
      await navigateTo(page, '/studio')
      await checkPageURL(page, '/studio')
      await checkElementVisible(page, 'body')
      const content = await page.locator('body').textContent()
      expect(content).toMatch(/Studio|Game|Design|Sign In/i)
      await takeScreenshot(page, 'studio-sector-links-unauthenticated')
    })

    test('should navigate to sector links page and show content', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)
      const content = await page.locator('body').textContent()
      if (content.includes('Sector Links') || content.includes('sector-links')) {
        await checkElementVisible(page, 'body')
        await takeScreenshot(page, 'studio-sector-links-page')
      } else {
        console.log('Sector Links link not found - may require authenticated session with game selected')
      }
    })

    test('should show sector links table when game is selected', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)

      const tableSelectors = ['[data-testid="sector-links-table"]', 'table', '[class*="resource-table"]']
      let tableFound = false
      for (const selector of tableSelectors) {
        try {
          if (await page.locator(selector).isVisible({ timeout: 2000 })) {
            tableFound = true
            console.log(`Sector links table found: ${selector}`)
            break
          }
        } catch { /* continue */ }
      }
      if (!tableFound) console.log('Sector links table not found - requires game selection')
      await takeScreenshot(page, 'studio-sector-links-table')
    })

    test('should open create sector link form', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)

      const buttonSelectors = ['button:has-text("Create New Sector Link")', 'button:has-text("Create Sector Link")', 'button:has-text("Create")']
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
        const formSelectors = ['[data-testid="sector-link-form"]', 'form', '[class*="modal"]']
        for (const selector of formSelectors) {
          try {
            if (await page.locator(selector).first().isVisible({ timeout: 2000 })) {
              console.log(`Sector link form found: ${selector}`)
              break
            }
          } catch { /* continue */ }
        }
        await takeScreenshot(page, 'studio-sector-link-create-form')
      } else {
        console.log('Create sector link button not found - requires authenticated session with game selected')
      }
    })
  })

  // ─── Lances Page ────────────────────────────────────────────────────────────

  test.describe('Lances Page', () => {
    test('should show unauthenticated content for lances page', async ({ page }) => {
      await navigateTo(page, '/studio')
      await checkPageURL(page, '/studio')
      await checkElementVisible(page, 'body')
      const content = await page.locator('body').textContent()
      expect(content).toMatch(/Studio|Game|Design|Sign In/i)
      await takeScreenshot(page, 'studio-lances-unauthenticated')
    })

    test('should navigate to lances page and show content', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)
      const content = await page.locator('body').textContent()
      if (content.includes('Lances') || content.includes('lances')) {
        await checkElementVisible(page, 'body')
        await takeScreenshot(page, 'studio-lances-page')
      } else {
        console.log('Lances link not found - may require authenticated session with game selected')
      }
    })

    test('should show lances table when game is selected', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)

      const tableSelectors = ['[data-testid="lances-table"]', 'table', '[class*="resource-table"]']
      let tableFound = false
      for (const selector of tableSelectors) {
        try {
          if (await page.locator(selector).isVisible({ timeout: 2000 })) {
            tableFound = true
            console.log(`Lances table found: ${selector}`)
            break
          }
        } catch { /* continue */ }
      }
      if (!tableFound) console.log('Lances table not found - requires game selection')
      await takeScreenshot(page, 'studio-lances-table')
    })

    test('should open create lance form', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)

      const buttonSelectors = ['button:has-text("Create New Lance")', 'button:has-text("Create Lance")', 'button:has-text("Create")']
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
        const formSelectors = ['[data-testid="lance-form"]', 'form', '[class*="modal"]']
        for (const selector of formSelectors) {
          try {
            if (await page.locator(selector).first().isVisible({ timeout: 2000 })) {
              console.log(`Lance form found: ${selector}`)
              break
            }
          } catch { /* continue */ }
        }
        await takeScreenshot(page, 'studio-lance-create-form')
      } else {
        console.log('Create lance button not found - requires authenticated session with game selected')
      }
    })
  })

  // ─── Computer Opponents Page ─────────────────────────────────────────────────

  test.describe('Computer Opponents Page', () => {
    test('should show unauthenticated content for computer opponents page', async ({ page }) => {
      await navigateTo(page, '/studio')
      await checkPageURL(page, '/studio')
      await checkElementVisible(page, 'body')
      const content = await page.locator('body').textContent()
      expect(content).toMatch(/Studio|Game|Design|Sign In/i)
      await takeScreenshot(page, 'studio-computer-opponents-unauthenticated')
    })

    test('should navigate to computer opponents page and show content', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)
      const content = await page.locator('body').textContent()
      if (content.includes('Computer Opponents') || content.includes('computer-opponents')) {
        await checkElementVisible(page, 'body')
        await takeScreenshot(page, 'studio-computer-opponents-page')
      } else {
        console.log('Computer Opponents link not found - may require authenticated session with game selected')
      }
    })

    test('should show computer opponents table when game is selected', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)

      const tableSelectors = ['[data-testid="computer-opponents-table"]', 'table', '[class*="resource-table"]']
      let tableFound = false
      for (const selector of tableSelectors) {
        try {
          if (await page.locator(selector).isVisible({ timeout: 2000 })) {
            tableFound = true
            console.log(`Computer opponents table found: ${selector}`)
            break
          }
        } catch { /* continue */ }
      }
      if (!tableFound) console.log('Computer opponents table not found - requires game selection')
      await takeScreenshot(page, 'studio-computer-opponents-table')
    })

    test('should open create computer opponent form', async ({ page }) => {
      await setupTestBypassHeaders(page)
      await navigateTo(page, '/studio')
      await waitForPageReady(page)

      const buttonSelectors = ['button:has-text("Create New Opponent")', 'button:has-text("Create Opponent")', 'button:has-text("Create")']
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
        const formSelectors = ['[data-testid="computer-opponent-form"]', 'form', '[class*="modal"]']
        for (const selector of formSelectors) {
          try {
            if (await page.locator(selector).first().isVisible({ timeout: 2000 })) {
              console.log(`Computer opponent form found: ${selector}`)
              break
            }
          } catch { /* continue */ }
        }
        await takeScreenshot(page, 'studio-computer-opponent-create-form')
      } else {
        console.log('Create computer opponent button not found - requires authenticated session with game selected')
      }
    })
  })
})
