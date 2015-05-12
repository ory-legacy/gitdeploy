"use strict";angular.module("gitdeployApp",["ngAnimate","ngCookies","ngResource","ngRoute","ngSanitize","btford.markdown","chart.js","ngTouch"]).config(["$routeProvider","$locationProvider","$httpProvider",function(a,b,c){c.interceptors.push("httpErrorInterceptor"),a.when("/",{templateUrl:"views/landing.html",controller:"LandingCtrl"}).when("/deploy",{templateUrl:"views/deploy.html",controller:"DeployCtrl"}).when("/dashboard/:app",{templateUrl:"views/dashboard.html",controller:"DashboardCtrl"}).when("/contact",{templateUrl:"views/contact.html",controller:"ContactCtrl"}).when("/docs",{templateUrl:"views/docs.html",controller:"DocsCtrl"}).when("/examples",{templateUrl:"views/examples.html",controller:"ExamplesCtrl"}).otherwise({redirectTo:"/"}),b.html5Mode(!0)}]),angular.module("gitdeployApp").controller("LandingCtrl",["$window","$scope",function(a,b){b.repository={},b.showIntegration=function(a){var c,d=document.createElement("a");return a.stopPropagation(),d.href=b.repository.url,b.repository.readmeUrl=d.href.substr(0,d.href.length-4)+"/edit/master/README.md",c=d.pathname.substr(1),c=c.substring(0,c.length-4),c=c.split("/")[1],b.repository.name=c,c=c.replace(/\-/g,"--")+"/master",b.repository.badge="https://img.shields.io/badge/gitdeploy.io-deploy%20"+c+"-green.svg",b.repository.deployUrl=window.location.protocol+"//"+window.location.host+"/deploy?repository="+encodeURIComponent(b.repository.url),b.repository.showBadge=!0,b.repository.deployUrl=window.location.protocol+"//"+window.location.host+"/deploy?repository="+encodeURIComponent(b.repository.url),a.stopPropagation(),$("#shieldModal").modal("show"),!1}}]),angular.module("gitdeployApp").controller("DeployCtrl",["$scope","$routeParams","$http","endpoint",function(a,b,c,d){var e=b.repository,f=function(b){var c=d.sse+"/deployments/"+b+"/events",e=new EventSource(c);return e.addEventListener("open",function(b){a.deploying=!0,console.log("Channel opened!",b)}),{addEventListener:function(b,c){e.addEventListener(b,function(b){a.$apply(function(){c(b)})})}}};return a.logs=[],a.app="",a.deploying=!1,a.retryUrl=window.location.href,a.newsletterMessage="Get a cup of coffee or sign up to our newsletter while you're waiting for the deployment to finish.",void 0===e||e.length<1?void(a.error="The repository query parameter is missing."):(a.error=!1,void c.post(d.deploy+"/deployments",{repository:e}).success(function(b){var c=f(b.data.id);a.app=b.data.id,c.addEventListener("message",function(b){var c;try{c=JSON.parse(b.data)}catch(d){return void console.log(d)}"app.deployed"===c.eventName&&(a.deploying=!1,window.location.href="/dashboard/"+a.app),a.logs.unshift(c.data.replace(/(\r\n|\r|\n)/gm,"\n")),a.logMessages=a.logs.join("\n")},!1),c.addEventListener("error",function(b){a.error="The backend server does not respond correctly or closed the connection.",b.currentTarget.close()})}).error(function(b){null===b||void 0===b.error?a.error="The backend server returned an error: No response was given, come back later.":a.error="The backend server returned an error: "+(b.error.message||"No response was given, come back later.")}))}]),angular.module("gitdeployApp").controller("NavCtrl",["$scope",function(a){a.awesomeThings=["HTML5 Boilerplate","AngularJS","Karma"]}]),angular.module("gitdeployApp").service("endpoint",[function(){var a="9000"===window.location.port?"7654":window.location.port,b=window.location.protocol+"//"+window.location.hostname+":"+a;return{sse:b,deploy:b,apps:b,config:b,authentication:b}}]),angular.module("gitdeployApp").controller("DashboardCtrl",["$scope","$routeParams","apps","config",function(a,b,c,d){var e=b.app;c.getApp(e).then(function(b){var f,g,h=[];angular.forEach(b.data.deployLogs,function(a){try{var b=JSON.parse(a.message);h.unshift(b.data.replace(/(\r\n|\r|\n)/gm,"\n"))}catch(c){console.log(c)}}),b.data.deployLogs=h.join("\n"),f=b.data.ps.split("\n"),b.data.ps=[],b.data.expiresAt=moment(b.data.expiresAt),b.data.createdAt=moment(b.data.createdAt),angular.forEach(f,function(c){c=c.split(/\s+/),b.data.ps.push({id:c[0],type:c[1]}),"web"===c[1]&&(a.noWebProcess=!1,a.$apply())}),b.data.ps.splice(0,1),a.app=b.data,a.$apply(),(g=function(){d.get().then(function(b){a.serverTime=moment(b.data.time);var d=Math.ceil(moment.duration(a.app.expiresAt.diff(a.serverTime)).asMinutes()),f=Math.ceil(moment.duration(a.app.expiresAt.diff(a.app.createdAt)).asMinutes())-d;a.data=[f,d],a.$apply(),a.serverTime.isAfter(a.app.expiresAt)?c.getApp(e).then(function(){}):window.setTimeout(function(){g()},2e4)})})()}),a.labels=["Time used","Time available"],a.data=[0,0],a.colors=["#DCDCDC","#97BBCD"],a.locationHref=window.location.href,a.noWebProcess=!0,a.newsletterMessage="You like Gitdeploy? Sign up to our newsletter and receive updates on new features!"}]),angular.module("gitdeployApp").service("apps",["endpoint","$http",function(a,b){return{getApp:function(c){return new Promise(function(d,e){b.get(a.apps+"/apps/"+c,{withCredentials:!0}).success(d).error(e)})}}}]),angular.module("gitdeployApp").factory("httpErrorInterceptor",["$q","$rootScope",function(a,b){return{responseError:function(c){return console.log("Error in response!",c),0===c.status?b.error={status:0,message:"The backend service is unavailable. Either the network is down or there are temporary issues with the backend. Try again later."}:void 0!==c.data&&void 0!==c.data.error&&void 0!==c.data.error.message?b.error={status:c.status,message:c.data.error.message}:b.error={status:c.status,message:c.data},b.$broadcast("error"),a.reject(c)}}}]),angular.module("gitdeployApp").controller("ErrorCtrl",["$rootScope","$scope",function(a,b){b.$on("error",function(){a.error=b.error}),a.reload=function(){window.location.reload()}}]),angular.module("gitdeployApp").service("config",["endpoint","$http",function(a,b){return{getServerTime:function(){return new Promise(function(c,d){b.get(a.apps+"/config",{withCredentials:!0}).success(function(a){c(a.now)}).error(d)})},getCluster:function(){return new Promise(function(c,d){b.get(a.apps+"/config",{withCredentials:!0}).success(function(a){c(a.cluster)}).error(d)})},get:function(){return new Promise(function(c,d){b.get(a.apps+"/config",{withCredentials:!0}).success(c).error(d)})}}}]),angular.module("gitdeployApp").controller("DocsCtrl",function(){}),angular.module("gitdeployApp").controller("ContactCtrl",["$scope",function(a){a.awesomeThings=["HTML5 Boilerplate","AngularJS","Karma"]}]),angular.module("gitdeployApp").controller("ExamplesCtrl",["$scope",function(a){a.awesomeThings=["HTML5 Boilerplate","AngularJS","Karma"]}]);