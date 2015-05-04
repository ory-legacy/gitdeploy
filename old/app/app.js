'use strict';

// Declare app level module which depends on views, and components
angular.module('gdApp', [
    'ngRoute',
    'gdApp.home',
    'gdApp.deploy'
]).config(['$routeProvider', '$interpolateProvider', '$locationProvider', function ($routeProvider, $interpolateProvider, $locationProvider) {
    $routeProvider.otherwise({redirectTo: '/'});
    $interpolateProvider.startSymbol('{(').endSymbol(')}');
    $locationProvider.html5Mode(true)
}]).filter('encodeURIComponent', function() {
    return window.encodeURIComponent;
});

var endpoint = 'http://localhost:7654';