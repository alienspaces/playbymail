import { baseUrl, getAuthHeaders, apiFetch, handleApiError } from './baseUrl';

/**
 * Fetch all creature placements for a game.
 * @param {string} gameId
 * @returns {Promise<CreaturePlacement[]>}
 */
export async function fetchAdventureGameCreaturePlacements(gameId, params = {}) {
  const url = new URL(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/creature-placements`);
  if (params.page_number) url.searchParams.set('page_number', params.page_number);
  const res = await apiFetch(url.toString(), {
    headers: { ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to fetch creature placements');
  const json = await res.json();
  const pagination = JSON.parse(res.headers.get('X-Pagination') || '{}');
  return { data: json.data || [], hasMore: !!pagination.has_more };
}

/**
 * Create a new creature placement for a game.
 * @param {string} gameId
 * @param {Partial<CreaturePlacement>} data
 * @returns {Promise<CreaturePlacement>}
 */
export async function createAdventureGameCreaturePlacement(gameId, data) {
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
export async function updateAdventureGameCreaturePlacement(gameId, placementId, data) {
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
export async function deleteAdventureGameCreaturePlacement(gameId, placementId) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/creature-placements/${encodeURIComponent(placementId)}`, {
    method: 'DELETE',
    headers: { ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to delete creature placement');
} 