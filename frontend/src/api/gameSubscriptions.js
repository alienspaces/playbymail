import { baseUrl, getAuthHeaders, apiFetch, handleApiError } from './baseUrl';

export async function getMyGameSubscriptions() {
  const res = await apiFetch(`${baseUrl}/api/v1/game-subscriptions`, {
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to fetch game subscriptions');
  return await res.json();
}

export async function createGameSubscription(gameId, subscriptionType, instanceLimit = null) {
  const body = {
    game_id: gameId,
    subscription_type: subscriptionType,
  };
  if (instanceLimit !== null) {
    body.instance_limit = instanceLimit;
  }
  const res = await apiFetch(`${baseUrl}/api/v1/game-subscriptions`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify(body),
  });
  await handleApiError(res, 'Failed to create game subscription');
  return await res.json();
}

export async function cancelGameSubscription(subscriptionId) {
  const res = await apiFetch(`${baseUrl}/api/v1/game-subscriptions/${subscriptionId}`, {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to cancel game subscription');
  if (res.status === 204) {
    return null;
  }
  return await res.json();
}

export async function linkGameInstanceToSubscription(subscriptionId, instanceId) {
  const res = await apiFetch(`${baseUrl}/api/v1/game-subscriptions/${subscriptionId}/instances`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify({
      game_subscription_id: subscriptionId,
      game_instance_id: instanceId,
    }),
  });
  await handleApiError(res, 'Failed to link game instance to subscription');
  return await res.json();
}

export async function unlinkGameInstanceFromSubscription(subscriptionId, instanceId) {
  const res = await apiFetch(`${baseUrl}/api/v1/game-subscriptions/${subscriptionId}/instances/${instanceId}`, {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to unlink game instance from subscription');
  if (res.status === 204) {
    return null;
  }
  return await res.json();
}

export async function getSubscriptionInstances(subscriptionId) {
  const res = await apiFetch(`${baseUrl}/api/v1/game-subscriptions/${subscriptionId}`, {
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
  });
  await handleApiError(res, 'Failed to get subscription instances');
  const data = await res.json();
  return data.data?.game_instance_ids || [];
}
