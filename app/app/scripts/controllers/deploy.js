'use strict';
/* global EventSource */

/**
 * @ngdoc function
 * @name gitdeployApp.controller:DeployCtrl
 * @description
 * # DeployCtrl
 * Controller of the gitdeployApp
 */
angular.module('gitdeployApp')
    .controller('DeployCtrl', [
        '$scope', '$routeParams', '$http', '$rootScope', 'endpoint',
        function ($scope, $routeParams, $http, $rootScope, endpoint) {
            var repository = $routeParams.repository,
                sse = function (app) {
                    var url = endpoint.sse + '/deployments/' + app + '/events', ev = new EventSource(url);
                    ev.addEventListener('open', function (e) {
                        console.log('Channel opened!', e);
                    });

                    return {
                        addEventListener: function (eventName, callback) {
                            ev.addEventListener(eventName, function () {
                                $rootScope.$apply(function () {
                                    callback.apply(sse);
                                });
                            });
                        }
                    };
                };

            $scope.logs = [];

            if (repository === undefined || repository.length < 1) {
                $scope.error = 'The repository query parameter is missing.';
                return;
            }

            $scope.error = false;
            $http.post(endpoint + '/deployments', {repository: repository}).
                success(function (data) {
                    var el = sse(data.data.id);
                    el.addEventListener('message', function (e) {
                        var message;
                        try {
                            message = JSON.parse(e.data);
                        } catch (exc) {
                            console.log(exc);
                            return;
                        }

                        if (message.eventName !== 'app.deployed') {
                            $scope.logs.push(message.data.replace(/(\r\n|\r|\n)/gm, '\n'));
                            $scope.logMessages = $scope.logs.join('\n');
                        } else {
                            $scope.deployUrl = message.data;
                        }
                    });

                    el.addEventListener('error', function (e) {
                        $scope.error = 'The backend server does not respond correctly or closed the connection.';
                        e.currentTarget.close();
                    });
                }).error(function (data) {
                    if (data === null || data.error === undefined) {
                        $scope.error = 'The backend server returned an error: No response was given, come back later.';
                    } else {
                        $scope.error =
                            'The backend server returned an error: ' + (data.error.message || 'No response was given, come back later.');
                    }
                });
        }
    ]);