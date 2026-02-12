const path = require('node:path');
const EXE = process.platform === 'win32' ? '.exe' : '';
const CLICKER = path.join(__dirname, '../clicker/bin/clicker') + EXE;
module.exports = { CLICKER };
