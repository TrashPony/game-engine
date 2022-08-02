const path = require('path');
const webpack = require('webpack');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const {CleanWebpackPlugin} = require('clean-webpack-plugin');
const UglifyJSPlugin = require('uglifyjs-webpack-plugin');

module.exports = {
  entry: './src/main.js',
  output: {
    path: path.resolve(__dirname, 'dist'),
    publicPath: '',
    filename: 'build-[hash:6].js'
  },
  module: {
    rules: [
      {
        test: /\.css$/,
        use: [
          'vue-style-loader',
          'css-loader'
        ],
      }, {
        test: /\.vue$/,
        loader: 'vue-loader',
        options: {
          loaders: {}
          // other vue-loader options go here
        }
      }, {
        test: /\.js$/,
        loader: 'babel-loader',
        exclude: /node_modules/
      }, {
        test: [/\.vert$/, /\.frag$/],
        use: "raw-loader"
      }, {
        test: /\.(png|svg|jpe?g)$/i,
        use: [
          {
            loader: 'img-optimize-loader',
            options: {
              compress: {
                mode: 'low', // 'lossless', 'low', 'high'
                disableOnDevelopment: true,
              },
              name: '[contenthash].[ext]',
              inline: {
                limit: false,
              },
            },
          },
        ],
      }, {
        test: /\.(gif|xml|json|ogg|mp3|wav|mpe?g)$/i,
        use: "file-loader",
      }
    ]
  },
  resolve: {
    alias: {
      'vue$': 'vue/dist/vue.esm.js',
      'assets': path.resolve(__dirname, 'src/assets/')
    },
    extensions: ['*', '.js', '.vue', '...']
  },
  devServer: {
    historyApiFallback: true,
    noInfo: true,
    overlay: true,
    port: 8083,
  },
  performance: {
    hints: false
  },
  devtool: '#eval-source-map',
  plugins: [
    new webpack.DefinePlugin({
      'NODE_ENV': JSON.stringify(process.env.NODE_ENV),
      CANVAS_RENDERER: JSON.stringify(true),
      WEBGL_RENDERER: JSON.stringify(true),
    }),
    new CleanWebpackPlugin(),
    new HtmlWebpackPlugin({
      template: 'index.html',
    }),
    new webpack.ProvidePlugin({
      $: "jquery",
      jQuery: "jquery",
      "window.jQuery": "jquery'",
      "window.$": "jquery"
    })
  ]
};

if (process.env.NODE_ENV === 'production') {
  module.exports.devtool = '#source-map';
  // http://vue-loader.vuejs.org/en/workflow/production.html
  module.exports.plugins = (module.exports.plugins || []).concat([
    new webpack.DefinePlugin({
      'process.env': {
        NODE_ENV: '"production"'
      }
    }),
    new UglifyJSPlugin({
        "uglifyOptions":
          {
            compress: {
              warnings: false
            },
            sourceMap: true
          }
      }
    ),
    new webpack.LoaderOptionsPlugin({
      minimize: true
    })
  ])
}
