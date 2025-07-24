import { baseUrl, getAuthHeaders, apiFetch } from './baseUrl';

/**
 * Fetch all item placements for a game.
 * @param {string} gameId
 * @returns {Promise<ItemPlacement[]>}
 */
export async function fetchItemPlacements(gameId) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/item-placements`, {
    headers: { ...getAuthHeaders() },
  });
  if (!res.ok) throw new Error('Failed to fetch item placements');
  const json = await res.json();
  return json.data || [];
}

/**
 * Create a new item placement for a game.
 * @param {string} gameId
 * @param {Partial<ItemPlacement>} data
 * @returns {Promise<ItemPlacement>}
 */
export async function createItemPlacement(gameId, data) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/item-placements`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data)
  });
  if (!res.ok) throw new Error('Failed to create item placement');
  const json = await res.json();
  return json.data;
}

/**
 * Update an item placement by ID.
 * @param {string} gameId
 * @param {string} placementId
 * @param {Partial<ItemPlacement>} data
 * @returns {Promise<ItemPlacement>}
 */
export async function updateItemPlacement(gameId, placementId, data) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/item-placements/${encodeURIComponent(placementId)}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data)
  });
  if (!res.ok) throw new Error('Failed to update item placement');
  const json = await res.json();
  return json.data;
}

/**
 * Delete an item placement by ID.
 * @param {string} gameId
 * @param {string} placementId
 * @returns {Promise<void>}
 */
export async function deleteItemPlacement(gameId, placementId) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/item-placements/${encodeURIComponent(placementId)}`, {
    method: 'DELETE',
    headers: { ...getAuthHeaders() },
  });
  if (!res.ok) throw new Error('Failed to delete item placement');
} 