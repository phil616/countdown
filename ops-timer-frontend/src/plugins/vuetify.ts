import 'vuetify/styles'
import '@mdi/font/css/materialdesignicons.css'
import { createVuetify } from 'vuetify'

export default createVuetify({
  theme: {
    defaultTheme: 'light',
    themes: {
      light: {
        dark: false,
        colors: {
          primary: '#1565C0',
          secondary: '#546E7A',
          accent: '#00BCD4',
          error: '#D32F2F',
          warning: '#F57C00',
          info: '#0288D1',
          success: '#388E3C',
          background: '#F5F5F5',
          surface: '#FFFFFF',
        },
      },
      dark: {
        dark: true,
        colors: {
          primary: '#42A5F5',
          secondary: '#78909C',
          accent: '#26C6DA',
          error: '#EF5350',
          warning: '#FFA726',
          info: '#29B6F6',
          success: '#66BB6A',
          background: '#121212',
          surface: '#1E1E1E',
        },
      },
    },
  },
  defaults: {
    VBtn: { variant: 'flat' },
    VCard: { elevation: 1 },
    VTextField: { variant: 'outlined', density: 'comfortable' },
    VSelect: { variant: 'outlined', density: 'comfortable' },
    VTextarea: { variant: 'outlined', density: 'comfortable' },
  },
})
