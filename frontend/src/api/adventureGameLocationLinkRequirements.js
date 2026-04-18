import { baseUrl, getAuthHeaders, apiFetch, handleApiError } from './baseUrl';

/**
 * Fetch all location link requirements for a game.
 * @param {string} gameId
 * @returns {Promise<GameLocationLinkRequirement[]>}
 */
export async function fetchAdventureGameLocationLinkRequirements(gameId, params = {}) {
  const url = new URL(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/location-link-requirements`);
  if (params.page_number) url.searchParams.set('page_number', params.page_number);
  const res = await apiFetch(url.toString(), {
    headers: { ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to fetch location link requirements');
  const json = await res.json();
  const pagination = JSON.parse(res.headers.get('X-Pagination') || '{}');
  return { data: json.data || [], hasMore: !!pagination.has_more };
}

/**
 * Create a new location link requirement for a game.
 * @param {string} gameId
 * @param {Partial<GameLocationLinkRequirement>} data
 * @returns {Promise<GameLocationLinkRequirement>}
 */
export async function createAdventureGameLocationLinkRequirement(gameId, data) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/location-link-requirements`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data)
  });
  await handleApiError(res, 'Failed to create location link requirement');
  const json = await res.json();
  return json.data;
}

/**
 * Update a location link requirement by ID.
 * @param {string} gameId
 * @param {string} requirementId
 * @param {Partial<GameLocationLinkRequirement>} data
 * @returns {Promise<GameLocationLinkRequirement>}
 */
export async function updateAdventureGameLocationLinkRequirement(gameId, requirementId, data) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/location-link-requirements/${encodeURIComponent(requirementId)}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data)
  });
  await handleApiError(res, 'Failed to update location link requirement');
  const json = await res.json();
  return json.data;
}

/**
 * Delete a location link requirement by ID.
 * @param {string} gameId
 * @param {string} requirementId
 * @returns {Promise<void>}
 */
export async function deleteAdventureGameLocationLinkRequirement(gameId, requirementId) {
  const res = await apiFetch(`${baseUrl}/api/v1/adventure-games/${encodeURIComponent(gameId)}/location-link-requirements/${encodeURIComponent(requirementId)}`, {
    method: 'DELETE',
    headers: { ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to delete location link requirement');
}
