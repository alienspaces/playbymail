import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
  {
    path: '/player/game-subscription-instances/:game_subscription_instance_id/login/:turn_sheet_token',
    name: 'PlayerTurnSheetLogin',
    component: () => import('../views/PlayerTurnSheetLoginView.vue'),
  },
  {
    path: '/player/game-subscription-instances/:game_subscription_instance_id/turn-sheets',
    name: 'PlayerTurnSheets',
    component: () => import('../views/PlayerTurnSheetView.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/player/join-game/:game_subscription_id',
    name: 'PlayerJoinGame',
    component: () => import('../views/PlayerJoinGameView.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  const sessionToken = authStore.sessionToken
  if (to.meta.requiresAuth && !sessionToken) {
    // If the route carries a GSI id, redirect to the token-request path for that instance
    // so the player can receive a new email link. Otherwise fall through to the games catalog.
    const gsiId = to.params.game_subscription_instance_id
    if (gsiId) {
      next({ name: 'PlayerTurnSheetLogin', params: { game_subscription_instance_id: gsiId, turn_sheet_token: 'expired' } })
    } else {
      next({ path: '/games' })
    }
  } else {
    next()
  }
})

export default router
