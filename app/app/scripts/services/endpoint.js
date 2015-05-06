'use strict';

/**
 * @ngdoc service
 * @name gitdeployApp.endpoint
 * @description
 * # endpoint
 * Service in the gitdeployApp.
 */
angular.module('gitdeployApp')
    .service('endpoint', function () {
        var endpoint, development = window.location.hostname === 'localhost';
        if (development) {
            endpoint = 'http://' + window.location.hostname + ':7654';
        } else {
            endpoint = 'https://api.' + window.location.hostname;
        }

        return {
            sse: endpoint,
            deploy: endpoint,
            apps: endpoint,
            authentication: endpoint
        };
    });
