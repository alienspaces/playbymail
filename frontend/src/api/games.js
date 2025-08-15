import { baseUrl, getAuthHeaders, apiFetch } from './baseUrl';

export async function listGames() {
  const res = await apiFetch(`${baseUrl}/api/v1/games`, {
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  if (!res.ok) throw new Error('Failed to fetch games');
  return await res.json();
}

export async function createGame({ name, game_type, turn_duration_hours }) {
  const res = await apiFetch(`${baseUrl}/api/v1/games`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify({ name, game_type, turn_duration_hours }),
  });
  if (!res.ok) throw new Error('Failed to create game');
  return await res.json();
}

export async function updateGame(id, { name, game_type, turn_duration_hours }) {
  const res = await apiFetch(`${baseUrl}/api/v1/games/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify({ name, game_type, turn_duration_hours }),
  });
  if (!res.ok) throw new Error('Failed to update game');
  return await res.json();
}

export async function deleteGame(id) {
  const res = await apiFetch(`${baseUrl}/api/v1/games/${id}`, {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  if (!res.ok) throw new Error('Failed to delete game');
  return await res.json();
} 