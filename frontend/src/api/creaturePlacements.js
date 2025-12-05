import { baseUrl, getAuthHeaders, apiFetch, handleApiError } from './baseUrl';

/**
 * Fetch all creature placements for a game.
 * @param {string} gameId
 * @returns {Promise<CreaturePlacement[]>}
 */
export async function fetchCreaturePlacements(gameId) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/creature-placements`, {
    headers: { ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to fetch creature placements');
  const json = await res.json();
  return json.data || [];
}

/**
 * Create a new creature placement for a game.
 * @param {string} gameId
 * @param {Partial<CreaturePlacement>} data
 * @returns {Promise<CreaturePlacement>}
 */
export async function createCreaturePlacement(gameId, data) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/creature-placements`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data)
  });
  await handleApiError(res, 'Failed to create creature placement');
  const json = await res.json();
  return json.data;
}

/**
 * Update a creature placement by ID.
 * @param {string} gameId
 * @param {string} placementId
 * @param {Partial<CreaturePlacement>} data
 * @returns {Promise<CreaturePlacement>}
 */
export async function updateCreaturePlacement(gameId, placementId, data) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/creature-placements/${encodeURIComponent(placementId)}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data)
  });
  await handleApiError(res, 'Failed to update creature placement');
  const json = await res.json();
  return json.data;
}

/**
 * Delete a creature placement by ID.
 * @param {string} gameId
 * @param {string} placementId
 * @returns {Promise<void>}
 */
export async function deleteCreaturePlacement(gameId, placementId) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/creature-placements/${encodeURIComponent(placementId)}`, {
    method: 'DELETE',
    headers: { ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to delete creature placement');
} 