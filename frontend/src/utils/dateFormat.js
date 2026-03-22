/**
 * Shared date/time formatting utilities.
 *
 * Timezone resolution order:
 *   1. Caller-supplied timezone (from account profile setting)
 *   2. Browser timezone (Intl.DateTimeFormat().resolvedOptions().timeZone)
 *   3. Fallback: 'UTC'
 *
 * Locale resolution order:
 *   1. Caller-supplied locale
 *   2. navigator.language
 *   3. Fallback: 'en-US'
 *
 * All date strings from the API are UTC ISO 8601; timezone conversion happens
 * only at display time here.
 *
 * Do NOT use toLocaleDateString(), toLocaleString(), or toLocaleTimeString()
 * directly in components. Always import from this module instead.
 */

function resolveTimezone(timezone) {
  return timezone || Intl.DateTimeFormat().resolvedOptions().timeZone || 'UTC'
}

function resolveLocale(locale) {
  return locale || navigator.language || 'en-US'
}

/**
 * Format a date+time string for display, including hours and minutes.
 * Returns 'N/A' for null/undefined/empty input.
 */
export function formatDateTime(dateString, { timezone, locale } = {}) {
  if (!dateString) return 'N/A'
  return new Date(dateString).toLocaleString(resolveLocale(locale), {
    timeZone: resolveTimezone(timezone),
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

/**
 * Format a date string for display, date only (no time).
 * Returns 'N/A' for null/undefined/empty input.
 */
export function formatDateOnly(dateString, { timezone, locale } = {}) {
  if (!dateString) return 'N/A'
  return new Date(dateString).toLocaleDateString(resolveLocale(locale), {
    timeZone: resolveTimezone(timezone),
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}

/**
 * Format a deadline date/time with contextual labels for near-future times.
 *
 * - Past deadlines:           'Overdue'
 * - Due within 24 hours:      'Today, HH:MM'
 * - Due within 48 hours:      'Tomorrow, HH:MM'
 * - Otherwise:                full date+time via formatDateTime
 *
 * Returns 'N/A' for null/undefined/empty input.
 */
export function formatDeadline(deadlineString, { timezone, locale } = {}) {
  if (!deadlineString) return 'N/A'
  const deadline = new Date(deadlineString)
  const now = new Date()
  const diff = deadline - now

  if (diff < 0) return 'Overdue'

  const timeStr = deadline.toLocaleTimeString(resolveLocale(locale), {
    timeZone: resolveTimezone(timezone),
    hour: '2-digit',
    minute: '2-digit',
  })

  if (diff < 24 * 60 * 60 * 1000) return `Today, ${timeStr}`
  if (diff < 48 * 60 * 60 * 1000) return `Tomorrow, ${timeStr}`

  return formatDateTime(deadlineString, { timezone, locale })
}
