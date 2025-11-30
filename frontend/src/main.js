import './assets/base.css'

import { createApp } from 'vue';
import App from './App.vue';
import router from './router';
import { createPinia } from 'pinia';
import { useAuthStore } from './stores/auth';

const app = createApp(App);
app.use(router);
app.use(createPinia());
app.mount('#app');

// Initialize session refresh polling if user has an existing session token
const authStore = useAuthStore();
authStore.initializeSessionRefresh();
