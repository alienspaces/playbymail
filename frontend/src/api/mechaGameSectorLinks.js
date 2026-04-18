import { baseUrl, getAuthHeaders, apiFetch, handleApiError } from './baseUrl';

const BASE = (gameId) => `${baseUrl}/api/v1/mecha-games/${encodeURIComponent(gameId)}/sector-links`

export async function fetchMechaGameSectorLinks(gameId, params = {}) {
  const url = new URL(BASE(gameId))
  if (params.page_number) url.searchParams.set('page_number', params.page_number)
  const res = await apiFetch(url.toString(), { headers: { ...getAuthHeaders() } })
  await handleApiError(res, 'Failed to fetch sector links')
  const json = await res.json()
  const pagination = JSON.parse(res.headers.get('X-Pagination') || '{}')
  return { data: json.data || [], hasMore: !!pagination.has_more }
}

export async function createMechaGameSectorLink(gameId, data) {
  const res = await apiFetch(BASE(gameId), {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data),
  })
  await handleApiError(res, 'Failed to create sector link')
  const json = await res.json()
  return json.data
}

export async function updateMechaGameSectorLink(gameId, sectorLinkId, data) {
  const res = await apiFetch(`${BASE(gameId)}/${encodeURIComponent(sectorLinkId)}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data),
  })
  await handleApiError(res, 'Failed to update sector link')
  const json = await res.json()
  return json.data
}

export async function deleteMechaGameSectorLink(gameId, sectorLinkId) {
  const res = await apiFetch(`${BASE(gameId)}/${encodeURIComponent(sectorLinkId)}`, {
    method: 'DELETE',
    headers: { ...getAuthHeaders() },
  })
  await handleApiError(res, 'Failed to delete sector link')
}
