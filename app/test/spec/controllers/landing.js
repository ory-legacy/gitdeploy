'use strict';

describe('Controller: LandingCtrl', function () {

    // load the controller's module
    beforeEach(module('gitdeployApp'));

    var LandingCtrl,
        scope;

    // Initialize the controller and a mock scope
    beforeEach(inject(function ($controller, $rootScope) {
        scope = $rootScope.$new();
        LandingCtrl = $controller('LandingCtrl', {
            $scope: scope
        });
    }));
});
