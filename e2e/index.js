const Mocha = require('mocha');
const fs = require('fs');
const path = require('path');

// Instantiate a Mocha instance.
const mocha = new Mocha({
  reporter: 'mochawesome',
  reporterOptions: {
    reportDir: '/report',
  },
});

const testDir = './suites/';

// Add each .js file to the mocha instance
fs.readdirSync(testDir).filter((file) => file.substr(-3) === '.js')
  .forEach(function(file) {
    mocha.addFile(
      path.join(testDir, file)
    );
  });

// Run the tests.
mocha.run(function(failures) {
  process.exitCode = failures ? -1 : 0;
});
