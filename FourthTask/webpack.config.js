const path = require('path');

module.exports = {
    entry: {
        main: ['./dist/editHtml.js', './dist/addEvents.ts', './dist/handleJson.js']
    },
    output: {
        filename: 'bundle.js',
        path: path.resolve(__dirname, 'static'),
    },
};