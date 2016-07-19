'use strict';

// BUCKETS FACTORY.
app.factory("Buckets", ['$http', 'Config', "Errors", "Snippets", 'Utils', function ($http, Config, Errors, Snippets, Utils) {
	var buckets = {}; // {bucketId: {id: "234234932-=", meta: {name: "snippets", timestamp: 1232354234}, }
	// var currentBucket = {};

	function getBuckets() {
		return buckets;
	}
	function storeOneBucket(bucket) {
		console.log('storing one bucket: ' + JSON.stringify(bucket))
		buckets[bucket.id] = bucket;
	}
	function storeManyBuckets(buckets) {
		console.log('storeing bucktses')
		if (Utils.typeOf(buckets) === 'object') {
			console.log('bucket is object');
		} else {
			console.log('bucket is not object');
			for (var i = 0; i < buckets.length; i++) {
				console.log('bucket i ' + buckets[i]);
				storeOneBucket(buckets[i]);
			}	
		}
		console.log('BUCKETS: ' + JSON.stringify(buckets));
	}
	function getMostRecent(libObj) {
		var timestamps = []; // [timestamps]
		var timesIdLookup = {}; // {timestamp: {snippet}}
		angular.forEach(libObj, function (val, key) {
			// console.log('key: ' + key + ", val: " + val);
			timestamps.push(val.meta.timestamp);
			timesIdLookup[val.meta.timestamp] = val;
		});
		var max = Math.max(...timestamps);
		return timesIdLookup[max];
	}
	function fetchAll() {
		return $http({
			method: "GET",
			url: Config.API_URL + Config.ENDPOINTS.BUCKETS
		});
	}
	function destroyBucket(bucket) {
		return $http({
			method: "DELETE",
			url: Config.API_URL + Config.ENDPOINTS.BUCKETS + "/" + bucket.id
		});
	}
	function createBucket(bucketName) {
		return $http({
			method: "POST",
			url: Config.API_URL + Config.ENDPOINTS.BUCKETS + "/" + bucketName
		});
	}
	function putBucket(bucket) {
		var url = Config.API_URL + Config.ENDPOINTS.BUCKETS + "/" + bucket.id;
		var param = JSON.stringify(bucket);
		return $http.put(url, param);
		// return $http({
		// 	method: "PUT",
		// 	url: Config.API_URL + Config.ENDPOINTS.BUCKETS + "/" + bucket.id + "?to=",
		// 	headers: {

		// 	}
		// });
	}
	return {
		storeOneBucket: storeOneBucket,
		storeManyBuckets: storeManyBuckets,
		fetchAll: fetchAll,
		getBuckets: getBuckets,
		getMostRecent: getMostRecent,
		createBucket: createBucket,
		destroyBucket: destroyBucket,
		putBucket: putBucket
	};
}]);