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
    .config([
        '$routeProvider', '$locationProvider', '$httpProvider', function ($routeProvider, $locationProvider, $httpProvider) {
            $httpProvider.interceptors.push('httpErrorInterceptor');
            $routeProvider.when('/', {
                templateUrl: 'views/landing.html',
                controller: 'LandingCtrl'
            }).when('/deploy', {
                templateUrl: 'views/deploy.html',
                controller: 'DeployCtrl'
            }).when('/dashboard/:app', {
                templateUrl: 'views/dashboard.html',
                controller: 'DashboardCtrl'
            }).when('/account/create', {
                templateUrl: 'views/accountcreate.html',
                controller: 'AccountCreateCtrl'
            }).when('/connect', {
                templateUrl: 'views/connect.html',
                controller: 'ConnectCtrl'
            }).when('/connect/callback', {
                templateUrl: 'views/githubcallback.html',
                controller: 'GithubCallbackCtrl'
            }).otherwise({
                redirectTo: '/'
            });
            $locationProvider.html5Mode(true);
        }
    ]);

