'use strict';

angular.module('gdApp.deploy', ['ngRoute']).config(['$routeProvider', function ($routeProvider) {
    $routeProvider.when('/deploy/', {
        controller: 'deployCtrl',
        templateUrl: 'app/deploy/deploy.html'
    });
}]).controller('deployCtrl', ['$scope', '$routeParams', '$http', '$rootScope',
    function ($scope, $routeParams, $http, $rootScope) {
        var repository = $routeParams.repository, sse;
        $scope.logs = [];
        sse = function (app) {
            var url = endpoint + '/deployments/' + app + '/events';
            var ev = new EventSource(url);
            ev.addEventListener('open', function (e) {
                console.log("Channel opened!", e);
            });
            return {
                addEventListener: function(eventName, callback) {
                    ev.addEventListener(eventName, function() {
                        var args = arguments;
                        $rootScope.$apply(function () {
                            callback.apply(sse, args);
                        });
                    });
                }
            };
        };

        if (repository === undefined || repository.length < 1) {
            console.log("No repository provided");
            $scope.error = "The repository query parameter is missing.";
            return;
        }
        $scope.error = false;

        $http.post(endpoint + '/deployments', {repository: repository}).
        success(function(data, status, headers, config) {
            var el = sse(data.data.id);
            el.addEventListener('message', function (e) {
                try {
                    var message = JSON.parse(e.data);
                } catch (exc) {
                    console.log(exc);
                    return;
                }
                var m = message.data;
                if (message.eventName !== 'app.deployed') {
                    $scope.logs.push(m.replace(/(\r\n|\r|\n)/gm, "\n"));
                    $scope.logMessages = $scope.logs.join("\n");
                } else {
                    $scope.deployUrl = m;
                }
            });
            el.addEventListener('error', function (e) {
                $scope.error = "The backend server does not respond correctly or closed the connection.";
                e.currentTarget.close();
            });
        }).error(function(data, status, headers, config) {
            if (data === null || data.error === undefined) {
                $scope.error = "The backend server returned an error: No response was given, come back later.";
            } else {
                $scope.error = "The backend server returned an error: " + (data.error.message || "No response was given, come back later.");
            }
        });
    }
]);
