import { createRouter, createWebHistory } from 'vue-router'
import GameView from '../views/GameView.vue'
import StudioLayout from '../components/StudioLayout.vue'
import LoginView from '../views/LoginView.vue';
import VerifyView from '../views/VerifyView.vue';
import { useAuthStore } from '../stores/auth';

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('../views/HomeView.vue'),
    meta: { requiresAuth: true },
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
  },
  {
    path: '/login',
    name: 'Login',
    component: LoginView,
  },
  {
    path: '/verify',
    name: 'Verify',
    component: VerifyView,
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// Navigation guard for auth
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore();
  const sessionToken = authStore.sessionToken;
  if (to.meta.requiresAuth && !sessionToken) {
    next({ path: '/login' });
  } else {
    next();
  }
});

export default router;
