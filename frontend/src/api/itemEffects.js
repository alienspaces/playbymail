import { baseUrl, getAuthHeaders, apiFetch, handleApiError } from './baseUrl';

/**
 * Fetch all item effects for a game.
 * @param {string} gameId
 * @returns {Promise<GameItemEffect[]>}
 */
export async function fetchItemEffects(gameId, params = {}) {
  const url = new URL(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/item-effects`);
  if (params.page_number) url.searchParams.set('page_number', params.page_number);
  const res = await apiFetch(url.toString(), {
    headers: { ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to fetch item effects');
  const json = await res.json();
  const pagination = JSON.parse(res.headers.get('X-Pagination') || '{}');
  return { data: json.data || [], hasMore: !!pagination.has_more };
}

/**
 * Create a new item effect for a game.
 * @param {string} gameId
 * @param {Partial<GameItemEffect>} data
 * @returns {Promise<GameItemEffect>}
 */
export async function createItemEffect(gameId, data) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/item-effects`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data),
  });
  await handleApiError(res, 'Failed to create item effect');
  const json = await res.json();
  return json.data;
}

/**
 * Update an item effect by ID.
 * @param {string} gameId
 * @param {string} itemEffectId
 * @param {Partial<GameItemEffect>} data
 * @returns {Promise<GameItemEffect>}
 */
export async function updateItemEffect(gameId, itemEffectId, data) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/item-effects/${encodeURIComponent(itemEffectId)}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data),
  });
  await handleApiError(res, 'Failed to update item effect');
  const json = await res.json();
  return json.data;
}

/**
 * Delete an item effect by ID.
 * @param {string} gameId
 * @param {string} itemEffectId
 * @returns {Promise<void>}
 */
export async function deleteItemEffect(gameId, itemEffectId) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/item-effects/${encodeURIComponent(itemEffectId)}`, {
    method: 'DELETE',
    headers: { ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to delete item effect');
}
