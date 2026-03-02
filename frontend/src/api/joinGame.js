import { baseUrl, apiFetch, handleApiError } from './baseUrl';

const joinPath = (gameSubscriptionId) =>
  `${baseUrl}/api/v1/game-subscriptions/${gameSubscriptionId}/join`;

export async function getJoinGameInfo(gameSubscriptionId) {
  const res = await apiFetch(joinPath(gameSubscriptionId), {
    headers: { 'Content-Type': 'application/json' },
  });
  await handleApiError(res, 'Failed to load game information');
  return await res.json();
}

export async function verifyJoinGameEmail(gameSubscriptionId, email) {
  const res = await apiFetch(`${joinPath(gameSubscriptionId)}/verify-email`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email }),
  });
  await handleApiError(res, 'Failed to verify email');
  return await res.json();
}

export async function submitJoinGame(gameSubscriptionId, data) {
  const res = await apiFetch(joinPath(gameSubscriptionId), {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  await handleApiError(res, 'Failed to submit join game');
  return await res.json();
}
