'use strict';

describe('Service: config', function () {

  // load the service's module
  beforeEach(module('gitdeployApp'));

  // instantiate service
  var config;
  beforeEach(inject(function (_config_) {
    config = _config_;
  }));

});
