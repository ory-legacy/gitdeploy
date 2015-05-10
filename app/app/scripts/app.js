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
        'btford.markdown',
        'chart.js',
        'ngTouch'
    ])
    .config([
        '$routeProvider', '$locationProvider', '$httpProvider', function ($routeProvider, $locationProvider, $httpProvider) {
            $httpProvider.interceptors.push('httpErrorInterceptor');
            $routeProvider.when('/', {
                templateUrl: 'views/landing.html',
                controller: 'LandingCtrl'
            })
                .when('/deploy', {
                    templateUrl: 'views/deploy.html',
                    controller: 'DeployCtrl'
                })
                .when('/dashboard/:app', {
                    templateUrl: 'views/dashboard.html',
                    controller: 'DashboardCtrl'
                })
                .when('/contact', {
                    templateUrl: 'views/contact.html',
                    controller: 'ContactCtrl'
                })
                .when('/docs', {
                    templateUrl: 'views/docs.html',
                    controller: 'DocsCtrl'
                })
                .when('/examples', {
                    templateUrl: 'views/examples.html',
                    controller: 'ExamplesCtrl'
                })
                .otherwise({
                    redirectTo: '/'
                });
            $locationProvider.html5Mode(true);
        }
    ]);

