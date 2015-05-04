'use strict';

/**
 * @ngdoc function
 * @name gitdeployApp.controller:DashboardCtrl
 * @description
 * # DashboardCtrl
 * Controller of the gitdeployApp
 */
angular.module('gitdeployApp')
    .controller('DashboardCtrl', ['scope', '$routeParams', 'apps', function ($scope, $routeParams, apps) {
        var id = $routeParams.id;
        apps.getApp(id).then(function (data){
            $scope.app = data.data;
            $scope.$apply();
        });
    }]);
