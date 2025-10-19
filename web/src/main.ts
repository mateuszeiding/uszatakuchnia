import { createApp } from 'vue'
import './scss/_!import.scss'
import App from './App.vue'
import router from './router'
import { VueQueryPlugin } from '@tanstack/vue-query'

const app = createApp(App)

app.use(router).use(VueQueryPlugin)

app.mount('#app')
