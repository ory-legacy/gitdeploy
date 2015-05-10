'use strict';

/**
 * @ngdoc service
 * @name gitdeployApp.endpoint
 * @description
 * # endpoint
 * Service in the gitdeployApp.
 */
angular.module('gitdeployApp')
    .service('endpoint', [function () {
        var port = window.location.port === '9000' ? '7654' : window.location.port,
            endpoint = window.location.protocol + '//' + window.location.hostname + ':' + port;

        return {
            sse: endpoint,
            deploy: endpoint,
            apps: endpoint,
            config: endpoint,
            authentication: endpoint
        };
    }]);
