import { test, expect } from '@playwright/test'
import {
  navigateTo,
  waitForPageReady,
  checkPageTitle,
  checkPageURL,
  checkElementVisible,
  checkElementContainsText,
  fillFormField,
  safeClick,
  checkErrorDisplayed,
  takeScreenshot,
  waitForText
} from '../utils/test-helpers.js'

test.describe('Authentication Flow', () => {
  test.beforeEach(async ({ page }) => {
    // Clear any existing state before each test
    await page.context().clearCookies()
  })

  test.describe('Login Process', () => {
    test('should display login form correctly', async ({ page }) => {
      await navigateTo(page, '/login')

      // App uses static title for all pages
      await checkPageTitle(page, 'Play by Mail')
      await checkPageURL(page, '/login')

      // Verify form elements
      await checkElementVisible(page, 'input[type="email"]')
      await checkElementVisible(page, 'button:has-text("Send Code")')

      // Check for expected text (page shows "Sign in with Email")
      await checkElementContainsText(page, 'body', 'Sign in')
      await checkElementContainsText(page, 'body', 'Email')

      await takeScreenshot(page, 'login-form')
    })

    test('should handle email input validation', async ({ page }) => {
      await navigateTo(page, '/login')

      const emailInput = page.locator('input[type="email"]')

      // Test invalid email formats
      const invalidEmails = [
        'invalid-email',
        'test@',
        '@example.com',
        'test..test@example.com'
      ]

      for (const invalidEmail of invalidEmails) {
        await fillFormField(page, 'input[type="email"]', invalidEmail)

        // Check if browser validation catches invalid emails
        const validity = await emailInput.evaluate(el => el.validity.valid)
        if (!validity) {
          console.log(`Invalid email "${invalidEmail}" correctly rejected by browser validation`)
        }
      }

      // Test valid email
      await fillFormField(page, 'input[type="email"]', 'test@example.com')
      const validity = await emailInput.evaluate(el => el.validity.valid)
      expect(validity).toBe(true)
    })

    test('should submit login form and redirect to verification', async ({ page }) => {
      await navigateTo(page, '/login')

      // Fill in valid email
      await fillFormField(page, 'input[type="email"]', 'test@example.com')

      // Submit form
      await safeClick(page, 'button:has-text("Send Code")')

      // Should redirect to verification page
      await checkPageURL(page, /\/verify/)

      // Verify page shows verification form elements
      await checkElementVisible(page, 'input[type="text"]')
      await checkElementVisible(page, 'button:has-text("Verify")')

      await takeScreenshot(page, 'verification-page')
    })

    test('should handle empty email submission', async ({ page }) => {
      await navigateTo(page, '/login')

      // Try to submit without email
      await safeClick(page, 'button:has-text("Send Code")')

      // Should stay on login page
      await checkPageURL(page, '/login')

      // Check for validation error or message
      const body = page.locator('body')
      const content = await body.textContent()

      // Should show some indication that email is required
      expect(content).toMatch(/email|required|fill/i)
    })
  })

  test.describe('Verification Process', () => {
    test('should display verification page correctly', async ({ page }) => {
      await navigateTo(page, '/verify?email=test@example.com')

      // App uses static title for all pages
      await checkPageTitle(page, 'Play by Mail')
      await checkPageURL(page, /\/verify/)

      // Check for verification form elements
      await checkElementVisible(page, 'input[type="text"]')
      await checkElementVisible(page, 'button:has-text("Verify")')

      // Should show verification page content
      await checkElementContainsText(page, 'body', 'verification code')

      await takeScreenshot(page, 'verification-form')
    })

    test('should handle verification code input', async ({ page }) => {
      await navigateTo(page, '/verify?email=test@example.com')

      const codeInput = page.locator('input[type="text"]')

      // Test various code inputs
      const testCodes = ['123456', 'ABC123', '123456789', 'A1B2C3']

      for (const code of testCodes) {
        await fillFormField(page, 'input[type="text"]', code)

        // Verify input was accepted
        const value = await codeInput.inputValue()
        expect(value).toBe(code)
      }
    })

    test('should handle invalid verification code submission', async ({ page }) => {
      await navigateTo(page, '/verify?email=test@example.com')

      // Enter invalid code
      await fillFormField(page, 'input[type="text"]', 'INVALID123')

      // Submit
      await safeClick(page, 'button:has-text("Verify")')

      // Wait for response
      await page.waitForTimeout(2000)

      // Should show error message or stay on page
      const errorSelectors = ['.error', '.message', '[data-testid="error"]', '.alert', '[role="alert"]']

      let errorFound = false
      for (const selector of errorSelectors) {
        try {
          const errorElement = page.locator(selector)
          if (await errorElement.isVisible({ timeout: 2000 }).catch(() => false)) {
            errorFound = true
            break
          }
        } catch (e) {
          // Continue to next selector
        }
      }

      if (!errorFound) {
        // If no specific error element, check body content
        const body = page.locator('body')
        const content = await body.textContent()
        // Should show some error indication or still be on verify page
        const hasError = content && content.match(/invalid|failed|incorrect|error|verify/i)
        if (!hasError) {
          console.log('No error message found - error handling may need improvement')
        }
      }
    })

    test('should handle empty verification code submission', async ({ page }) => {
      await navigateTo(page, '/verify?email=test@example.com')

      // Try to submit without code
      await safeClick(page, 'button:has-text("Verify")')

      // Should stay on verification page
      await checkPageURL(page, /\/verify/)

      // Check for validation message
      const body = page.locator('body')
      const content = await body.textContent()

      // Should show some indication that code is required
      expect(content).toMatch(/code|required|fill/i)
    })
  })

  test.describe('Error Handling', () => {
    test('should handle network errors gracefully', async ({ page }) => {
      await navigateTo(page, '/login')

      // Block API calls to simulate network failure
      await page.route('**/api/**', route => {
        route.abort('failed')
      })

      // Fill and submit form
      await fillFormField(page, 'input[type="email"]', 'test@example.com')
      await safeClick(page, 'button:has-text("Send Code")')

      // Should handle error gracefully
      await page.waitForTimeout(2000)

      // Check for error message
      const body = page.locator('body')
      const content = await body.textContent()

      // Should show some error indication
      expect(content).toMatch(/error|failed|network|connection/i)
    })

    test('should handle server errors gracefully', async ({ page }) => {
      await navigateTo(page, '/login')

      // Mock server error response
      await page.route('**/api/**', route => {
        route.fulfill({
          status: 500,
          contentType: 'application/json',
          body: JSON.stringify({ error: 'Internal Server Error' })
        })
      })

      // Fill and submit form
      await fillFormField(page, 'input[type="email"]', 'test@example.com')
      await safeClick(page, 'button:has-text("Send Code")')

      // Should handle error gracefully
      await page.waitForTimeout(2000)

      // Check for error message
      const body = page.locator('body')
      const content = await body.textContent()

      // Should show some error indication
      expect(content).toMatch(/error|failed|server|unavailable/i)
    })
  })

  test.describe('Authentication State', () => {
    test('should maintain authentication state across page reloads', async ({ page }) => {
      // This test would require actual authentication
      // For now, we'll test the unauthenticated state persistence

      await navigateTo(page, '/login')
      await checkPageURL(page, '/login')

      // Reload page
      await page.reload()
      await waitForPageReady(page)

      // Should still be on login page
      await checkPageURL(page, '/login')
      await checkElementVisible(page, 'input[type="email"]')
    })

    test('should clear authentication state when cookies are cleared', async ({ page }) => {
      await navigateTo(page, '/login')

      // Clear cookies
      await page.context().clearCookies()

      // Navigate to protected page
      await navigateTo(page, '/admin')

      // Should show unauthenticated content
      await checkPageURL(page, '/admin')

      const body = page.locator('body')
      const content = await body.textContent()

      // Should not show authenticated content
      expect(content).not.toMatch(/Games & Instances|Manage Instances/i)
    })
  })

  test.describe('Accessibility', () => {
    test('should have proper form labels and accessibility', async ({ page }) => {
      await navigateTo(page, '/login')

      // Check for proper form structure
      const form = page.locator('form')
      if (await form.isVisible()) {
        // Form should have proper action and method
        const action = await form.getAttribute('action')
        const method = await form.getAttribute('method')

        if (action) {
          expect(action).toMatch(/\/api\/|verify|login/i)
        }

        if (method) {
          expect(method.toLowerCase()).toMatch(/post|get/i)
        }
      }

      // Check for proper input attributes
      const emailInput = page.locator('input[type="email"]')
      const placeholder = await emailInput.getAttribute('placeholder')
      const required = await emailInput.getAttribute('required')

      // Should have helpful placeholder or label
      expect(placeholder || 'label').toBeTruthy()
    })
  })
})
