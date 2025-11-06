import vuetify from 'eslint-config-vuetify'

export default [
  ...vuetify(), // bring in Vuetify’s rules first
  {
    rules: {
      'antfu/top-level-function': 'off',
      '@typescript-eslint/func-style': 'off',
      // override Vuetify's func-style rule
      'func-style': ['error', 'declaration', { allowArrowFunctions: true }],
    },
  },
]
