# Playwright End-to-End Tests

This directory contains Playwright tests for end-to-end testing of the PlayByMail application.

## Test Organization

### Directory Structure

```
playbymail/                    # Project root
├── playwright.config.js       # Playwright configuration
├── package.json               # Package configuration with test scripts
├── package-lock.json          # Dependency lock file
├── playwright/                # Test directory
│   ├── core/                  # Core functionality tests
│   │   ├── navigation.spec.js # Page navigation and routing
│   │   └── authentication.spec.js # Login and verification flows
│   ├── ui/                    # UI component tests
│   │   └── components.spec.js # Buttons, forms, responsive design
│   ├── workflows/             # User workflow tests
│   │   ├── game-creation.spec.js # Game creation flows
│   │   └── admin-workflows.spec.js # Admin dashboard workflows
│   ├── utils/                 # Test utilities and helpers
│   │   ├── test-helpers.js   # Common test functions
│   │   └── log-parser.js     # Log parsing utilities
│   ├── legacy/                # Old test files (for reference)
│   │   ├── admin-loading-bug.spec.js
│   │   ├── admin-navigation.spec.js
│   │   ├── auth-flow.spec.js
│   │   ├── simple-auth.spec.js
│   │   ├── game-creation.spec.js
│   │   ├── home.spec.js
│   │   └── log-monitoring.spec.js
│   ├── global-setup.js        # Global test setup
│   ├── global-teardown.js     # Global test cleanup
│   └── README.md              # This file
├── frontend/                  # Frontend application
└── backend/                   # Backend application
```

### Test Categories

- **Core**: Basic functionality, navigation, authentication
- **UI**: Component behavior, responsive design, accessibility
- **Workflows**: User journey testing, game creation, admin flows
- **Legacy**: Old test files maintained for reference

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
npm run test:e2e:ui-components
npm run test:e2e:workflows
npm run test:e2e:legacy
```

### Using npm Scripts

```bash
# Run all Playwright tests
npm run test:e2e

# Run specific categories
npm run test:e2e:core
npm run test:e2e:ui-components
npm run test:e2e:workflows
npm run test:e2e:legacy

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
npx playwright test legacy/

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
```

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
  await navigateTo(page, '/login')
  await fillFormField(page, 'input[type="email"]', 'test@example.com')
  await safeClick(page, 'button:has-text("Send Code")')
  await checkPageURL(page, /\/verify/)
  await takeScreenshot(page, 'form-submitted')
})
```

## Test Design

### Coverage
- Focus on user workflows, not implementation details
- Test what users see and do, not internal code structure
- Cover main user journeys and edge cases

### Maintainability
- Use common test utilities and helper functions
- Avoid hardcoded selectors when possible
- Keep tests independent and isolated

### Flexibility
- Tests adapt to UI changes gracefully
- Use multiple selector strategies for robustness
- Handle optional UI elements appropriately

## Testing Capabilities

### Browser Support
Tests run in multiple browsers by default:
- **Chromium** (Chrome/Edge)
- **Firefox**
- **WebKit** (Safari)
- **Mobile Chrome**
- **Mobile Safari**

### Responsive Testing
- **Mobile viewport** testing (375x667)
- **Tablet viewport** testing (768x1024)
- **Desktop viewport** testing (1280x720)

### Accessibility Testing
- **ARIA labels** and attributes
- **Focus management** testing
- **Semantic HTML** validation

## Debugging and Reporting

### Screenshots and Videos
- **Screenshots**: Automatically captured on test failure
- **Videos**: Recorded for failed tests
- **Traces**: Generated for retried tests

### Test Reports
```bash
# HTML report
npm run test:e2e:report

# JUnit report
npx playwright test --reporter=junit
```

### Debug Mode
```bash
# Interactive UI
npx playwright test --ui

# Visible browser
npx playwright test --headed

# Step-by-step debugging
npx playwright test --debug
```

## Error Handling

### Network Error Testing
```javascript
test('should handle network errors gracefully', async ({ page }) => {
  await page.route('**/api/**', route => {
    route.abort('failed')
  })
  
  await navigateTo(page, '/login')
  await fillFormField(page, 'input[type="email"]', 'test@example.com')
  await safeClick(page, 'button:has-text("Send Code")')
  
  await page.waitForTimeout(2000)
  const content = await page.textContent('body')
  expect(content).toMatch(/error|failed|network/i)
})
```

### Server Error Testing
```javascript
test('should handle server errors gracefully', async ({ page }) => {
  await page.route('**/api/**', route => {
    route.fulfill({
      status: 500,
      contentType: 'application/json',
      body: JSON.stringify({ error: 'Internal Server Error' })
    })
  })
  
  // Test error handling
})
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
The `playwright.config.js` configures:
- Test directories and patterns
- Browser projects and viewports
- Timeouts and retries
- Screenshot and video capture
- Global setup and teardown
- Web server for testing

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
    await page.context().clearCookies()
  })

  test('should do something specific', async ({ page }) => {
    await navigateTo(page, '/path')
    await checkElementVisible(page, '.selector')
    await safeClick(page, 'button')
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
2. **Test data missing**: Run `./tools/db-load-seed-data`
3. **Port conflicts**: Check if ports 3000 and 8080 are available
4. **Browser dependencies**: Run `npx playwright install-deps`

### Getting Help
- Check test output for detailed error messages
- Use `--headed` mode to see what's happening
- Review screenshots and videos in test output
- Check browser console for JavaScript errors
- Use `--ui` mode for interactive debugging

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

## Additional Resources

- [Playwright Documentation](https://playwright.dev/)
- [Playwright Testing Best Practices](https://playwright.dev/docs/best-practices)
- [Playwright Configuration](https://playwright.dev/docs/configuration)
- [Playwright API Reference](https://playwright.dev/docs/api/class-playwright)

---

**Note**: This test suite covers user workflows and critical functionality. Focus on testing what users do rather than implementation details.
