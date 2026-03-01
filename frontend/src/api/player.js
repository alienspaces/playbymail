import { baseUrl, apiFetch, getAuthHeaders, handleApiError } from './baseUrl';

const gsiPath = (gsiId) =>
  `${baseUrl}/api/v1/player/game-subscription-instances/${gsiId}`;

/**
 * Verify game subscription instance turn sheet token
 * @param {string} gameSubscriptionInstanceID - Game subscription instance ID
 * @param {string} email - Account email address
 * @param {string} turnSheetToken - Turn sheet token from email link
 * @returns {Promise<string>} Session token
 */
export async function verifyGameSubscriptionToken(gameSubscriptionInstanceID, email, turnSheetToken) {
  const res = await fetch(
    `${baseUrl}/api/v1/player/game-subscription-instances/${gameSubscriptionInstanceID}/verify-token`,
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...getAuthHeaders(),
      },
      body: JSON.stringify({ email, turn_sheet_token: turnSheetToken }),
    }
  );
  await handleApiError(res, 'Token verification failed');
  const data = await res.json();
  return data.session_token;
}

/**
 * Get the list of turn sheets for a game subscription instance.
 * @param {string} gsiId
 * @returns {Promise<object>}
 */
export async function getGSITurnSheets(gsiId) {
  const res = await apiFetch(`${gsiPath(gsiId)}/turn-sheets`, {
    headers: { 'Content-Type': 'application/json' },
  });
  await handleApiError(res, 'Failed to load turn sheets');
  return await res.json();
}

/**
 * Get a specific turn sheet for a game subscription instance.
 * @param {string} gsiId
 * @param {string} turnSheetId
 * @returns {Promise<object>}
 */
export async function getGSITurnSheet(gsiId, turnSheetId) {
  const res = await apiFetch(`${gsiPath(gsiId)}/turn-sheets/${turnSheetId}`, {
    headers: { 'Content-Type': 'application/json' },
  });
  await handleApiError(res, 'Failed to load turn sheet');
  return await res.json();
}

/**
 * Save (auto-save) form data for a turn sheet.
 * @param {string} gsiId
 * @param {string} turnSheetId
 * @param {object} scannedData
 * @returns {Promise<object>}
 */
export async function saveGSITurnSheet(gsiId, turnSheetId, scannedData) {
  const res = await apiFetch(`${gsiPath(gsiId)}/turn-sheets/${turnSheetId}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ scanned_data: scannedData }),
  });
  await handleApiError(res, 'Failed to save turn sheet');
  return await res.json();
}

/**
 * Submit all turn sheets for a game subscription instance.
 * @param {string} gsiId
 * @returns {Promise<object>}
 */
export async function submitGSITurnSheets(gsiId) {
  const res = await apiFetch(`${gsiPath(gsiId)}/turn-sheets/submit`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
  });
  await handleApiError(res, 'Failed to submit turn sheets');
  return await res.json();
}

