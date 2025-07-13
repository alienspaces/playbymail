import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import GameView from '../views/GameView.vue'
import StudioLayout from '../components/StudioLayout.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/games',
      name: 'games',
      component: GameView
    },
    {
      path: '/studio',
      redirect: '/studio/games',
    },
    {
      path: '/studio/:gameId',
      component: StudioLayout,
      props: true,
      children: [
        {
          path: 'locations',
          name: 'studio-locations',
          component: () => import('../views/StudioLocationsView.vue'),
          props: true
        },
        {
          path: 'items',
          name: 'studio-items',
          component: () => import('../views/StudioItemsView.vue'),
          props: true
        },
        {
          path: 'creatures',
          name: 'studio-creatures',
          component: () => import('../views/StudioCreaturesView.vue'),
          props: true
        },
        {
          path: 'placement',
          name: 'studio-placement',
          component: () => import('../views/StudioPlacementView.vue'),
          props: true
        }
      ]
    }
  ]
})

export default router
