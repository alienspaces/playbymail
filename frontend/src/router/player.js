import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
  {
    path: '/player/game-subscription-instances/:game_subscription_instance_id/login/:turn_sheet_token',
    name: 'PlayerTurnSheetLogin',
    component: () => import('../views/PlayerTurnSheetLoginView.vue'),
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
    next({ path: '/player/login' })
  } else {
    next()
  }
})

export default router
