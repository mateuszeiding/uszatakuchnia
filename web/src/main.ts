import { createApp } from 'vue';

import './scss/_!import.scss';
import App from './App.vue';
import { createAuth } from './auth/createAuth';
import router from './router';
import { VueQueryPlugin } from '@tanstack/vue-query';

const app = createApp(App);

app.use(createAuth()).use(router).use(VueQueryPlugin);

app.mount('#app');
