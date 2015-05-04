'use strict';

describe('gdApp.version module', function() {
  beforeEach(module('gdApp.version'));

  describe('version service', function() {
    it('should return current version', inject(function(version) {
      expect(version).toEqual('0.1');
    }));
  });
});
