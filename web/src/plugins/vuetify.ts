/**
 * plugins/vuetify.ts
 *
 * Framework documentation: https://vuetifyjs.com`
 */

// Styles
import '@mdi/font/css/materialdesignicons.css'
import 'vuetify/styles'
import { aliases, mdi } from 'vuetify/iconsets/mdi'

// Composables
import { createVuetify } from 'vuetify'

// https://vuetifyjs.com/en/introduction/why-vuetify/#feature-guides

/**
export default createVuetify({
  theme: {
    defaultTheme: 'system',
  },
})**/

// plugins/vuetify.ts
export default createVuetify({
  defaults: {
    global: {
      font: 'ui-sans-serif, system-ui -apple-system BlinkMacSystemFont, Segoe UI, Roboto'
    }
  },
  theme: {
    defaultTheme: 'light',
    themes: {
      light: {
        dark: false,
        colors: {
          background: '#F5F5F5',
          surface: '#FFFFFF',
          primary: '#0066FF',
          secondary: '#00D4FF',
          accent: '#007BFF',
          info: '#00A8FF',
          success: '#4CAF50',
          warning: '#FFC107',
          error: '#FF5252',
          text: '#0A0A0A',
          border: '#E0E0E0',
          'gray-600': '#4B5563'
        },
        variables: {
          'primary-gradient': 'linear-gradient(135deg, #00D4FF 0%, #0066FF 100%)',
          'secondary-gradient': 'linear-gradient(135deg, rgb(0, 212, 255) 0%, rgb(0, 102, 255) 100%)',
          'card-gradient': 'linear-gradient(135deg, rgba(0,212,255,0.1) 0%, rgba(0,102,255,0.1) 100%)',
          'gray-600': '#4B5563'
        },
      },
      dark: {
        dark: true,
        colors: {
          background: '#0D1117',
          surface: '#161B22',
          primary: '#2196F3',
          secondary: '#00B8FF',
          accent: '#3399FF',
          info: '#33BFFF',
          success: '#4CAF50',
          warning: '#FFC107',
          error: '#FF5252',
          text: '#EAEAEA',
          border: '#2D333B',
          'gray-600': '#4B5563'
        },
        variables: {
          'primary-gradient': 'linear-gradient(135deg, rgb(0, 102, 255) 0%, rgb(0, 212, 255) 100%)',
          'card-gradient': 'linear-gradient(135deg, rgba(0,212,255,0.05) 0%, rgba(0,102,255,0.05) 100%)',
        },
      },
    },
  },
  icons: {
    defaultSet: 'mdi',
    aliases,
    sets: { mdi },
  },
})

