import { baseUrl, apiFetch, handleApiError } from './baseUrl';

const joinGamePath = (gameInstanceId) =>
  `${baseUrl}/api/v1/player/game-instances/${gameInstanceId}/join-game`;

export async function getJoinGameInfo(gameInstanceId) {
  const res = await apiFetch(joinGamePath(gameInstanceId), {
    headers: { 'Content-Type': 'application/json' },
  });
  await handleApiError(res, 'Failed to load game information');
  return await res.json();
}

export async function verifyJoinGameEmail(gameInstanceId, email) {
  const res = await apiFetch(`${joinGamePath(gameInstanceId)}/verify-email`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email }),
  });
  await handleApiError(res, 'Failed to verify email');
  return await res.json();
}

export async function submitJoinGame(gameInstanceId, data) {
  const res = await apiFetch(joinGamePath(gameInstanceId), {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  await handleApiError(res, 'Failed to submit join game');
  return await res.json();
}
