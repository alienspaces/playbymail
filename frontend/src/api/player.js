import { baseUrl, getAuthHeaders, handleApiError } from './baseUrl';

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



