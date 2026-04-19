import { baseUrl, getAuthHeaders, apiFetch, handleApiError } from './baseUrl';

const BASE = (gameId) => `${baseUrl}/api/v1/mecha-games/${encodeURIComponent(gameId)}/equipments`

export async function fetchMechaGameEquipment(gameId, params = {}) {
  const url = new URL(BASE(gameId))
  if (params.page_number) url.searchParams.set('page_number', params.page_number)
  const res = await apiFetch(url.toString(), { headers: { ...getAuthHeaders() } })
  await handleApiError(res, 'Failed to fetch equipment')
  const json = await res.json()
  const pagination = JSON.parse(res.headers.get('X-Pagination') || '{}')
  return { data: json.data || [], hasMore: !!pagination.has_more }
}

export async function createMechaGameEquipment(gameId, data) {
  const res = await apiFetch(BASE(gameId), {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data),
  })
  await handleApiError(res, 'Failed to create equipment')
  const json = await res.json()
  return json.data
}

export async function updateMechaGameEquipment(gameId, equipmentId, data) {
  const res = await apiFetch(`${BASE(gameId)}/${encodeURIComponent(equipmentId)}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data),
  })
  await handleApiError(res, 'Failed to update equipment')
  const json = await res.json()
  return json.data
}

export async function deleteMechaGameEquipment(gameId, equipmentId) {
  const res = await apiFetch(`${BASE(gameId)}/${encodeURIComponent(equipmentId)}`, {
    method: 'DELETE',
    headers: { ...getAuthHeaders() },
  })
  await handleApiError(res, 'Failed to delete equipment')
}
