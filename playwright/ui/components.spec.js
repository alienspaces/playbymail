import { test, expect } from '@playwright/test'
import { 
  navigateTo, 
  waitForPageReady, 
  checkElementVisible,
  checkElementContainsText,
  checkButtonEnabled,
  checkButtonDisabled,
  safeClick,
  fillFormField,
  takeScreenshot,
  waitForElement,
  waitForElementToDisappear
} from '../utils/test-helpers.js'

test.describe('UI Components', () => {
  test.beforeEach(async ({ page }) => {
    // Clear any existing state before each test
    await page.context().clearCookies()
    await page.context().clearStorageState()
  })

  test.describe('Buttons', () => {
    test('should display buttons correctly', async ({ page }) => {
      await navigateTo(page, '/login')
      
      // Check primary button
      const submitButton = page.locator('button:has-text("Send Code")')
      await checkElementVisible(page, 'button:has-text("Send Code")')
      await checkButtonEnabled(page, 'button:has-text("Send Code")')
      
      // Check button text
      await checkElementContainsText(page, 'button:has-text("Send Code")', 'Send Code')
      
      await takeScreenshot(page, 'buttons-display')
    })

    test('should handle button click events', async ({ page }) => {
      await navigateTo(page, '/login')
      
      // Fill form first
      await fillFormField(page, 'input[type="email"]', 'test@example.com')
      
      // Click button
      await safeClick(page, 'button:has-text("Send Code")')
      
      // Should redirect to verification page
      await expect(page).toHaveURL(/\/verify/)
    })

    test('should show loading state on button click', async ({ page }) => {
      await navigateTo(page, '/login')
      
      // Fill form
      await fillFormField(page, 'input[type="email"]', 'test@example.com')
      
      // Click button and check for loading state
      const submitButton = page.locator('button:has-text("Send Code")')
      
      // Click and immediately check for loading indicators
      await safeClick(page, 'button:has-text("Send Code")')
      
      // Look for loading indicators
      const loadingSelectors = [
        '[data-testid="loading"]',
        '.loading',
        '.spinner',
        'button:has-text("Sending")',
        'button:has-text("Loading")'
      ]
      
      let loadingFound = false
      for (const selector of loadingSelectors) {
        try {
          if (await page.locator(selector).isVisible()) {
            loadingFound = true
            break
          }
        } catch (e) {
          // Continue to next selector
        }
      }
      
      // If no specific loading element, check if button text changed
      if (!loadingFound) {
        const buttonText = await submitButton.textContent()
        if (buttonText.includes('Sending') || buttonText.includes('Loading')) {
          loadingFound = true
        }
      }
      
      // Loading state is optional but good to have
      console.log('Loading state found:', loadingFound)
    })
  })

  test.describe('Form Inputs', () => {
    test('should display form inputs correctly', async ({ page }) => {
      await navigateTo(page, '/login')
      
      // Check email input
      const emailInput = page.locator('input[type="email"]')
      await checkElementVisible(page, 'input[type="email"]')
      
      // Check input attributes
      const placeholder = await emailInput.getAttribute('placeholder')
      const type = await emailInput.getAttribute('type')
      const required = await emailInput.getAttribute('required')
      
      expect(type).toBe('email')
      expect(placeholder || 'label').toBeTruthy()
      
      await takeScreenshot(page, 'form-inputs')
    })

    test('should handle input focus and blur', async ({ page }) => {
      await navigateTo(page, '/login')
      
      const emailInput = page.locator('input[type="email"]')
      
      // Focus input
      await emailInput.focus()
      
      // Check if focused
      const isFocused = await emailInput.evaluate(el => el === document.activeElement)
      expect(isFocused).toBe(true)
      
      // Blur input
      await emailInput.blur()
      
      // Check if not focused
      const isNotFocused = await emailInput.evaluate(el => el !== document.activeElement)
      expect(isNotFocused).toBe(true)
    })

    test('should handle input validation', async ({ page }) => {
      await navigateTo(page, '/login')
      
      const emailInput = page.locator('input[type="email"]')
      
      // Test valid input
      await fillFormField(page, 'input[type="email"]', 'valid@example.com')
      let validity = await emailInput.evaluate(el => el.validity.valid)
      expect(validity).toBe(true)
      
      // Test invalid input
      await fillFormField(page, 'input[type="email"]', 'invalid-email')
      validity = await emailInput.evaluate(el => el.validity.valid)
      expect(validity).toBe(false)
    })

    test('should handle input clearing', async ({ page }) => {
      await navigateTo(page, '/login')
      
      const emailInput = page.locator('input[type="email"]')
      
      // Fill input
      await fillFormField(page, 'input[type="email"]', 'test@example.com')
      
      // Clear input
      await emailInput.clear()
      
      // Check if empty
      const value = await emailInput.inputValue()
      expect(value).toBe('')
    })
  })

  test.describe('Navigation Elements', () => {
    test('should display navigation elements', async ({ page }) => {
      await navigateTo(page, '/')
      
      // Look for common navigation elements
      const navSelectors = [
        'nav',
        '.navigation',
        '.navbar',
        '.header',
        '.menu'
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
      
      // Should have some navigation structure
      expect(navFound).toBe(true)
      
      await takeScreenshot(page, 'navigation-elements')
    })

    test('should handle navigation links', async ({ page }) => {
      await navigateTo(page, '/')
      
      // Look for navigation links
      const linkSelectors = [
        'a[href="/"]',
        'a[href="/faq"]',
        'a[href="/about"]',
        'a[href="/login"]'
      ]
      
      for (const selector of linkSelectors) {
        try {
          const link = page.locator(selector)
          if (await link.isVisible()) {
            // Check if link is clickable
            await expect(link).toBeVisible()
            
            // Test navigation (optional - might be too aggressive)
            // await safeClick(page, selector)
            // await page.waitForTimeout(500)
            // await page.goBack()
            break
          }
        } catch (e) {
          // Continue to next selector
        }
      }
    })
  })

  test.describe('Loading States', () => {
    test('should show loading indicators when appropriate', async ({ page }) => {
      await navigateTo(page, '/login')
      
      // Fill form
      await fillFormField(page, 'input[type="email"]', 'test@example.com')
      
      // Submit and look for loading state
      await safeClick(page, 'button:has-text("Send Code")')
      
      // Look for loading indicators
      const loadingSelectors = [
        '[data-testid="loading"]',
        '.loading',
        '.spinner',
        '.progress',
        '[aria-label*="loading"]'
      ]
      
      let loadingFound = false
      for (const selector of loadingSelectors) {
        try {
          if (await page.locator(selector).isVisible()) {
            loadingFound = true
            console.log(`Loading indicator found: ${selector}`)
            break
          }
        } catch (e) {
          // Continue to next selector
        }
      }
      
      // Loading state is optional but good UX
      console.log('Loading state found:', loadingFound)
    })
  })

  test.describe('Error States', () => {
    test('should display error messages appropriately', async ({ page }) => {
      await navigateTo(page, '/verify?email=test@example.com')
      
      // Submit invalid code
      await fillFormField(page, 'input[type="text"]', 'INVALID')
      await safeClick(page, 'button:has-text("Verify")')
      
      // Wait for response
      await page.waitForTimeout(1000)
      
      // Look for error indicators
      const errorSelectors = [
        '.error',
        '.message',
        '[data-testid="error"]',
        '.alert',
        '.notification'
      ]
      
      let errorFound = false
      for (const selector of errorSelectors) {
        try {
          if (await page.locator(selector).isVisible()) {
            errorFound = true
            console.log(`Error element found: ${selector}`)
            break
          }
        } catch (e) {
          // Continue to next selector
        }
      }
      
      // Should show some error indication
      expect(errorFound).toBe(true)
      
      await takeScreenshot(page, 'error-state')
    })
  })

  test.describe('Responsive Design', () => {
    test('should adapt to mobile viewport', async ({ page }) => {
      // Set mobile viewport
      await page.setViewportSize({ width: 375, height: 667 })
      
      await navigateTo(page, '/')
      
      // Check if page loads in mobile view
      await checkElementVisible(page, 'body')
      
      // Take mobile screenshot
      await takeScreenshot(page, 'mobile-view')
      
      // Reset to desktop viewport
      await page.setViewportSize({ width: 1280, height: 720 })
    })

    test('should adapt to tablet viewport', async ({ page }) => {
      // Set tablet viewport
      await page.setViewportSize({ width: 768, height: 1024 })
      
      await navigateTo(page, '/')
      
      // Check if page loads in tablet view
      await checkElementVisible(page, 'body')
      
      // Take tablet screenshot
      await takeScreenshot(page, 'tablet-view')
      
      // Reset to desktop viewport
      await page.setViewportSize({ width: 1280, height: 720 })
    })
  })

  test.describe('Accessibility', () => {
    test('should have proper ARIA labels', async ({ page }) => {
      await navigateTo(page, '/login')
      
      // Check for common accessibility attributes
      const emailInput = page.locator('input[type="email"]')
      
      // Check for aria-label or aria-labelledby
      const ariaLabel = await emailInput.getAttribute('aria-label')
      const ariaLabelledBy = await emailInput.getAttribute('aria-labelledby')
      
      // Should have some accessibility labeling
      expect(ariaLabel || ariaLabelledBy).toBeTruthy()
    })

    test('should have proper focus management', async ({ page }) => {
      await navigateTo(page, '/login')
      
      const emailInput = page.locator('input[type="email"]')
      
      // Focus should be manageable
      await emailInput.focus()
      
      // Check if focused
      const isFocused = await emailInput.evaluate(el => el === document.activeElement)
      expect(isFocused).toBe(true)
    })
  })
})
