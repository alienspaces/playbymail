import { createApp } from 'vue'
import { createPinia } from 'pinia'

import PlayerApp from './PlayerApp.vue'
import playerRouter from './router/player'

const app = createApp(PlayerApp)

app.use(createPinia())
app.use(playerRouter)

app.mount('#player-app')
