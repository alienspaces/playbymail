import { baseUrl, handleApiError } from './baseUrl';

const joinPath = (gameSubscriptionId) =>
  `${baseUrl}/api/v1/game-subscriptions/${gameSubscriptionId}/join`;

export async function getJoinGameInfo(gameSubscriptionId) {
  const res = await fetch(joinPath(gameSubscriptionId), {
    headers: { 'Content-Type': 'application/json' },
  });
  await handleApiError(res, 'Failed to load game information');
  return await res.json();
}

export async function getJoinSheet(gameSubscriptionId) {
  const res = await fetch(`${joinPath(gameSubscriptionId)}/sheet`);
  await handleApiError(res, 'Failed to load join game turn sheet');
  return await res.text();
}

export async function submitJoinGame(gameSubscriptionId, data) {
  const res = await fetch(joinPath(gameSubscriptionId), {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  await handleApiError(res, 'Failed to submit join game');
  return await res.json();
}
