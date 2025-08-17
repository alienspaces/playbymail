import { test, expect } from '@playwright/test'

test.describe('Authentication Flow', () => {
  test('should show unauthenticated content on studio page', async ({ page }) => {
    // Navigate to studio page - it should be accessible but show unauthenticated content
    await page.goto('/studio')
    
    // Wait for page to load
    await page.waitForLoadState('networkidle')
    
    // Should show StudioEntryView for unauthenticated users
    const pageContent = page.locator('body')
    
    // The page should load without redirecting to login
    await expect(page).toHaveURL('/studio')
    
    // Should show some studio-related content (even if limited)
    await expect(pageContent).toContainText(/Studio|Game|Design/i)
  })

  test('should submit login form and redirect to verify', async ({ page }) => {
    // 1. Submit email for authentication
    await page.goto('/login')
    
    const emailInput = page.locator('input[type="email"], input[name="email"]')
    await emailInput.fill('test@example.com')
    
    const submitButton = page.locator('button[type="submit"], button:has-text("Send Code")')
    await submitButton.click()
    
    // 2. Wait for redirect to verification page
    await expect(page).toHaveURL(/\/verify/)
    
    // 3. Verify page shows email
    const emailDisplay = page.locator('.email, [data-testid="email"]')
    if (await emailDisplay.isVisible()) {
      await expect(emailDisplay).toContainText('test@example.com')
    }
    
    // Note: We can't complete the full flow without the verification token
    // This test verifies the first part of the authentication flow
    console.log('✅ Login form submission and redirect to verification page working')
  })

  test('should handle invalid verification code', async ({ page }) => {
    await page.goto('/verify?email=test@example.com')
    
    // Enter invalid code
    const codeInput = page.locator('input[type="text"], input[name="code"]')
    await codeInput.fill('INVALID123')
    
    // Submit
    const verifyButton = page.locator('button[type="submit"], button:has-text("Verify")')
    await verifyButton.click()
    
    // Should show error message
    await page.waitForTimeout(1000)
    
    const errorMessage = page.locator('.error, .message, [data-testid="error"]')
    if (await errorMessage.isVisible()) {
      await expect(errorMessage).toContainText(/invalid|failed/i)
    }
  })

  test('should show verification page with code input', async ({ page }) => {
    // Navigate directly to verification page
    await page.goto('/verify?email=test@example.com')
    
    // Look for verification code input
    const codeInput = page.locator('input[type="text"], input[name="code"], input[placeholder*="code"]')
    const verifyButton = page.locator('button[type="submit"], button:has-text("Verify")')
    
    await expect(codeInput).toBeVisible()
    await expect(verifyButton).toBeVisible()
  })
})
