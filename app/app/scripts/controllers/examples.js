'use strict';

/**
 * @ngdoc function
 * @name gitdeployApp.controller:ExamplesCtrl
 * @description
 * # ExamplesCtrl
 * Controller of the gitdeployApp
 */
angular.module('gitdeployApp')
    .controller('ExamplesCtrl', function ($scope) {
        $scope.awesomeThings = [
            'HTML5 Boilerplate',
            'AngularJS',
            'Karma'
        ];
    });
