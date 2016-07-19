'use strict';

app.factory('IP', ['$http', '$log', function ($http, $log) {

	var ipData = {};

	return {
		fetchIp: function () {
			var url = '//freegeoip.net/json/?';
			return $http.get(url);
		},
		storeIp: function (ipdata) {
			angular.extend(ipData, ipdata, ipdata);
			$log.log('storeIp', ipData);
		},
		getIp: function () {
			return ipData;
		}
	}
}]);