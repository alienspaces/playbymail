import { createRouter, createWebHistory } from 'vue-router'
import GameView from '../views/GameView.vue'
import LoginView from '../views/LoginView.vue';
import VerifyView from '../views/VerifyView.vue';
import FaqView from '../views/FaqView.vue';
import StudioLayout from '../components/StudioLayout.vue';

import { useAuthStore } from '../stores/auth';

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('../views/HomeView.vue'),
  },
  {
    path: '/studio',
    component: StudioLayout,
    children: [
      { path: '', name: 'StudioGames', component: GameView },
      // Adventure game type studio views
      { path: ':gameId/locations', component: () => import('../views/studio/adventure/StudioLocationsView.vue') },
      { path: ':gameId/location-links', component: () => import('../views/studio/adventure/StudioLocationLinksView.vue') },
      { path: ':gameId/items', component: () => import('../views/studio/adventure/StudioItemsView.vue') },
      { path: ':gameId/creatures', component: () => import('../views/studio/adventure/StudioCreaturesView.vue') },
      { path: ':gameId/item-placements', component: () => import('../views/studio/adventure/StudioItemPlacementsView.vue') },
      { path: ':gameId/creature-placements', component: () => import('../views/studio/adventure/StudioCreaturePlacementsView.vue') },
      { path: ':gameId/turn-sheet-backgrounds', component: () => import('../views/studio/adventure/StudioTurnSheetBackgroundsView.vue') },

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
    path: '/player/game-subscriptions/:game_subscription_id/game-instances/:game_instance_id/login/:turn_sheet_token',
    name: 'PlayerTurnSheetLogin',
    component: () => import('../views/PlayerTurnSheetLoginView.vue'),
  },
  {
    path: '/faq',
    name: 'Faq',
    component: FaqView,
  },
  {
    path: '/account',
    component: () => import('../components/AccountLayout.vue'),
    children: [
      { path: '', name: 'AccountProfile', component: () => import('../views/account/AccountProfileView.vue') },
      { path: 'contacts', name: 'AccountContacts', component: () => import('../views/account/AccountContactsView.vue') },
    ],
  },
  {
    path: '/admin',
    component: () => import('../components/ManagementLayout.vue'),
    children: [
      { path: '', name: 'ManagementDashboard', component: () => import('../views/management/ManagementGamesDashboardView.vue') },
      { path: 'games/:gameId/instances', name: 'ManagementGameInstances', component: () => import('../views/management/ManagementGameInstancesView.vue') },
      { path: 'games/:gameId/instances/create', name: 'ManagementCreateInstance', component: () => import('../views/management/ManagementCreateInstanceView.vue') },
      { path: 'games/:gameId/instances/:instanceId', name: 'ManagementInstanceDetail', component: () => import('../views/management/ManagementInstanceDetailView.vue') },
      { path: 'games/:gameId/turn-sheets', name: 'ManagementTurnSheets', component: () => import('../views/management/ManagementTurnSheetsView.vue') },
    ],
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
