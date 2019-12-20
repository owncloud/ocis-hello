const VueLoaderPlugin = require('vue-loader/lib/plugin')

module.exports = {
  plugins: [
    new VueLoaderPlugin()
  ],
  entry: {
    files: [
      './ui/default.js'
    ]
  },
  output: {
    path: require('path').resolve(__dirname, 'assets'),
    publicPath: 'apps/hello/',
    chunkFilename: '[name].hello.bundle.js',
    filename: 'hello.bundle.js'
  },
  module: {
    rules: [{
      test: /\.js?$/,
      exclude: /node_modules/,
      loader: 'babel-loader'
    }, {
      test: /\.vue$/,
      loader: 'vue-loader'
    }, {
      enforce: 'pre',
      test: /\.(js|vue)$/,
      exclude: /node_modules/,
      loader: 'eslint-loader'
    }, {
      test: /\.css$/,
      use: [
        'vue-style-loader',
        'css-loader'
      ]
    }]
  }
}
