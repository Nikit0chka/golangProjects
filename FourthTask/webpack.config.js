const path = require('path');

module.exports = {
    mode: 'development',
    entry: {
        main: './jsFiles/draw.js',
        vendor1: './jsFiles/eventHolders.js',
        vendor2: './jsFiles/model.js'
    },
    output: {
        filename: '[name].bundle.js',
        path: path.resolve(__dirname, '/static')
    }
};