import { baseUrl, getAuthHeaders, apiFetch, handleApiError } from './baseUrl';

const BASE = (gameId, squadId) =>
  `${baseUrl}/api/v1/mecha-games/${encodeURIComponent(gameId)}/squads/${encodeURIComponent(squadId)}/mechs`

export async function fetchMechaGameSquadMechs(gameId, squadId, params = {}) {
  const url = new URL(BASE(gameId, squadId))
  if (params.page_number) url.searchParams.set('page_number', params.page_number)
  const res = await apiFetch(url.toString(), { headers: { ...getAuthHeaders() } })
  await handleApiError(res, 'Failed to fetch squad mechs')
  const json = await res.json()
  const pagination = JSON.parse(res.headers.get('X-Pagination') || '{}')
  return { data: json.data || [], hasMore: !!pagination.has_more }
}

export async function createMechaGameSquadMech(gameId, squadId, data) {
  const res = await apiFetch(BASE(gameId, squadId), {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data),
  })
  await handleApiError(res, 'Failed to create squad mech')
  const json = await res.json()
  return json.data
}

export async function updateMechaGameSquadMech(gameId, squadId, mechId, data) {
  const res = await apiFetch(`${BASE(gameId, squadId)}/${encodeURIComponent(mechId)}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data),
  })
  await handleApiError(res, 'Failed to update squad mech')
  const json = await res.json()
  return json.data
}

export async function deleteMechaGameSquadMech(gameId, squadId, mechId) {
  const res = await apiFetch(`${BASE(gameId, squadId)}/${encodeURIComponent(mechId)}`, {
    method: 'DELETE',
    headers: { ...getAuthHeaders() },
  })
  await handleApiError(res, 'Failed to delete squad mech')
}
