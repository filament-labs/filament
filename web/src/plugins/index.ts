/**
 * plugins/index.ts
 *
 * Automatically included in `./src/main.ts`
 */

// Plugins
import vuetify from './vuetify'
import pinia from '../store'
import { useStore } from '@/store/app'
import router from '../router'
import '../styles/style.scss'
// Types
import type { App } from 'vue'

export function registerPlugins (app: App) {
  app
    .use(vuetify)
    .use(router)
    .use(pinia)
  
  // connect to websocket server
  const store = useStore()
  if (!store.isAppInitialized) {
    store.init()
  }
}
