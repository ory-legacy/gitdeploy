'use strict';
/* global moment */

/**
 * @ngdoc function
 * @name gitdeployApp.controller:DashboardCtrl
 * @description
 * # DashboardCtrl
 * Controller of the gitdeployApp
 */
angular.module('gitdeployApp')
    .controller('DashboardCtrl', [
        '$scope', '$routeParams', 'apps', 'config', function ($scope, $routeParams, apps, config) {
            var id = $routeParams.app;
            apps.getApp(id).then(function (data) {
                var dl = [], ps;
                angular.forEach(data.data.deployLogs, function (v) {
                    try {
                        var message = JSON.parse(v.message);
                        dl.unshift(message.data.replace(/(\r\n|\r|\n)/gm, '\n'));
                    } catch (exc) {
                        console.log(exc);
                    }
                });
                data.data.deployLogs = dl.join('\n');
                ps = data.data.ps.split('\n');
                data.data.ps = [];
                data.data.expiresAt = moment(data.data.expiresAt);
                data.data.createdAt = moment(data.data.createdAt);
                angular.forEach(ps, function (v) {
                    v = v.split(/\s+/);
                    data.data.ps.push({id: v[0], type: v[1]});
                    if (v[1] === 'web') {
                        $scope.noWebProcess = false;
                        $scope.$apply();
                    }
                });
                data.data.ps.splice(0, 1);
                $scope.app = data.data;
                $scope.$apply();

                config.get().then(function (response) {
                    $scope.serverTime = moment(response.data.time);
                    var currentTTL = Math.round(moment.duration($scope.app.expiresAt.diff($scope.serverTime)).asMinutes()),
                        ttl = Math.round(moment.duration($scope.app.expiresAt.diff($scope.app.createdAt)).asMinutes()) - currentTTL;
                    $scope.data = [ttl, currentTTL];
                    $scope.$apply();
                });
            });

            $scope.labels = ['Time used', 'Time available'];
            $scope.data = [0, 0];
            $scope.colors = ['#DCDCDC', '#97BBCD']; // grey, blue
            $scope.noWebProcess = true;
        }
    ]);
