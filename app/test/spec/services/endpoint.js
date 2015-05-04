'use strict';

describe('Service: endpoint', function () {

  // load the service's module
  beforeEach(module('gitdeployApp'));

  // instantiate service
  var endpoint;
  beforeEach(inject(function (_endpoint_) {
    endpoint = _endpoint_;
  }));

  it('should do something', function () {
    expect(!!endpoint).toBe(true);
  });

});
