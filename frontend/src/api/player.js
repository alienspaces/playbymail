import { baseUrl, getAuthHeaders, handleApiError } from './baseUrl';

/**
 * Verify game subscription turn sheet token
 * @param {string} gameSubscriptionID - Game subscription ID
 * @param {string} gameInstanceID - Game instance ID
 * @param {string} email - Account email address
 * @param {string} turnSheetToken - Turn sheet token from email link
 * @returns {Promise<string>} Session token
 */
export async function verifyGameSubscriptionToken(gameSubscriptionID, gameInstanceID, email, turnSheetToken) {
  const res = await fetch(
    `${baseUrl}/api/v1/player/game-subscriptions/${gameSubscriptionID}/game-instances/${gameInstanceID}/verify-token`,
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


