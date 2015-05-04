'use strict';

angular.module('gdApp.version', [
  'gdApp.version.interpolate-filter',
  'gdApp.version.version-directive'
])

.value('version', '0.1');
