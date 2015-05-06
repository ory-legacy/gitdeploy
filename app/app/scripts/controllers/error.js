'use strict';

/**
 * @ngdoc function
 * @name gitdeployApp.controller:ErrorCtrl
 * @description
 * # ErrorCtrl
 * Controller of the gitdeployApp
 */
angular.module('gitdeployApp')
    .controller('ErrorCtrl', [
        '$rootScope', '$scope', function ($scope, $rootScope) {
            $rootScope.$on('error', function () {
                console.log('err', $rootScope.error);
                $scope.error = $rootScope.error;
            });
            $scope.reload = function () {
                window.location.reload();
            };
        }
    ]);
