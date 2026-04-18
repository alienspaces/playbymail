import { baseUrl, getAuthHeaders, apiFetch, handleApiError } from './baseUrl';

/**
 * Fetch all states for a location object.
 * @param {string} gameId
 * @param {string} locationObjectId
 * @returns {Promise<GameLocationObjectState[]>}
 */
export async function fetchAdventureGameLocationObjectStates(gameId, locationObjectId, params = {}) {
  const url = new URL(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/location-objects/${encodeURIComponent(locationObjectId)}/states`);
  if (params.page_number) url.searchParams.set('page_number', params.page_number);
  const res = await apiFetch(url.toString(), { headers: { ...getAuthHeaders() } });
  await handleApiError(res, 'Failed to fetch location object states');
  const json = await res.json();
  const pagination = JSON.parse(res.headers.get('X-Pagination') || '{}');
  return { data: json.data || [], hasMore: !!pagination.has_more };
}

/**
 * Create a new state for a location object.
 * @param {string} gameId
 * @param {string} locationObjectId
 * @param {Partial<GameLocationObjectState>} data
 * @returns {Promise<GameLocationObjectState>}
 */
export async function createAdventureGameLocationObjectState(gameId, locationObjectId, data) {
  const res = await apiFetch(
    `${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/location-objects/${encodeURIComponent(locationObjectId)}/states`,
    {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
      body: JSON.stringify(data),
    }
  );
  await handleApiError(res, 'Failed to create location object state');
  const json = await res.json();
  return json.data;
}

/**
 * Update a location object state by ID.
 * @param {string} gameId
 * @param {string} locationObjectId
 * @param {string} stateId
 * @param {Partial<GameLocationObjectState>} data
 * @returns {Promise<GameLocationObjectState>}
 */
export async function updateAdventureGameLocationObjectState(gameId, locationObjectId, stateId, data) {
  const res = await apiFetch(
    `${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/location-objects/${encodeURIComponent(locationObjectId)}/states/${encodeURIComponent(stateId)}`,
    {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
      body: JSON.stringify(data),
    }
  );
  await handleApiError(res, 'Failed to update location object state');
  const json = await res.json();
  return json.data;
}

/**
 * Delete a location object state by ID.
 * @param {string} gameId
 * @param {string} locationObjectId
 * @param {string} stateId
 * @returns {Promise<void>}
 */
export async function deleteAdventureGameLocationObjectState(gameId, locationObjectId, stateId) {
  const res = await apiFetch(
    `${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/location-objects/${encodeURIComponent(locationObjectId)}/states/${encodeURIComponent(stateId)}`,
    {
      method: 'DELETE',
      headers: { ...getAuthHeaders() },
    }
  );
  await handleApiError(res, 'Failed to delete location object state');
}
