const path = require('path');

module.exports = {
    mode: 'development',

    entry: {  entry: './tsFiles/eventHolders.ts',
    },
    module: {
        rules: [
            {
                test: /\.tsx?$/,
                use: 'ts-loader',
            },]
    },
    output: {
        filename: 'bundle.js',
        path: path.resolve(__dirname, 'static')
    },
    resolve: {
        extensions: ['.ts', '.js'],
    },
};