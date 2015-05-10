'use strict';
/* global $ */

/**
 * @ngdoc function
 * @name gitdeployApp.controller:LandingCtrl
 * @description
 * # LandingCtrl
 * Controller of the gitdeployApp
 */
angular.module('gitdeployApp')
    .controller('LandingCtrl', [
        '$window', '$scope',
        function ($window, $scope) {
            $scope.repository = {};
            $scope.showIntegration = function createShield($event) {
                var parser = document.createElement('a'), name;
                $event.stopPropagation();
                parser.href = $scope.repository.url;
                $scope.repository.readmeUrl = parser.href.substr(0, parser.href.length - 4) + '/edit/master/README.md';
                name = parser.pathname.substr(1);
                name = name.substring(0, name.length - 4);
                name = name.split('/')[1];

                // Don't change the order!
                $scope.repository.name = name;
                name = name.replace(/\-/g, '--') + '/master';
                // end

                $scope.repository.badge = 'https://img.shields.io/badge/gitdeploy.io-deploy%20' + (name) + '-green.svg';
                $scope.repository.deployUrl =
                    window.location.protocol + '//' + window.location.host + '/deploy?repository=' + encodeURIComponent($scope.repository.url);
                $scope.repository.showBadge = true;
                $scope.repository.deployUrl =
                    window.location.protocol + '//' + window.location.host + '/deploy?repository=' + encodeURIComponent($scope.repository.url);
                $event.stopPropagation();
                $('#shieldModal').modal('show');
                return false;
            };
        }
    ]);
