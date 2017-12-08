module.exports = {
  js: {
    assets: true,
    eslint: false
  },
  css: {
    px2rem: false,
    autoprefixer: true,
  },
  webpack: {
    // module: {
    //   loaders: [
    //     {test: /.(png|woff(2)?|eot|ttf|svg)(?[a-z0-9=.]+)?$/, loader: 'url-loader?limit=100000',exclude: /node_modules/},
    //     {test: /\.css$/, loader: 'css-loader',exclude: /node_modules/},
    //   ]
    // },
    externals: {
      vue: 'Vue',
      jquery: 'jQuery',
    },
    provide: {
      Vue: 'vue',
      $: 'jquery',
    }
  }
};
