# Frontend

This template should help get you started developing with Vue 3 in Vite.

## Recommended IDE Setup

[VSCode](https://code.visualstudio.com/) + [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar) (and disable Vetur).

## Customize configuration

See [Vite Configuration Reference](https://vite.dev/config/).

## Project Setup

```sh
npm install
```

### Compile and Hot-Reload for Development

```sh
npm run dev
```

### Compile and Minify for Production

```sh
npm run build
```

### Run Unit Tests with [Vitest](https://vitest.dev/)

```sh
npm run test:unit
```

### Lint with [ESLint](https://eslint.org/)

```sh
npm run lint
```

## Development Authentication

When running the backend with `./tools/start-backend`, the frontend automatically supports simplified authentication for local development.

### How It Works

When `APP_ENV=develop` (set in `.env.develop`), the frontend includes test bypass headers in authentication requests. This allows you to **use the email address as the verification code** instead of needing to receive an actual email.

### Login Flow in Development

1. Enter one of the test email addresses on the login page
2. Click "Send verification code"
3. On the verification page, enter the **same email address** as the verification code
4. You'll be logged in with a valid session

### Test User Accounts

The following test accounts are seeded when you run `./tools/start-backend`:

| Email | Description |
|-------|-------------|
| `test-account-one@example.com` | Test Account One |
| `test-account-two@example.com` | Test Account Two |
| `test-account-three@example.com` | Test Account Three |

### Configuration

The bypass is controlled by these environment variables (set in `.env.develop`):

- `APP_ENV=develop` - Enables development mode
- `TEST_BYPASS_HEADER_NAME` - The HTTP header name for bypass authentication
- `TEST_BYPASS_HEADER_VALUE` - The required header value to enable bypass

### Production Behavior

In production (`APP_ENV=production`), the bypass headers are not included and real email verification is required. The environment variables are empty or not set, so no bypass functionality is available.
