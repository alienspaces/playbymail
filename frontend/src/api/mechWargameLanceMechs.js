import { baseUrl, getAuthHeaders, apiFetch, handleApiError } from './baseUrl';

const BASE = (gameId, lanceId) =>
  `${baseUrl}/api/v1/mech-wargame-games/${encodeURIComponent(gameId)}/lances/${encodeURIComponent(lanceId)}/mechs`

export async function fetchLanceMechs(gameId, lanceId, params = {}) {
  const url = new URL(BASE(gameId, lanceId))
  if (params.page_number) url.searchParams.set('page_number', params.page_number)
  const res = await apiFetch(url.toString(), { headers: { ...getAuthHeaders() } })
  await handleApiError(res, 'Failed to fetch lance mechs')
  const json = await res.json()
  const pagination = JSON.parse(res.headers.get('X-Pagination') || '{}')
  return { data: json.data || [], hasMore: !!pagination.has_more }
}

export async function createLanceMech(gameId, lanceId, data) {
  const res = await apiFetch(BASE(gameId, lanceId), {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data),
  })
  await handleApiError(res, 'Failed to create lance mech')
  const json = await res.json()
  return json.data
}

export async function updateLanceMech(gameId, lanceId, mechId, data) {
  const res = await apiFetch(`${BASE(gameId, lanceId)}/${encodeURIComponent(mechId)}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data),
  })
  await handleApiError(res, 'Failed to update lance mech')
  const json = await res.json()
  return json.data
}

export async function deleteLanceMech(gameId, lanceId, mechId) {
  const res = await apiFetch(`${BASE(gameId, lanceId)}/${encodeURIComponent(mechId)}`, {
    method: 'DELETE',
    headers: { ...getAuthHeaders() },
  })
  await handleApiError(res, 'Failed to delete lance mech')
}
