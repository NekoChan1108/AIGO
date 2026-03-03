import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { MotionPlugin } from '@vueuse/motion'
import router from './router/index.ts'
import './style.css'
import App from './App.vue'

// Polyfill for structuredClone if it doesn't exist
if (typeof globalThis.structuredClone === 'undefined') {
  globalThis.structuredClone = <T>(value: T): T => {
    return JSON.parse(JSON.stringify(value))
  }
}

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(MotionPlugin)

app.mount('#app')