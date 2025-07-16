import { createRouter, createWebHistory } from 'vue-router'
import GameView from '../views/GameView.vue'
import LoginView from '../views/LoginView.vue';
import VerifyView from '../views/VerifyView.vue';
import FaqView from '../views/FaqView.vue';
import StudioEntryView from '../views/StudioEntryView.vue';
import AdminEntryView from '../views/AdminEntryView.vue';
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
    name: 'StudioEntry',
    component: StudioEntryView,
    children: [
      // Existing studio routes (e.g., games, locations, items, creatures, placement)
      // These will be rendered when logged in
      {
        path: ':gameId/locations',
        component: () => import('../views/StudioLocationsView.vue'),
      },
      {
        path: ':gameId/items',
        component: () => import('../views/StudioItemsView.vue'),
      },
      {
        path: ':gameId/creatures',
        component: () => import('../views/StudioCreaturesView.vue'),
      },
      {
        path: ':gameId/placement',
        component: () => import('../views/StudioPlacementView.vue'),
      },
    ],
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
  {
    path: '/faq',
    name: 'Faq',
    component: FaqView,
  },
  {
    path: '/admin',
    name: 'AdminEntry',
    component: AdminEntryView,
    // Placeholder: add children for admin features as implemented
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
