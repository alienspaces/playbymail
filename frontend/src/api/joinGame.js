import { baseUrl, apiFetch, getAuthHeaders, handleApiError } from './baseUrl';

const joinPath = (gameSubscriptionId) =>
  `${baseUrl}/api/v1/game-subscriptions/${gameSubscriptionId}/join`;

export async function getJoinGameInfo(gameSubscriptionId) {
  const res = await apiFetch(joinPath(gameSubscriptionId), {
    headers: { 'Content-Type': 'application/json' },
  });
  await handleApiError(res, 'Failed to load game information');
  return await res.json();
}

export async function getJoinSheet(gameSubscriptionId) {
  const res = await apiFetch(`${joinPath(gameSubscriptionId)}/sheet`, {
    headers: { ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to load join game turn sheet');
  return await res.text();
}

export async function submitJoinGame(gameSubscriptionId, data) {
  const res = await apiFetch(joinPath(gameSubscriptionId), {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data),
  });
  await handleApiError(res, 'Failed to submit join game');
  return await res.json();
}
