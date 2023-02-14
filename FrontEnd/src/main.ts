import { createApp } from 'vue'
import App from './App.vue'
import mitt from 'mitt';
import './assets/main.css';
import VueDraggable from 'vue-draggable'

const emitter = mitt();
const app = createApp(App);
app.use(VueDraggable)
app.config.globalProperties.emitter = emitter;
app.mount('#app');