module.exports = {
  root: true,
  parser: 'babel-eslint',
  parserOptions: {
    sourceType: 'module'
  },
  env: {
    'browser': true
  },
  globals: {
    $: true,
    Vue: true,
    Regular: true,
    __DEV__: true,
    __PROD__: true
  },
  extends: 'airbnb-base',
  plugins: [
    'html'
  ],
  rules: {
    'no-debugger': global.FEB_PROD ? 2 : 0,
    'no-console': global.FEB_PROD ? 1 : 0,
    'func-names': 0,
    'import/no-unresolved': 0,
    'import/extensions': 0,
    'import/no-extraneous-dependencies': 0,
    'linebreak-style': 0
  }
};
