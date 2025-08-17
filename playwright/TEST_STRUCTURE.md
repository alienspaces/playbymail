# Playwright Test Structure Overview

This document describes the organized Playwright test suite for the PlayByMail application.

## Architecture Overview

The test suite is organized into logical categories that provide coverage while maintaining flexibility for UI updates. Each category focuses on specific aspects of the application and can be run independently.

## Directory Structure

```
playwright/
├── core/                           # Core functionality tests
│   ├── navigation.spec.js         # Page navigation and routing
│   └── authentication.spec.js     # Login and verification flows
├── ui/                            # UI component tests
│   └── components.spec.js         # Buttons, forms, responsive design
├── workflows/                     # User workflow tests
│   ├── game-creation.spec.js     # Game creation flows
│   └── admin-workflows.spec.js   # Admin dashboard workflows
├── integration/                   # Integration tests (future)
│   ├── api.spec.js               # API integration testing
│   └── database.spec.js          # Database integration testing
├── utils/                         # Test utilities and helpers
│   ├── test-helpers.js           # Common test functions
│   └── log-parser.js             # Log parsing utilities
├── legacy/                        # Existing test files (for reference)
├── playwright.config.js          # Playwright configuration
├── global-setup.js               # Global test setup
├── global-teardown.js            # Global test cleanup
├── README.md                     # Comprehensive documentation
└── TEST_STRUCTURE.md             # This file
```

## Test Categories

### 1. Core Tests (`core/`)
**Purpose**: Test fundamental application functionality
**Coverage**: Navigation, authentication, basic page loading
**Maintainability**: High - focuses on core functionality that rarely changes

**Files**:
- `navigation.spec.js` - Page navigation, routing, browser behavior
- `authentication.spec.js` - Login flows, verification, error handling

**Key Features**:
- Basic page loading and navigation
- Authentication flow testing
- Error handling for network/server issues
- Browser compatibility testing

### 2. UI Tests (`ui/`)
**Purpose**: Test UI component behavior and responsiveness
**Coverage**: Buttons, forms, responsive design, accessibility
**Maintainability**: Medium - adapts to UI changes gracefully

**Files**:
- `components.spec.js` - Button behavior, form inputs, responsive design

**Key Features**:
- Component interaction testing
- Responsive design validation
- Accessibility testing
- Form validation testing

### 3. Workflow Tests (`workflows/`)
**Purpose**: Test complete user journeys and business processes
**Coverage**: Game creation, admin workflows, user management
**Maintainability**: Medium - focuses on user workflows rather than specific UI elements

**Files**:
- `game-creation.spec.js` - Studio access, game creation, configuration
- `admin-workflows.spec.js` - Admin dashboard, game management, navigation

**Key Features**:
- End-to-end user workflows
- Business process validation
- Navigation flow testing
- Error handling in workflows

### 4. Integration Tests (`integration/`)
**Purpose**: Test API integration and database operations
**Coverage**: Backend integration, data persistence, API contracts
**Maintainability**: High - focuses on integration points that are stable

**Files**:
- `api.spec.js` - API endpoint testing, response validation
- `database.spec.js` - Database operations, data consistency

**Key Features**:
- API contract validation
- Database operation testing
- Data consistency checks
- Performance testing

### 5. Legacy Tests (`legacy/`)
**Purpose**: Maintain existing test coverage during transition
**Coverage**: Existing test scenarios and edge cases
**Maintainability**: Low - will be gradually replaced by organized tests

**Files**:
- All existing `.spec.js` files in the root directory

## Test Utilities

### Common Helper Functions (`utils/test-helpers.js`)

The test suite provides helper functions that make tests more maintainable and readable:

**Navigation Helpers**:
- `navigateTo(page, path)` - Navigate and wait for page ready
- `waitForPageReady(page)` - Wait for page to be fully loaded
- `checkPageTitle(page, title)` - Verify page title
- `checkPageURL(page, url)` - Verify page URL

**Element Interaction**:
- `safeClick(page, selector)` - Click element safely with scrolling
- `fillFormField(page, selector, value)` - Fill form field safely
- `checkElementVisible(page, selector)` - Verify element is visible
- `checkElementContainsText(page, selector, text)` - Verify element contains text

**State Management**:
- `waitForText(page, text)` - Wait for text to appear
- `waitForElementToDisappear(page, selector)` - Wait for element to disappear
- `waitForLoadingState(page, selector)` - Handle loading states

**Debugging**:
- `takeScreenshot(page, name)` - Take screenshots for debugging
- `getElementText(page, selector)` - Get element text content

**Form Handling**:
- `submitForm(page, selector)` - Submit form safely
- `checkButtonEnabled(page, selector)` - Check button state
- `checkButtonDisabled(page, selector)` - Check button state

**Error Handling**:
- `checkErrorDisplayed(page, selector, message)` - Verify error messages
- `checkSuccessMessage(page, selector, message)` - Verify success messages

## Test Execution

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
- Use common test utilities and helper functions
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
- **Screen reader** compatibility

## Debugging and Reporting

### Screenshots and Videos
- **Screenshots**: Automatically captured on test failure
- **Videos**: Recorded for failed tests
- **Traces**: Generated for retried tests

### Test Reports
- **HTML Report**: Interactive test results
- **JUnit Report**: CI/CD integration
- **JSON Report**: Programmatic access

### Debug Mode
- **UI Mode**: Interactive test runner
- **Debug Mode**: Step-by-step debugging
- **Headed Mode**: Visible browser execution

## Error Handling

### Network Error Testing
- **API failure simulation**
- **Network timeout handling**
- **Server error responses**

### Form Validation Testing
- **Input validation** testing
- **Error message** display
- **Form submission** handling

### Loading State Testing
- **Loading indicators** display
- **State transitions** handling
- **Timeout handling**

## Coverage Strategy

### What We Test
1. **User Workflows**: Complete user journeys from start to finish
2. **Critical Paths**: Essential functionality that must work
3. **Error Scenarios**: How the app handles failures gracefully
4. **Responsive Design**: Mobile and tablet compatibility
5. **Accessibility**: Basic accessibility requirements

### What We Don't Test
1. **Implementation Details**: Internal code structure
2. **Every Edge Case**: Focus on likely scenarios
3. **Perfect Coverage**: Aim for 80/20 rule
4. **Overly Specific Selectors**: Use flexible selectors

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

### Test Structure Template
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

## Summary

This organized test structure provides:

- **Coverage** of user workflows and critical functionality
- **Maintainable tests** that adapt to UI changes gracefully
- **Reusable utilities** and helper functions
- **Flexible execution** with multiple test categories
- **Documentation** and examples
- **CI/CD integration** ready for automated testing

The tests focus on **what users do** rather than **how the code works**, making them more valuable for ensuring the application works correctly from a user perspective while being easier to maintain as the UI evolves.
