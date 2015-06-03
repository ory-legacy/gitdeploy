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
        '$window', '$scope', '$http',
        function ($window, $scope, $http) {
            var truncate = function (string, len){
                 if (string.length > len) {
                    return string.substring(0, len) + '...';
                 } else {
                    return string;
                 }
            }, errcb = function () {
                $scope.error = false;
            }, getInfo = function (name, type, cb){
                $http.get('https://api.github.com/repos/' + name + '/' + type).success(function(data) {
                    angular.forEach(data, cb);
                    $scope.ref = $scope.refs[0];
                }).error(errcb);
            }, lastUrl;
            $scope.repository = {url: ''};
            $scope.refs = [];
            $scope.$watchCollection('repository', function() {
                var parser = document.createElement('a'), name;
                if ($scope.repository.url === undefined || !$scope.repository.url.match(/(http|https)\:\/\/github\.com\/[a-zA-Z0-9\-]+\/[a-zA-Z0-9\-\.]+/)){
                    return;
                }
                if ($scope.repository.url === lastUrl) {
                    return;
                }
                lastUrl = $scope.repository.url;
                $scope.refs = [];
                parser.href = $scope.repository.url;
                name = parser.pathname.substr(1);
                getInfo(name, 'branches', function (v){
                    $scope.refs.push({
                        value: 'origin/' + v.name,
                        label: v.name + ' (branch)',
                        name: 'branch/' + v.name
                    });
                });
                getInfo(name, 'tags', function (v){
                   $scope.refs.push({
                       value: 'tags/' + v.name,
                       label: v.name + ' (tag)',
                       name: 'tag/' + v.name
                   });
                });
                getInfo(name, 'commits', function (v){
                   $scope.refs.push({
                       value: v.sha,
                       label: v.committer.login + ': ' + truncate(v.commit.message, 40) + ' (' + v.sha.substr(0,7) + ')',
                       name: 'commit/' + v.sha.substr(0,7)
                   });
                });
            });
            $scope.showIntegration = function ($event) {
                var parser = document.createElement('a'),
                    name = $scope.ref.name;
                $scope.repository.name = name;
                name = name.replace(/\-/g, '--');
                $event.stopPropagation();
                parser.href = $scope.repository.url;
                $scope.repository.readmeUrl = parser.href.substr(0, parser.href.length - 4) + '/edit/master/README.md';
                $scope.repository.badge = 'https://img.shields.io/badge/gitdeploy.io-deploy%20' + (name) + '-green.svg';
                $scope.repository.deployUrl =
                    window.location.protocol + '//' + window.location.host + '/deploy?repository=' + encodeURIComponent($scope.repository.url) + '&ref=' + encodeURIComponent($scope.ref.value);
                $scope.repository.showBadge = true;
                $event.stopPropagation();
                $('#shieldModal').modal('show');
                return false;
            };
        }
    ]);
