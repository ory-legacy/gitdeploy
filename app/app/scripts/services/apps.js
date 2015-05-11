'use strict';

/**
 * @ngdoc service
 * @name gitdeployApp.apps
 * @description
 * # apps
 * Service in the gitdeployApp.
 */
angular.module('gitdeployApp')
    .service('apps', [
        'endpoint', '$http', function (endpoint, $http) {
            return {
                getApp: function (id) {
                    return new Promise(function (resolve, reject) {
                        $http.get(endpoint.apps + '/apps/' + id,
                            {withCredentials: true}).success(resolve).error(reject);
                    });
                }
            };
        }
    ]);
