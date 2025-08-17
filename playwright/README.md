# Playwright End-to-End Tests

This directory contains Playwright tests for end-to-end testing of the PlayByMail application. The tests cover core functionality, UI components, and user workflows.

## Test Organization

### Directory Structure

```
playwright/
├── core/                    # Core functionality tests
│   ├── navigation.spec.js   # Page navigation and routing
│   └── authentication.spec.js # Login and verification flows
├── ui/                      # UI component tests
│   └── components.spec.js   # Buttons, forms, responsive design
├── workflows/               # User workflow tests
│   ├── game-creation.spec.js # Game creation flows
│   └── admin-workflows.spec.js # Admin dashboard workflows
├── integration/             # Integration tests (future)
│   ├── api.spec.js         # API integration testing
│   └── database.spec.js    # Database integration testing
├── utils/                   # Test utilities and helpers
│   ├── test-helpers.js     # Common test functions
│   └── log-parser.js       # Log parsing utilities
├── legacy/                  # Existing test files (for reference)
├── playwright.config.js    # Playwright configuration
└── README.md               # This file
```

### Test Categories

- **Core**: Basic functionality, navigation, authentication
- **UI**: Component behavior, responsive design, accessibility
- **Workflows**: User journey testing, game creation, admin flows
- **Integration**: API testing, database operations
- **Legacy**: Existing tests (maintained for compatibility)

## Running Tests

### Prerequisites

1. **Backend must be running** with test data loaded
2. **Frontend must be accessible** at `http://localhost:3000`
3. **Test user account** must exist in the database

### Quick Start

```bash
# Start backend with test data
cd playbymail
./tools/start-backend

# Run all tests
npm run test:e2e

# Run specific test category
npm run test:e2e:core
npm run test:e2e:ui
npm run test:e2e:workflows
```

### Using npm Scripts

Package.json provides convenient test execution:

```bash
# Run all Playwright tests
npm run test:e2e

# Run specific categories
npm run test:e2e:core
npm run test:e2e:ui
npm run test:e2e:workflows

# Run with options
npm run test:e2e:ui      # Interactive UI
npm run test:e2e:headed  # Visible browser
npm run test:e2e:debug   # Debug mode
```

### Direct Playwright Commands

```bash
# Run all tests
npx playwright test

# Run specific test category
npx playwright test core/
npx playwright test ui/
npx playwright test workflows/

# Run specific test file
npx playwright test core/navigation.spec.js

# Run tests with UI (interactive)
npx playwright test --ui

# Run tests in headed mode (see browser)
npx playwright test --headed

# Run tests in debug mode
npx playwright test --debug

# Run tests matching pattern
npx playwright test --grep "login"

# Run tests for specific browser
npx playwright test --project chromium

# Run tests with specific file pattern
npx playwright test "**/*.spec.js"
npx playwright test "**/core/**/*.spec.js"
```

## Test Design

### 1. Coverage
- Focus on user workflows, not implementation details
- Test what users see and do, not internal code structure
- Cover main user journeys and edge cases

### 2. Maintainability
- Use common test utilities and helpers
- Avoid hardcoded selectors when possible
- Keep tests independent and isolated

### 3. Reusability
- Common setup and teardown procedures
- Shared test utilities and helper functions
- Consistent test patterns across suites

### 4. Flexibility
- Tests adapt to UI changes gracefully
- Use multiple selector strategies for robustness
- Handle optional UI elements appropriately

## Test Utilities

### Common Helper Functions

```javascript
import { 
  navigateTo,           // Navigate to page and wait for ready
  waitForPageReady,      // Wait for page to be fully loaded
  checkElementVisible,   // Verify element is visible
  checkElementContainsText, // Verify element contains text
  safeClick,            // Click element safely with scrolling
  fillFormField,        // Fill form field safely
  takeScreenshot,       // Take screenshot for debugging
  waitForText,          // Wait for text to appear
  checkPageTitle,       // Verify page title
  checkPageURL          // Verify page URL
} from '../utils/test-helpers.js'
```

### Using Test Helpers

```javascript
test('should handle form submission', async ({ page }) => {
  // Navigate and wait for page to be ready
  await navigateTo(page, '/login')
  
  // Fill form field safely
  await fillFormField(page, 'input[type="email"]', 'test@example.com')
  
  // Click button safely
  await safeClick(page, 'button:has-text("Send Code")')
  
  // Verify navigation
  await checkPageURL(page, /\/verify/)
  
  // Take screenshot for debugging
  await takeScreenshot(page, 'form-submitted')
})
```

## Testing Responsiveness

### Viewport Testing

```javascript
test('should adapt to mobile viewport', async ({ page }) => {
  // Set mobile viewport
  await page.setViewportSize({ width: 375, height: 667 })
  
  await navigateTo(page, '/')
  await takeScreenshot(page, 'mobile-view')
  
  // Reset to desktop
  await page.setViewportSize({ width: 1280, height: 720 })
})
```

### Browser Testing

Tests run in multiple browsers by default:
- **Chromium** (Chrome/Edge)
- **Firefox**
- **WebKit** (Safari)
- **Mobile Chrome**
- **Mobile Safari**

