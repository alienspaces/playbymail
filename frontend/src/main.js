import './assets/base.css'

import { createApp } from 'vue';
import App from './App.vue';
import router from './router';
import { createPinia } from 'pinia';
import { useAuthStore } from './stores/auth';

console.log('[main] bootstrapping MAIN app, URL:', window.location.href)
console.log('[main] document has #app:', !!document.getElementById('app'))
console.log('[main] document has #player-app:', !!document.getElementById('player-app'))

const app = createApp(App);
app.use(router);
app.use(createPinia());
app.mount('#app');
console.log('[main] MAIN app mounted')

// Initialize session refresh polling if user has an existing session token
const authStore = useAuthStore();
authStore.initializeSessionRefresh();
