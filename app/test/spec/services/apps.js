'use strict';

describe('Service: apps', function () {

  // load the service's module
  beforeEach(module('gitdeployApp'));

  // instantiate service
  var apps;
  beforeEach(inject(function (_apps_) {
    apps = _apps_;
  }));

});
