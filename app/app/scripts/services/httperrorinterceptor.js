'use strict';

/**
 * @ngdoc service
 * @name gitdeployApp.httperrorinterceptor
 * @description
 * # httperrorinterceptor
 * Factory in the gitdeployApp.
 */
angular.module('gitdeployApp')
    .factory('httpErrorInterceptor', [
        '$q', '$rootScope', function ($q, $rootScope) {
            return {
                'responseError': function (rejection) {
                    console.log('Error in response!', rejection);
                    if (rejection.status === 0) {
                        $rootScope.error = {
                            status: 0,
                            message: 'The backend service is unavailable. Either the network is down or there are temporary issues with the backend. Try again later.'
                        };
                    } else if (rejection.data.error !== undefined && rejection.data.error.message !== undefined) {
                        $rootScope.error = {
                            status: rejection.status,
                            message: rejection.data.error.message
                        };
                    } else {
                        $rootScope.error = {
                            status: rejection.status,
                            message: rejection.data
                        };
                    }
                    $rootScope.$broadcast('error');
                    return $q.reject(rejection);
                }
            };
        }
    ]);
