import { createRouter, createWebHistory } from 'vue-router'
import GameView from '../views/GameView.vue'
import LoginView from '../views/LoginView.vue';
import VerifyView from '../views/VerifyView.vue';
import FaqView from '../views/FaqView.vue';
import StudioLayout from '../components/StudioLayout.vue';
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
    path: '/studio',
    component: StudioLayout,
    meta: { requiresAuth: true },
    children: [
      { path: '', name: 'StudioGames', component: GameView },
      // Adventure game type studio views
      { path: ':gameId/locations', component: () => import('../views/studio/adventure/StudioLocationsView.vue') },
      { path: ':gameId/location-links', component: () => import('../views/studio/adventure/StudioLocationLinksView.vue') },
      { path: ':gameId/items', component: () => import('../views/studio/adventure/StudioItemsView.vue') },
      { path: ':gameId/creatures', component: () => import('../views/studio/adventure/StudioCreaturesView.vue') },
      { path: ':gameId/item-placements', component: () => import('../views/studio/adventure/StudioItemPlacementsView.vue') },
      { path: ':gameId/creature-placements', component: () => import('../views/studio/adventure/StudioCreaturePlacementsView.vue') },
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
    path: '/account',
    name: 'Account',
    component: () => import('../views/AccountView.vue'),
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
