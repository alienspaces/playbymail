// MailPit REST API helpers for Playwright E2E tests.
// Assumes MailPit is running on localhost:8025 (started via tools/mailpit-start).

const MAILPIT_API = process.env.MAILPIT_API_URL || 'http://localhost:8025/api/v1'

/**
 * Poll MailPit until an email matching the criteria arrives.
 * @param {string} recipientEmail - Email address to match in the "To" field
 * @param {string} subjectContains - Substring to match in the subject
 * @param {object} [options]
 * @param {number} [options.timeout=30000] - Max wait in ms
 * @param {number} [options.interval=1000] - Poll interval in ms
 * @returns {Promise<object>} The matching MailPit message summary
 */
export async function waitForEmail(recipientEmail, subjectContains, { timeout = 30000, interval = 1000 } = {}) {
  const deadline = Date.now() + timeout
  const query = `to:${recipientEmail} subject:${subjectContains}`

  while (Date.now() < deadline) {
    const res = await fetch(`${MAILPIT_API}/search?query=${encodeURIComponent(query)}`)
    if (res.ok) {
      const data = await res.json()
      if (data.messages && data.messages.length > 0) {
        return data.messages[0]
      }
    }
    await new Promise(r => setTimeout(r, interval))
  }

  throw new Error(`Timed out waiting for email to=${recipientEmail} subject~="${subjectContains}" after ${timeout}ms`)
}

/**
 * Get the full message (including HTML body) by message ID.
 * @param {string} messageId
 * @returns {Promise<object>} Full MailPit message object
 */
export async function getEmailBody(messageId) {
  const res = await fetch(`${MAILPIT_API}/message/${messageId}`)
  if (!res.ok) {
    throw new Error(`Failed to get message ${messageId}: ${res.status}`)
  }
  return res.json()
}

/**
 * Extract the first href matching a pattern from an HTML body.
 * @param {string} htmlBody
 * @param {RegExp|string} pattern - Pattern to match inside href values
 * @returns {string|null} The matched URL or null
 */
export function extractLink(htmlBody, pattern) {
  const re = pattern instanceof RegExp ? pattern : new RegExp(pattern)
  const hrefRegex = /href=["']([^"']+)["']/gi
  let match
  while ((match = hrefRegex.exec(htmlBody)) !== null) {
    if (re.test(match[1])) {
      return match[1]
    }
  }
  return null
}

/**
 * Delete all messages in MailPit.
 */
export async function clearAllEmails() {
  await fetch(`${MAILPIT_API}/messages`, { method: 'DELETE' })
}

/**
 * Get the latest email for a recipient.
 * @param {string} recipientEmail
 * @returns {Promise<object|null>} Latest message summary or null
 */
export async function getLatestEmail(recipientEmail) {
  const query = `to:${recipientEmail}`
  const res = await fetch(`${MAILPIT_API}/search?query=${encodeURIComponent(query)}`)
  if (!res.ok) return null
  const data = await res.json()
  if (data.messages && data.messages.length > 0) {
    return data.messages[0]
  }
  return null
}

/**
 * Convenience: wait for an email, fetch its body, and extract a link.
 * @param {string} recipientEmail
 * @param {string} subjectContains
 * @param {RegExp|string} linkPattern
 * @param {object} [options] - Passed to waitForEmail
 * @returns {Promise<string>} The extracted URL
 */
export async function waitForEmailLink(recipientEmail, subjectContains, linkPattern, options) {
  const summary = await waitForEmail(recipientEmail, subjectContains, options)
  const full = await getEmailBody(summary.ID)
  const html = full.HTML || full.Text || ''
  const link = extractLink(html, linkPattern)
  if (!link) {
    throw new Error(`No link matching ${linkPattern} found in email "${summary.Subject}"`)
  }
  return link
}
