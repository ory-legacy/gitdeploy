'use strict';

/**
 * @ngdoc function
 * @name gitdeployApp.controller:NavCtrl
 * @description
 * # NavCtrl
 * Controller of the gitdeployApp
 */
angular.module('gitdeployApp')
    .controller('NavCtrl', function ($scope) {
        $scope.awesomeThings = [
            'HTML5 Boilerplate',
            'AngularJS',
            'Karma'
        ];
    });
