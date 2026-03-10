import { baseUrl, getAuthHeaders, apiFetch, handleApiError } from './baseUrl';

const gameSubscriptionInstancePath = (gameSubscriptionInstanceId) =>
  `${baseUrl}/api/v1/player/game-subscription-instances/${gameSubscriptionInstanceId}`;

/**
 * Verify game subscription instance turn sheet token.
 * @param {string} gameSubscriptionInstanceID
 * @param {string} turnSheetToken
 * @returns {Promise<string>} Session token
 */
export async function verifyGameSubscriptionToken(gameSubscriptionInstanceID, turnSheetToken) {
  const res = await fetch(
    `${gameSubscriptionInstancePath(gameSubscriptionInstanceID)}/verify-token`,
    {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ turn_sheet_token: turnSheetToken }),
    }
  );
  await handleApiError(res, 'Token verification failed');
  const data = await res.json();
  return data.session_token;
}

/**
 * Request a new turn sheet token (e.g. when the current one has expired).
 * The backend sends a fresh email with a new link.
 * @param {string} gameSubscriptionInstanceID
 * @param {string} email
 */
export async function requestNewTurnSheetToken(gameSubscriptionInstanceID, email) {
  const res = await fetch(
    `${gameSubscriptionInstancePath(gameSubscriptionInstanceID)}/request-token`,
    {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email }),
    }
  );
  await handleApiError(res, 'Failed to request new link');
}

/**
 * Get the list of turn sheets for a game subscription instance.
 * @param {string} gameSubscriptionInstanceId
 * @returns {Promise<object>}
 */
export async function getGameSubscriptionInstanceTurnSheets(gameSubscriptionInstanceId) {
  const res = await apiFetch(`${gameSubscriptionInstancePath(gameSubscriptionInstanceId)}/turn-sheets`, {
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to load turn sheets');
  return await res.json();
}

/**
 * Get a specific turn sheet for a game subscription instance.
 * @param {string} gameSubscriptionInstanceId
 * @param {string} turnSheetId
 * @returns {Promise<object>}
 */
export async function getGameSubscriptionInstanceTurnSheet(gameSubscriptionInstanceId, turnSheetId) {
  const res = await apiFetch(`${gameSubscriptionInstancePath(gameSubscriptionInstanceId)}/turn-sheets/${turnSheetId}`, {
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to load turn sheet');
  return await res.json();
}

/**
 * Save (auto-save) form data for a turn sheet.
 * @param {string} gameSubscriptionInstanceId
 * @param {string} turnSheetId
 * @param {object} scannedData
 * @returns {Promise<object>}
 */
export async function saveGameSubscriptionInstanceTurnSheet(gameSubscriptionInstanceId, turnSheetId, scannedData) {
  const res = await apiFetch(`${gameSubscriptionInstancePath(gameSubscriptionInstanceId)}/turn-sheets/${turnSheetId}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify({ scanned_data: scannedData }),
  });
  await handleApiError(res, 'Failed to save turn sheet');
  return await res.json();
}

/**
 * Submit all turn sheets for a game subscription instance.
 * @param {string} gameSubscriptionInstanceId
 * @returns {Promise<object>}
 */
export async function submitGameSubscriptionInstanceTurnSheets(gameSubscriptionInstanceId) {
  const res = await apiFetch(`${gameSubscriptionInstancePath(gameSubscriptionInstanceId)}/turn-sheet-upload`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to submit turn sheets');
  return await res.json();
}

/**
 * Get a turn sheet rendered as HTML (for inline viewer).
 * @param {string} gameSubscriptionInstanceId
 * @param {string} turnSheetId
 * @returns {Promise<string>} HTML string
 */
export async function getGameSubscriptionInstanceTurnSheetHTML(gameSubscriptionInstanceId, turnSheetId) {
  const res = await apiFetch(`${gameSubscriptionInstancePath(gameSubscriptionInstanceId)}/turn-sheets/${turnSheetId}`, {
    headers: { Accept: 'text/html', ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to load turn sheet HTML');
  return await res.text();
}

/**
 * Download a printable PDF for a turn sheet.
 * Returns the raw Response so the caller can trigger a file download.
 * @param {string} gameSubscriptionInstanceId
 * @param {string} turnSheetId
 * @returns {Promise<Response>}
 */
export async function downloadGameSubscriptionInstanceTurnSheetPDF(gameSubscriptionInstanceId, turnSheetId) {
  const res = await apiFetch(`${gameSubscriptionInstancePath(gameSubscriptionInstanceId)}/turn-sheets/${turnSheetId}/download`, {
    headers: { Accept: 'application/pdf', ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to download turn sheet PDF');
  return res;
}

/**
 * Upload a scanned turn sheet image for OCR processing.
 * @param {string} gameSubscriptionInstanceId
 * @param {string} turnSheetId
 * @param {File} imageFile
 * @returns {Promise<object>}
 */
export async function uploadGameSubscriptionInstanceTurnSheetScan(gameSubscriptionInstanceId, turnSheetId, imageFile) {
  const formData = new FormData()
  formData.append('image', imageFile)

  const res = await apiFetch(`${gameSubscriptionInstancePath(gameSubscriptionInstanceId)}/turn-sheets/${turnSheetId}/scan`, {
    method: 'POST',
    headers: { ...getAuthHeaders() },
    body: formData,
  });
  await handleApiError(res, 'Failed to upload scanned turn sheet');
  return await res.json();
}

