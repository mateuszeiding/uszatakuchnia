import { VueQueryPlugin } from '@tanstack/vue-query';
import { createApp } from 'vue';

import App from './App.vue';
import { createAuth } from './auth/createAuth';
import router from './router';
import { vAutosize } from './shared/directives/autosize';

const app = createApp(App);

app.use(createAuth()).use(router).use(VueQueryPlugin);
app.directive('autosize', vAutosize);

app.mount('#app');