## Debugging Tests

### Screenshots and Videos

- **Screenshots**: Automatically captured on test failure
- **Videos**: Recorded for failed tests
- **Traces**: Generated for retried tests

### Debug Mode

```bash
# Run single test in debug mode
npx playwright test --debug core/navigation.spec.js

# Run with UI for step-by-step debugging
npx playwright test --ui
```

### Common Debugging Patterns

```javascript
test('should debug element visibility', async ({ page }) => {
  await navigateTo(page, '/login')
  
  // Debug: log page content
  const content = await page.textContent('body')
  console.log('Page content:', content)
  
  // Debug: take screenshot
  await takeScreenshot(page, 'debug-login-page')
  
  // Debug: check element state
  const emailInput = page.locator('input[type="email"]')
  console.log('Email input visible:', await emailInput.isVisible())
  console.log('Email input enabled:', await emailInput.isEnabled())
})
```

## Error Handling

### Network Error Testing

```javascript
test('should handle network errors gracefully', async ({ page }) => {
  // Block API calls to simulate network failure
  await page.route('**/api/**', route => {
    route.abort('failed')
  })
  
  // Test error handling
  await navigateTo(page, '/login')
  await fillFormField(page, 'input[type="email"]', 'test@example.com')
  await safeClick(page, 'button:has-text("Send Code")')
  
  // Should show error message
  await page.waitForTimeout(2000)
  const content = await page.textContent('body')
  expect(content).toMatch(/error|failed|network/i)
})
```

### Server Error Testing

```javascript
test('should handle server errors gracefully', async ({ page }) => {
  // Mock server error response
  await page.route('**/api/**', route => {
    route.fulfill({
      status: 500,
      contentType: 'application/json',
      body: JSON.stringify({ error: 'Internal Server Error' })
    })
  })
  
  // Test error handling
  // ... test implementation
})
```

## Test Reporting

### HTML Report

```bash
# Generate HTML report
npm run test:e2e:report

# Or use direct command
npx playwright show-report

# Or open from playwright-report directory
open playwright-report/index.html
```

### JUnit Report

```bash
# Generate JUnit XML report
npx playwright test --reporter=junit

# Report will be in test-results/results.xml
```

### Custom Reporting

```bash
# Multiple reporters
npx playwright test --reporter=html,junit,json

# Custom reporter configuration
npx playwright test --reporter=html --reporter=junit
```

## Configuration

### Environment Variables

```bash
# Override base URL for tests
export TEST_BASE_URL=http://localhost:3000

# Run in CI mode
export CI=true

# Custom timeout
export PLAYWRIGHT_TIMEOUT=30000
```

### Playwright Config

The `playwright.config.js` file configures:
- Test directories and patterns
- Browser projects and viewports
- Timeouts and retries
- Screenshot and video capture
- Global setup and teardown
- Web server for testing

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Playwright Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: 18
      - run: npm ci
      - run: npx playwright install-deps
      - run: npx playwright test
      - uses: actions/upload-artifact@v3
        if: always()
        with:
          name: playwright-report
          path: playwright-report/
```

## Writing New Tests

### Test Structure

```javascript
import { test, expect } from '@playwright/test'
import { 
  navigateTo, 
  checkElementVisible,
  safeClick 
} from '../utils/test-helpers.js'

test.describe('Feature Name', () => {
  test.beforeEach(async ({ page }) => {
    // Setup for each test
    await page.context().clearCookies()
  })

  test('should do something specific', async ({ page }) => {
    // Test implementation
    await navigateTo(page, '/path')
    await checkElementVisible(page, '.selector')
    await safeClick(page, 'button')
    
    // Assertions
    await expect(page).toHaveURL('/expected-path')
  })
})
```

### Best Practices

1. **Use descriptive test names** that explain what is being tested
2. **Keep tests independent** - no shared state between tests
3. **Use test helpers** for common operations
4. **Take screenshots** for debugging complex failures
5. **Test user workflows** not implementation details
6. **Handle errors gracefully** in tests
7. **Use appropriate timeouts** for different operations

## Troubleshooting

### Common Issues

1. **Backend not running**: Ensure `./tools/start-backend` is running
2. **Test data missing**: Run `./tools/db-load-test-data`
3. **Port conflicts**: Check if ports 3000 and 8080 are available
4. **Browser dependencies**: Run `npx playwright install-deps`

### Getting Help

- Check test output for detailed error messages
- Use `--headed` mode to see what's happening
- Review screenshots and videos in test output
- Check browser console for JavaScript errors
- Use `--ui` mode for interactive debugging

## Additional Resources

- [Playwright Documentation](https://playwright.dev/)
- [Playwright Testing Best Practices](https://playwright.dev/docs/best-practices)
- [Playwright Configuration](https://playwright.dev/docs/configuration)
- [Playwright API Reference](https://playwright.dev/docs/api/class-playwright)

---

**Note**: This test suite covers user workflows and critical functionality. Focus on testing what users do rather than implementation details.
