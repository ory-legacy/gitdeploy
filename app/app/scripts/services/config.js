'use strict';

/**
 * @ngdoc service
 * @name gitdeployApp.config
 * @description
 * # config
 * Service in the gitdeployApp.
 */
angular.module('gitdeployApp')
    .service('config', [
        'endpoint', '$http', function (endpoint, $http) {
            return {
                getServerTime: function () {
                    return new Promise(function (resolve, reject) {
                        $http.get(endpoint.apps + '/config',
                            {withCredentials: true}).success(function (data) {
                                resolve(data.now);
                            }).error(reject);
                    });
                },
                getCluster: function () {
                    return new Promise(function (resolve, reject) {
                        $http.get(endpoint.apps + '/config',
                            {withCredentials: true}).success(function (data) {
                                resolve(data.cluster);
                            }).error(reject);
                    });
                },
                get: function () {
                    return new Promise(function (resolve, reject) {
                        $http.get(endpoint.apps + '/config',
                            {withCredentials: true}).success(resolve).error(reject);
                    });
                }
            };
        }
    ]);