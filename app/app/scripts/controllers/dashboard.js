'use strict';

/**
 * @ngdoc function
 * @name gitdeployApp.controller:DashboardCtrl
 * @description
 * # DashboardCtrl
 * Controller of the gitdeployApp
 */
angular.module('gitdeployApp')
    .controller('DashboardCtrl', [
        '$scope', '$routeParams', 'apps', function ($scope, $routeParams, apps) {
            var id = $routeParams.app;
            apps.getApp(id).then(function (data) {
                $scope.$apply(function () {
                    $scope.app = data.data;
                });
            }).catch(function (data, status) {
                if (status === 404) {

                }
            });
        }
    ]);
