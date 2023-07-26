const path = require('path')

module.exports = {
    entry: {
        main: path.resolve(__dirname, './dist'),
    },
    output: {
        path: path.resolve(__dirname, './out'),
        filename: '[name].bundle.js',
    },
}