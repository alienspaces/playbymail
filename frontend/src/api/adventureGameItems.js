import { baseUrl, getAuthHeaders, apiFetch, handleApiError } from './baseUrl';

/**
 * Fetch all items for a game.
 * @param {string} gameId
 * @returns {Promise<GameItem[]>}
 */
export async function fetchAdventureGameItems(gameId, params = {}) {
  const url = new URL(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/items`);
  if (params.page_number) url.searchParams.set('page_number', params.page_number);
  const res = await apiFetch(url.toString(), {
    headers: { ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to fetch items');
  const json = await res.json();
  const pagination = JSON.parse(res.headers.get('X-Pagination') || '{}');
  return { data: json.data || [], hasMore: !!pagination.has_more };
}

/**
 * Create a new item for a game.
 * @param {string} gameId
 * @param {Partial<GameItem>} data
 * @returns {Promise<GameItem>}
 */
export async function createAdventureGameItem(gameId, data) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/items`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data)
  });
  await handleApiError(res, 'Failed to create item');
  const json = await res.json();
  return json.data;
}

/**
 * Update an item by ID.
 * @param {string} gameId
 * @param {string} itemId
 * @param {Partial<GameItem>} data
 * @returns {Promise<GameItem>}
 */
export async function updateAdventureGameItem(gameId, itemId, data) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/items/${encodeURIComponent(itemId)}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data)
  });
  await handleApiError(res, 'Failed to update item');
  const json = await res.json();
  return json.data;
}

/**
 * Delete an item by ID.
 * @param {string} gameId
 * @param {string} itemId
 * @returns {Promise<void>}
 */
export async function deleteAdventureGameItem(gameId, itemId) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/items/${encodeURIComponent(itemId)}`, {
    method: 'DELETE',
    headers: { ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to delete item');
} 