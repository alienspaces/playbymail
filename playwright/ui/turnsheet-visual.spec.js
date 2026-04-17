/**
 * Visual regression tests for turn sheet HTML rendering.
 *
 * Captures full-page screenshots of each turn sheet at desktop (1280x900),
 * mobile (375x812), and print-emulated viewports. Uses toHaveScreenshot() for
 * baseline comparison so regressions are detected automatically.
 *
 * Prerequisites:
 *   Generate HTML fixtures first: ./tools/render-turnsheets
 *   The HTML files are self-contained (base64 background images) so no server
 *   is required. Tests skip gracefully if a fixture file is missing.
 *
 * First run creates baselines in turnsheet-visual.spec.js-snapshots/.
 * Update baselines after intentional changes: npm run test:e2e:visual -- --update-snapshots
 */

import { test, expect } from '@playwright/test'
import { readFileSync, existsSync } from 'fs'
import { resolve } from 'path'

const TESTDATA = resolve('backend/internal/turnsheet/testdata')

const SHEETS = [
  { name: 'adventure_game_join_game',            file: 'adventure_game_join_game_turnsheet.html' },
  { name: 'adventure_game_location_choice',      file: 'adventure_game_location_choice_turnsheet.html' },
  { name: 'adventure_game_inventory_management', file: 'adventure_game_inventory_management_turnsheet.html' },
  { name: 'adventure_game_monster_encounter',    file: 'adventure_game_monster_encounter_turnsheet.html' },
  { name: 'mecha_join_game',                     file: 'mecha_join_game_turnsheet.html' },
  { name: 'mecha_squad_management',              file: 'mecha_squad_management_turnsheet.html' },
  { name: 'mecha_orders',                        file: 'mecha_orders_turnsheet.html' },
]

const VIEWPORTS = [
  { label: 'desktop', width: 1280, height: 900 },
  { label: 'mobile',  width: 375,  height: 812 },
]

for (const sheet of SHEETS) {
  const filePath = resolve(TESTDATA, sheet.file)

  for (const viewport of VIEWPORTS) {
    test(`${sheet.name} - ${viewport.label}`, async ({ page }) => {
      if (!existsSync(filePath)) {
        test.skip(true, `Missing fixture: ${sheet.file} -- run ./tools/render-turnsheets first`)
        return
      }

      const html = readFileSync(filePath, 'utf-8')

      await page.setViewportSize({ width: viewport.width, height: viewport.height })
      await page.setContent(html, { waitUntil: 'networkidle' })

      await expect(page).toHaveScreenshot(`${sheet.name}-${viewport.label}.png`, {
        fullPage: true,
      })
    })
  }

  test(`${sheet.name} - print`, async ({ page }) => {
    if (!existsSync(filePath)) {
      test.skip(true, `Missing fixture: ${sheet.file} -- run ./tools/render-turnsheets first`)
      return
    }

    const html = readFileSync(filePath, 'utf-8')

    await page.setViewportSize({ width: 1280, height: 900 })
    await page.setContent(html, { waitUntil: 'networkidle' })
    await page.emulateMedia({ media: 'print' })

    await expect(page).toHaveScreenshot(`${sheet.name}-print.png`, {
      fullPage: true,
    })
  })
}
