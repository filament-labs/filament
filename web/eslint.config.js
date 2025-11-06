import vuetify from 'eslint-config-vuetify'

export default vuetify({
  vue: true,
  ts: {
    preset: 'recommended',
  },
  rules: {
    '@typescript-eslint/method-signature-style': 'off',
  },
})
