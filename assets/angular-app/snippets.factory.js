'use strict';

// SNIPPETS FACTORY.
app.factory("Snippets", ['$http', '$log', "Config", "Errors",
	function ($http, $log, Config, Errors) {
		// var currentSnippet = {};
		var snippetsLib = {}; // {snippetId: {snippet}, snippetId: {snippet}, snippetId: [{snippet}, {snippet}, ... ]}

		function setOneToSnippetsLib(snippet) {
			if (snippet.id !== "") {
				snippetsLib[snippet.id] = snippet;
			}
		}
		function setManyToSnippetsLib(snippets) {
			// console.log("Got many snippets: " + JSON.stringify(snippets));
			
			for (var i = 0; i < snippets.length; i++) {
				setOneToSnippetsLib(snippets[i]);
			}	
			$log.log('SNIPPETS: ', snippetsLib);
		}
		function getMostRecent(libObj) {
			var timestamps = []; // [timestamps]
			var timesIdLookup = {}; // {timestamp: {snippet}}
			angular.forEach(libObj, function (val, key) {
				// console.log('key: ' + key + ", val: " + val);
				timestamps.push(val.timestamp);
				timesIdLookup[val.timestamp] = val;
			});
			var max = Math.max(...timestamps);
			return timesIdLookup[max];
		}
		function getSnippetsLib() {
			return snippetsLib;
		}
		function getUberAll() {
			return $http({
				method: "GET",
				url: Config.API_URL + Config.ENDPOINTS.SNIPPETS
			});
		}
		function buildNewSnippet() {
			return {
				id: Math.random().toString(36).substring(7),
				// bucketId: "c25pcHBldHM=", // This will be set by the controller pending either the currentBucket (or later any given bucket).
				name: "boots.go",
				language: "go",
				content: "",
				timestamp: Date.now(),
				description: "is a cat",
			};
		}
		function deleteSnippet(snippet) {
			return $http({
				method: "DELETE",
				url: Config.API_URL + Config.ENDPOINTS.SNIPPETS + "/" + snippet.id + "?bucketId=" + snippet.bucketId
			});
		}
		return {
			setOneToSnippetsLib: setOneToSnippetsLib,
			setManyToSnippetsLib: setManyToSnippetsLib,
			getMostRecent: getMostRecent,
			getSnippetsLib: getSnippetsLib,
			getUberAll: getUberAll,
			buildNewSnippet: buildNewSnippet,
			deleteSnippet: deleteSnippet
		};
}]);