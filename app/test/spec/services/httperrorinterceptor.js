'use strict';

describe('Service: httperrorinterceptor', function () {

  // load the service's module
  beforeEach(module('gitdeployApp'));

  // instantiate service
  var httperrorinterceptor;
  beforeEach(inject(function (_httperrorinterceptor_) {
    httperrorinterceptor = _httperrorinterceptor_;
  }));

  it('should do something', function () {
    expect(!!httperrorinterceptor).toBe(true);
  });

});
