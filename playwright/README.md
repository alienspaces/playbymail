# Playwright End-to-End Tests

This directory contains Playwright tests for end-to-end testing of the PlayByMail application.

## Test Structure

- **`auth.setup.js`** - Authentication setup using development bypass
- **`home.spec.js`** - Basic home page tests
- **`game-creation.spec.js`** - Game creation flow tests
- **`auth-flow.spec.js`** - Real authentication flow tests (without bypass)
- **`utils/test-helpers.js`** - Common test utility functions

## Running Tests

### Prerequisites

1. **Backend must be running** with test data loaded
2. **Frontend must be accessible** at `http://localhost:3000`
3. **Test user account** must exist in the database

### Basic Test Commands

```bash
# Run all tests
npm run test:e2e

# Run tests with UI (interactive)
npm run test:e2e:ui

# Run tests in headed mode (see browser)
npm run test:e2e:headed

# Run tests in debug mode
npm run test:e2e:debug

# View test report
npm run test:e2e:report
```

### Running Specific Tests

```bash
# Run only home page tests
npx playwright test home.spec.js

# Run only authentication tests
npx playwright test auth-flow.spec.js

# Run tests matching a pattern
npx playwright test --grep "home page"
```

## Authentication Strategy

### Development Bypass (Default)

Tests use the development authentication bypass by setting the `X-Bypass-Authentication` header. This:

- ✅ **Fast execution** - no email delays
- ✅ **Real data** - uses actual account records
- ✅ **No frontend changes** - bypass happens at HTTP level
- ✅ **Consistent state** - same user for all tests

### Real Authentication Flow

The `auth-flow.spec.js` tests demonstrate testing the actual email verification flow:

- Login form submission
- Verification code handling
- Error scenarios
- Redirect behavior

## Test Configuration

### Environment Variables

```bash
# Override base URL for tests
TEST_BASE_URL=http://localhost:3000

# Run in CI mode
CI=true
```

### Browser Support

Tests run in multiple browsers:
- **Chromium** (Chrome/Edge)
- **Firefox**
- **WebKit** (Safari)

## Writing New Tests

### Basic Test Structure

```javascript
import { test, expect } from '@playwright/test'

test.describe('Feature Name', () => {
  test('should do something', async ({ page }) => {
    await page.goto('/path')
    await expect(page.locator('.selector')).toBeVisible()
  })
})
```

### Using Test Helpers

```javascript
import { waitForElement, safeFill } from '../utils/test-helpers.js'

test('should fill form', async ({ page }) => {
  await page.goto('/form')
  await safeFill(page, 'input[name="field"]', 'value')
})
```

### Authentication in Tests

Tests automatically run with authentication. To test unauthenticated behavior:

```javascript
test('should redirect unauthenticated users', async ({ page }) => {
  await page.context().clearCookies()
  await page.goto('/protected-page')
  await expect(page).toHaveURL('/login')
})
```

## Debugging Tests

### Screenshots and Videos

- **Screenshots**: Automatically captured on test failure
- **Videos**: Recorded for failed tests
- **Traces**: Generated for retried tests

### Debug Mode

```bash
npm run test:e2e:debug
```

This opens Playwright Inspector for step-by-step debugging.

### UI Mode

```bash
npm run test:e2e:ui
```

Interactive test runner with real-time feedback.

## Best Practices

1. **Use data-testid attributes** for reliable element selection
2. **Wait for network idle** before assertions
3. **Take screenshots** for debugging complex failures
4. **Test user workflows** not implementation details
5. **Keep tests independent** - no shared state between tests

## Troubleshooting

### Common Issues

1. **Backend not running**: Ensure `./tools/start-backend` is running
2. **Test data missing**: Run `./tools/db-load-test-data`
3. **Port conflicts**: Check if port 3000 is available
4. **Browser dependencies**: Run `npx playwright install-deps`

### Getting Help

- Check test output for detailed error messages
- Use `--headed` mode to see what's happening
- Review screenshots and videos in `playwright/`
- Check browser console for JavaScript errors
