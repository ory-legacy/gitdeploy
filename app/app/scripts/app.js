'use strict';

/**
 * @ngdoc overview
 * @name gitdeployApp
 * @description
 * # gitdeployApp
 *
 * Main module of the application.
 */
angular
    .module('gitdeployApp', [
        'ngAnimate',
        'ngCookies',
        'ngResource',
        'ngRoute',
        'ngSanitize',
        'ngTouch'
    ])
    .config(function ($routeProvider) {
        $routeProvider
            .when('/', {
                templateUrl: 'views/landing.html',
                controller: 'LandingCtrl'
            })
            .when('/deploy', {
                templateUrl: 'views/deploy.html',
                controller: 'DeployCtrl'
            })
            .otherwise({
                redirectTo: '/'
            });
    });
