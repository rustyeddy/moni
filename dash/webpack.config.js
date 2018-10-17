const Webpack = require('webpack');
const Path = require('path')
const UglifyJsPlugin = require('uglifyjs-webpack-plugin');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const OptimizeCSSAssetsPlugin = require('optimize-css-assets-webpack-plugin');
const CopyWebpackPlugin = require('copy-webpack-plugin');
const HardSourceWebpackPlugin = require('hard-source-webpack-plugin');

const opts = {
  rootDir: process.cwd(),
  devBuild: process.env.NODE_ENV !== 'production',
};

module.exports = {
  entry: {
    app: './src/js/app',
    charts: './src/js/charts',
    forms: './src/js/forms',
    maps: './src/js/maps',
    tables: './src/js/tables'
  },
  mode: process.env.NODE_ENV === 'production' ? 'production' : 'development',
  devtool: process.env.NODE_ENV === 'production' ? 'source-map' : 'inline-source-map',
  output: {
    path: Path.join(opts.rootDir, 'dist'),
    pathinfo: opts.devBuild,
    filename: 'js/[name].js'
  },
  performance: { hints: false },
  optimization: {
    minimizer: [
      new UglifyJsPlugin({
        cache: true,
        parallel: true,
        sourceMap: false
      }),
      new OptimizeCSSAssetsPlugin({})
    ]
  },
  plugins: [
    // Extract css files to seperate bundle
    new MiniCssExtractPlugin({
      filename: 'css/[name].css',
      chunkFilename: 'css/[id].css'
    }),
    // jQuery and PopperJS
    new Webpack.ProvidePlugin({
      $: 'jquery',
      jQuery: 'jquery',
      'window.$': 'jquery',
      'window.jQuery': 'jquery',
      Popper: ['popper.js', 'default']
    }),
    // Copy fonts and images to dist
    new CopyWebpackPlugin([
      {from: 'src/fonts', to: 'fonts'},
      {from: 'src/img', to: 'img'}
    ]),
    // Copy dist folders to docs/dist (for demo pages)
    new CopyWebpackPlugin([
      {from: 'dist', to: '../docs'}
    ]),
    // Speed up webpack build
    new HardSourceWebpackPlugin()
  ],
  module: {
    rules: [
      // Babel-loader
      {
        test: /\.js$/,
        exclude: /(node_modules)/,
        loader: ['babel-loader?cacheDirectory=true']
      },
      // Css-loader & sass-loader
      {
        test: /\.scss$/,
        use: [
          MiniCssExtractPlugin.loader,
          'css-loader',
          'postcss-loader',
          'sass-loader'
        ]
      },
      // Load fonts
      {
        test   : /\.(ttf|eot|svg|woff(2)?)(\?[a-z0-9=&.]+)?$/,
        loader : 'url-loader',
        options: {
          limit: 100000
        }
      },
      // Load images
      {
        test   : /\.(png|jpg|jpeg|gif?)(\?[a-z0-9=&.]+)?$/,
        loader : 'url-loader',
        options: {
          limit: 100000
        }
      },
      // Expose loader
      {
        test: require.resolve('jquery'),
        use: [{
            loader: 'expose-loader',
            options: 'jQuery'
        },{
            loader: 'expose-loader',
            options: '$'
        }]
      }
    ]
  },
  resolve: {
    extensions: ['.js', '.scss'],
    modules: ['node_modules'],
    alias: {
      request$: 'xhr'
    }
  }
}
