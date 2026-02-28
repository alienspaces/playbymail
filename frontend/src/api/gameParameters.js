import { baseUrl, getAuthHeaders, apiFetch, handleApiError } from './baseUrl';

export async function listGameParameters(params = {}) {
  const queryParams = new URLSearchParams();
  if (params.gameType) queryParams.append('game_type', params.gameType);
  if (params.configKey) queryParams.append('config_key', params.configKey);
  
  const url = `${baseUrl}/api/v1/game-parameters${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;
  const res = await apiFetch(url, {
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to fetch game parameters');
  return await res.json();
}
