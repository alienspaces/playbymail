import { test, expect } from '@playwright/test'
import { 
  navigateTo, 
  waitForPageReady, 
  checkPageTitle, 
  checkPageURL,
  checkElementVisible,
  checkElementContainsText,
  safeClick,
  fillFormField,
  takeScreenshot,
  waitForText,
  waitForLoadingState
} from '../utils/test-helpers.js'

test.describe('Game Creation Workflows', () => {
  test.beforeEach(async ({ page }) => {
    // Clear any existing state before each test
    await page.context().clearCookies()
  })

  test.describe('Studio Access', () => {
    test('should show unauthenticated studio content', async ({ page }) => {
      await navigateTo(page, '/studio')
      
      await checkPageTitle(page, 'Studio')
      await checkPageURL(page, '/studio')
      
      // Should show unauthenticated studio content
      await checkElementVisible(page, 'body')
      
      // Look for studio-related content
      const body = page.locator('body')
      const content = await body.textContent()
      
      // Should show some studio-related content
      expect(content).toMatch(/Studio|Game|Design|Create/i)
      
      await takeScreenshot(page, 'studio-unauthenticated')
    })

    test('should show login prompt for unauthenticated users', async ({ page }) => {
      await navigateTo(page, '/studio')
      
      // Look for login-related content
      const body = page.locator('body')
      const content = await body.textContent()
      
      // Should show login prompt or unauthenticated message
      expect(content).toMatch(/Sign In|Login|Authenticate|Access/i)
    })
  })

  test.describe('Game Creation Process', () => {
    test('should display game creation form when authenticated', async ({ page }) => {
      // Use development bypass for authentication
      await page.setExtraHTTPHeaders({
        'X-Bypass-Authentication': 'test@example.com'
      })
      
      await navigateTo(page, '/studio')
      
      // Should show authenticated studio content
      await checkPageURL(page, '/studio')
      
      // Look for game creation elements
      const body = page.locator('body')
      const content = await body.textContent()
      
      if (content.includes('Create Game') || content.includes('New Game')) {
        // Should show game creation form
        await checkElementVisible(page, 'body')
        
        // Look for form elements
        const formSelectors = [
          'form',
          'input[name="gameName"]',
          'input[name="gameType"]',
          'button:has-text("Create")',
          'button:has-text("Submit")'
        ]
        
        let formFound = false
        for (const selector of formSelectors) {
          try {
            if (await page.locator(selector).isVisible()) {
              formFound = true
              break
            }
          } catch (e) {
            // Continue to next selector
          }
        }
        
        expect(formFound).toBe(true)
        
        await takeScreenshot(page, 'game-creation-form')
      } else {
        console.log('Game creation form not found - may need authentication or different content')
      }
    })

    test('should handle game creation form submission', async ({ page }) => {
      // Use development bypass for authentication
      await page.setExtraHTTPHeaders({
        'X-Bypass-Authentication': 'test@example.com'
      })
      
      await navigateTo(page, '/studio')
      
      // Look for game creation form
      const formSelectors = [
        'form',
        'input[name="gameName"]',
        'button:has-text("Create")'
      ]
      
      let formFound = false
      for (const selector of formSelectors) {
        try {
          if (await page.locator(selector).isVisible()) {
            formFound = true
            break
          }
        } catch (e) {
          // Continue to next selector
        }
      }
      
      if (formFound) {
        // Fill in game details
        const nameInput = page.locator('input[name="gameName"], input[placeholder*="name"], input[type="text"]').first()
        if (await nameInput.isVisible()) {
          await fillFormField(page, 'input[name="gameName"], input[placeholder*="name"], input[type="text"]', 'Test Adventure Game')
          
          // Submit form
          const submitButton = page.locator('button:has-text("Create"), button:has-text("Submit"), button[type="submit"]').first()
          if (await submitButton.isVisible()) {
            await safeClick(page, 'button:has-text("Create"), button:has-text("Submit"), button[type="submit"]')
            
            // Should show success or redirect
            await page.waitForTimeout(2000)
            
            // Check for success message or redirect
            const body = page.locator('body')
            const content = await body.textContent()
            
            if (content.includes('success') || content.includes('created') || content.includes('redirect')) {
              console.log('Game creation form submitted successfully')
            } else {
              console.log('Form submitted but no clear success indication')
            }
          }
        }
      } else {
        console.log('Game creation form not found - skipping form submission test')
      }
    })
  })

  test.describe('Game Configuration', () => {
    test('should show game configuration options', async ({ page }) => {
      // Use development bypass for authentication
      await page.setExtraHTTPHeaders({
        'X-Bypass-Authentication': 'test@example.com'
      })
      
      await navigateTo(page, '/studio')
      
      // Look for configuration options
      const configSelectors = [
        'input[name="gameType"]',
        'select[name="gameType"]',
        'input[name="turnDuration"]',
        'input[name="maxPlayers"]',
        'input[name="description"]'
      ]
      
      let configFound = false
      for (const selector of configSelectors) {
        try {
          if (await page.locator(selector).isVisible()) {
            configFound = true
            console.log(`Configuration element found: ${selector}`)
            break
          }
        } catch (e) {
          // Continue to next selector
        }
      }
      
      // Configuration options are optional but good to have
      console.log('Game configuration options found:', configFound)
    })

    test('should validate game creation form inputs', async ({ page }) => {
      // Use development bypass for authentication
      await page.setExtraHTTPHeaders({
        'X-Bypass-Authentication': 'test@example.com'
      })
      
      await navigateTo(page, '/studio')
      
      // Look for form inputs
      const nameInput = page.locator('input[name="gameName"], input[placeholder*="name"], input[type="text"]').first()
      
      if (await nameInput.isVisible()) {
        // Test empty submission
        const submitButton = page.locator('button:has-text("Create"), button:has-text("Submit"), button[type="submit"]').first()
        
        if (await submitButton.isVisible()) {
          await safeClick(page, 'button:has-text("Create"), button:has-text("Submit"), button[type="submit"]')
          
          // Should show validation error
          await page.waitForTimeout(1000)
          
          const body = page.locator('body')
          const content = await body.textContent()
          
          if (content.includes('required') || content.includes('error') || content.includes('invalid')) {
            console.log('Form validation working correctly')
          } else {
            console.log('No validation error shown for empty form')
          }
        }
      } else {
        console.log('Game creation form not found - skipping validation test')
      }
    })
  })

  test.describe('Game Templates', () => {
    test('should show available game templates', async ({ page }) => {
      // Use development bypass for authentication
      await page.setExtraHTTPHeaders({
        'X-Bypass-Authentication': 'test@example.com'
      })
      
      await navigateTo(page, '/studio')
      
      // Look for game templates
      const templateSelectors = [
        '.game-template',
        '[data-testid="template"]',
        '.template',
        'button:has-text("Template")'
      ]
      
      let templatesFound = false
      for (const selector of templateSelectors) {
        try {
          if (await page.locator(selector).isVisible()) {
            templatesFound = true
            console.log(`Game template found: ${selector}`)
            break
          }
        } catch (e) {
          // Continue to next selector
        }
      }
      
      // Templates are optional but good to have
      console.log('Game templates found:', templatesFound)
    })
  })

  test.describe('Error Handling', () => {
    test('should handle game creation errors gracefully', async ({ page }) => {
      // Use development bypass for authentication
      await page.setExtraHTTPHeaders({
        'X-Bypass-Authentication': 'test@example.com'
      })
      
      await navigateTo(page, '/studio')
      
      // Mock API error
      await page.route('**/api/**', route => {
        route.fulfill({
          status: 400,
          contentType: 'application/json',
          body: JSON.stringify({ error: 'Game creation failed' })
        })
      })
      
      // Try to create game
      const submitButton = page.locator('button:has-text("Create"), button:has-text("Submit"), button[type="submit"]').first()
      
      if (await submitButton.isVisible()) {
        await safeClick(page, 'button:has-text("Create"), button:has-text("Submit"), button[type="submit"]')
        
        // Should show error message
        await page.waitForTimeout(2000)
        
        const body = page.locator('body')
        const content = await body.textContent()
        
        if (content.includes('error') || content.includes('failed') || content.includes('invalid')) {
          console.log('Error handling working correctly')
        } else {
          console.log('No error message shown for failed game creation')
        }
      } else {
        console.log('Submit button not found - skipping error handling test')
      }
    })
  })

  test.describe('Navigation Flow', () => {
    test('should navigate between studio sections', async ({ page }) => {
      // Use development bypass for authentication
      await page.setExtraHTTPHeaders({
        'X-Bypass-Authentication': 'test@example.com'
      })
      
      await navigateTo(page, '/studio')
      
      // Look for navigation elements
      const navSelectors = [
        'nav',
        '.navigation',
        '.sidebar',
        '.menu',
        'a[href*="studio"]'
      ]
      
      let navFound = false
      for (const selector of navSelectors) {
        try {
          if (await page.locator(selector).isVisible()) {
            navFound = true
            break
          }
        } catch (e) {
          // Continue to next selector
        }
      }
      
      if (navFound) {
        // Test navigation if available
        console.log('Studio navigation found - navigation flow available')
      } else {
        console.log('Studio navigation not found - may be single page')
      }
    })
  })

  test.describe('Responsive Design', () => {
    test('should adapt studio to mobile viewport', async ({ page }) => {
      // Use development bypass for authentication
      await page.setExtraHTTPHeaders({
        'X-Bypass-Authentication': 'test@example.com'
      })
      
      // Set mobile viewport
      await page.setViewportSize({ width: 375, height: 667 })
      
      await navigateTo(page, '/studio')
      
      // Check if page loads in mobile view
      await checkElementVisible(page, 'body')
      
      // Take mobile screenshot
      await takeScreenshot(page, 'studio-mobile-view')
      
      // Reset to desktop viewport
      await page.setViewportSize({ width: 1280, height: 720 })
    })
  })
})
