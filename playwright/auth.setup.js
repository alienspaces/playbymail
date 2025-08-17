import { test as setup, expect } from '@playwright/test'

setup('setup test environment', async ({ context }) => {
  // This setup is now minimal since we're testing the real authentication flow
  // We'll authenticate through the normal flow in individual tests
  
  // Navigate to home page to verify the app loads
  const page = await context.newPage()
  await page.goto('/')
  
  // Wait for the page to load
  await page.waitForLoadState('networkidle')
  
  // Verify the app loads correctly
  const pageContent = page.locator('body')
  await expect(pageContent).toContainText(/Play by Mail|Welcome|Game/i)
  
  // Store the basic context state for other tests
  await context.storageState({ path: 'playwright/.auth/user.json' })
  
  await page.close()
})
