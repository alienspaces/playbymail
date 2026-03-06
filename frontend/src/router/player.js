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
    meta: { requiresAuth: true },
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
    const gsiId = to.params.game_subscription_instance_id
    if (gsiId) {
      next({ name: 'PlayerTurnSheetLogin', params: { game_subscription_instance_id: gsiId, turn_sheet_token: 'expired' } })
    } else {
      // Redirect to main app login with a return path so the user comes back after authenticating
      window.location.href = '/login?redirect=' + encodeURIComponent(to.fullPath)
    }
  } else {
    next()
  }
})

export default router
