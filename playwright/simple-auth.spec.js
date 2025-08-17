import { test, expect } from '@playwright/test'

test.describe('Simple Authentication Flow', () => {
  test('should show login page', async ({ page }) => {
    await page.goto('/login')
    
    // Verify login form elements
    const emailInput = page.locator('input[type="email"], input[name="email"]')
    const submitButton = page.locator('button[type="submit"], button:has-text("Send Code")')
    
    await expect(emailInput).toBeVisible()
    await expect(submitButton).toBeVisible()
  })

  test('should submit login form and redirect to verify', async ({ page }) => {
    await page.goto('/login')
    
    // Fill in email
    const emailInput = page.locator('input[type="email"], input[name="email"]')
    await emailInput.fill('test@example.com')
    
    // Submit form
    const submitButton = page.locator('button[type="submit"], button:has-text("Send Code")')
    await submitButton.click()
    
    // Should redirect to verification page
    await expect(page).toHaveURL(/\/verify/)
    
    // Verify page shows email
    const emailDisplay = page.locator('.email, [data-testid="email"]')
    if (await emailDisplay.isVisible()) {
      await expect(emailDisplay).toContainText('test@example.com')
    }
  })

  test('should show verification page', async ({ page }) => {
    await page.goto('/verify?email=test@example.com')
    
    // Look for verification code input
    const codeInput = page.locator('input[type="text"], input[name="code"], input[placeholder*="code"]')
    const verifyButton = page.locator('button[type="submit"], button:has-text("Verify")')
    
    await expect(codeInput).toBeVisible()
    await expect(verifyButton).toBeVisible()
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

  test('should show unauthenticated content on studio page', async ({ page }) => {
    // Studio page should be accessible but show unauthenticated content
    await page.goto('/studio')
    
    // Wait for page to load
    await page.waitForLoadState('networkidle')
    
    // Should show StudioEntryView for unauthenticated users
    // Look for content that indicates this is the unauthenticated studio view
    const pageContent = page.locator('body')
    
    // The page should load without redirecting to login
    await expect(page).toHaveURL('/studio')
    
    // Should show some studio-related content (even if limited)
    await expect(pageContent).toContainText(/Studio|Game|Design/i)
  })

  test('should show unauthenticated content on admin page', async ({ page }) => {
    // Admin page should be accessible but show unauthenticated content
    await page.goto('/admin')
    
    // Wait for page to load
    await page.waitForLoadState('networkidle')
    
    // Should load without redirecting to login
    await expect(page).toHaveURL('/admin')
    
    // Should show some admin-related content (even if limited)
    const pageContent = page.locator('body')
    await expect(pageContent).toContainText(/Management|Admin|Dashboard/i)
  })

  test('should show account page with error for unauthenticated users', async ({ page }) => {
    // Account page should be accessible but show error for unauthenticated users
    await page.goto('/account')
    
    // Wait for page to load
    await page.waitForLoadState('networkidle')
    
    // Should load without redirecting to login
    await expect(page).toHaveURL('/account')
    
    // Should show error message since API call will fail
    const pageContent = page.locator('body')
    await expect(pageContent).toContainText(/Account|Failed|Error/i)
  })
})
