import { createApp } from 'vue'
import './scss/_!import.scss'
import App from './App.vue'
import router from './router'
import { VueQueryPlugin } from '@tanstack/vue-query'
import { createAuth } from './auth/createAuth'

const app = createApp(App)

app.use(createAuth()).use(router).use(VueQueryPlugin)

app.mount('#app')
