import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/player/game-subscription-instances/:game_subscription_instance_id/turn-sheets/:turn_sheet_token',
    name: 'PlayerTurnSheets',
    component: () => import('../views/PlayerTurnSheetView.vue'),
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
  console.log('[player-router] navigating to:', to.fullPath, 'matched:', to.matched.length, 'routes')
  next()
})

export default router
