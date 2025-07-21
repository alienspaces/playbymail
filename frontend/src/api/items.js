import { baseUrl, getAuthHeaders, apiFetch } from './baseUrl';

/**
 * Fetch all items for a game.
 * @param {string} gameId
 * @returns {Promise<GameItem[]>}
 */
export async function fetchItems(gameId) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/items`, {
    headers: { ...getAuthHeaders() },
  });
  if (!res.ok) throw new Error('Failed to fetch items');
  const json = await res.json();
  return json.data || [];
}

/**
 * Create a new item for a game.
 * @param {string} gameId
 * @param {Partial<GameItem>} data
 * @returns {Promise<GameItem>}
 */
export async function createItem(gameId, data) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/items`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data)
  });
  if (!res.ok) throw new Error('Failed to create item');
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
export async function updateItem(gameId, itemId, data) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/items/${encodeURIComponent(itemId)}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data)
  });
  if (!res.ok) throw new Error('Failed to update item');
  const json = await res.json();
  return json.data;
}

/**
 * Delete an item by ID.
 * @param {string} gameId
 * @param {string} itemId
 * @returns {Promise<void>}
 */
export async function deleteItem(gameId, itemId) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/items/${encodeURIComponent(itemId)}`, {
    method: 'DELETE',
    headers: { ...getAuthHeaders() },
  });
  if (!res.ok) throw new Error('Failed to delete item');
} 