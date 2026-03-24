/**
 * Playwright configuration for visual regression tests.
 *
 * Loads static turn sheet HTML files directly via page.setContent() --
 * no backend or frontend server is needed.
 *
 * Usage:
 *   npm run test:e2e:visual
 *
 * Generate HTML fixtures first:
 *   ./tools/render-turnsheets
 *
 * Update baselines after intentional template changes:
 *   npm run test:e2e:visual -- --update-snapshots
 */

import { defineConfig, devices } from '@playwright/test'

export default defineConfig({
  testDir: './playwright/ui',
  testMatch: /.*turnsheet-visual.*/,
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: 0,
  reporter: 'line',
  timeout: 30000,
  expect: {
    timeout: 10000,
  },

  use: {
    trace: 'off',
    screenshot: 'off',
    video: 'off',
  },

  projects: [
    {
      name: 'visual',
      use: { ...devices['Desktop Chrome'] },
    },
  ],

  // No webServer -- visual tests use page.setContent() with static HTML files.
})
