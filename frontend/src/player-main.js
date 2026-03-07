import { createApp } from 'vue'
import { createPinia } from 'pinia'

import PlayerApp from './PlayerApp.vue'
import playerRouter from './router/player'

console.log('[player-main] bootstrapping player app, URL:', window.location.href)

const app = createApp(PlayerApp)

app.use(createPinia())
app.use(playerRouter)

const el = document.getElementById('player-app')
console.log('[player-main] mount target #player-app:', el ? 'found' : 'NOT FOUND')

app.mount('#player-app')
console.log('[player-main] app mounted')
