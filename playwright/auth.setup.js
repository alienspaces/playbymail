import { test as setup, expect } from '@playwright/test'
import { TEST_BYPASS_HEADER_NAME, TEST_BYPASS_HEADER_VALUE } from './utils/test-helpers.js'

setup('setup test environment', async ({ context }) => {
  // Set up test bypass headers for API requests
  // This allows using email as verification code in tests
  await context.setExtraHTTPHeaders({
    [TEST_BYPASS_HEADER_NAME]: TEST_BYPASS_HEADER_VALUE
  })

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
