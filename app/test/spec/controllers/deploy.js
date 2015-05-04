'use strict';

describe('Controller: DeployCtrl', function () {

    // load the controller's module
    beforeEach(module('gitdeployApp'));

    var DeployCtrl,
        scope;

    // Initialize the controller and a mock scope
    beforeEach(inject(function ($controller, $rootScope) {
        scope = $rootScope.$new();
        DeployCtrl = $controller('DeployCtrl', {
            $scope: scope
        });
    }));
});
