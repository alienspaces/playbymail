import { baseUrl, getAuthHeaders, apiFetch, handleApiError } from './baseUrl';

const BASE = (gameId) => `${baseUrl}/api/v1/mecha-games/${encodeURIComponent(gameId)}/computer-opponents`

export async function fetchMechaGameComputerOpponents(gameId, params = {}) {
  const url = new URL(BASE(gameId))
  if (params.page_number) url.searchParams.set('page_number', params.page_number)
  const res = await apiFetch(url.toString(), { headers: { ...getAuthHeaders() } })
  await handleApiError(res, 'Failed to fetch computer opponents')
  const json = await res.json()
  const pagination = JSON.parse(res.headers.get('X-Pagination') || '{}')
  return { data: json.data || [], hasMore: !!pagination.has_more }
}

export async function createMechaGameComputerOpponent(gameId, data) {
  const res = await apiFetch(BASE(gameId), {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data),
  })
  await handleApiError(res, 'Failed to create computer opponent')
  const json = await res.json()
  return json.data
}

export async function updateMechaGameComputerOpponent(gameId, opponentId, data) {
  const res = await apiFetch(`${BASE(gameId)}/${encodeURIComponent(opponentId)}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(data),
  })
  await handleApiError(res, 'Failed to update computer opponent')
  const json = await res.json()
  return json.data
}

export async function deleteMechaGameComputerOpponent(gameId, opponentId) {
  const res = await apiFetch(`${BASE(gameId)}/${encodeURIComponent(opponentId)}`, {
    method: 'DELETE',
    headers: { ...getAuthHeaders() },
  })
  await handleApiError(res, 'Failed to delete computer opponent')
}
