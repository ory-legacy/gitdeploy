'use strict';

angular.module('gdApp.home', ['ngRoute']).config(['$routeProvider', function ($routeProvider) {
    $routeProvider.when('/', {
        controller: 'homeCtrl'
    });
}]).controller('homeCtrl', ['$window', '$scope',
    function ($window, $scope) {
        $scope.repository = {};
        $scope.createShield = function createShield($event) {
            var parser = document.createElement('a'), name;
            $event.stopPropagation();
            parser.href = $scope.repository.url;
            $scope.repository.readmeUrl = parser.href.substr(0, parser.href.length - 4) + '/edit/master/README.md';
            name = parser.pathname.substr(1);
            name = name.substring(0, name.length - 4);
            name = name.split('/')[1];

            // Don't change the order!
            $scope.repository.name = name;
            name = name.replace(/\-/g, '--') + "/master";
            // end

            $scope.repository.badge = "https://img.shields.io/badge/gitdeploy.io-deploy%20" + (name) + "-green.svg";
            $scope.repository.deployUrl = window.location.protocol + '//' + window.location.host + '/deploy?repository=' +  encodeURIComponent($scope.repository.url);
            $scope.repository.showBadge = true;
            $scope.repository.deployUrl = window.location.protocol + "//" + window.location.host + "/deploy?repository="
                + encodeURIComponent($scope.repository.url);
            return false;
        };
        $scope.showIntegration = function showIntegration($event) {
            $event.stopPropagation();
            $('#shieldModal').modal('show');
            return false;
        }
    }
]);