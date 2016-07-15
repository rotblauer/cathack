'use strict';

var app = angular.module("cathack", [
	'ngWebSocket',
	'ui.codemirror',
	'contenteditable'
	]);

// CONFIG. 
app.constant('Config', {
  // gulp environment: injects environment vars
  "WS_URL": "ws://" + window.location.host + "/hack/ws",
  "API_URL": "http://" + window.location.host + "/hack",
  "ENDPOINTS": {
  	"BUCKETS": "/b",
  	"SNIPPETS": "/s",
  	"FS": "/fs"
  },
  "EDITOROPTIONS": {
		mode: {name: 'markdown'},
		lineNumbers: true,
		tabSize: 2,
		inputStyle: "contenteditable",
		styleSelectedText: true,
		matchBrackets: true,
		autoCloseBrackets: true,
		showHint: true,
		allowDropFileTypes: [
			'text/plain',
			'text/html',
			'text/javascript',
			'application/javascript',
			'text/css'
		]
	},
	"DEFAULTSNIPPET": {
		id: Math.random().toString(36).substring(7),
		// bucketId: "c25pcHBldHM=", // This will be set by the controller pending either the currentBucket (or later any given bucket).
		name: "boots.go",
		language: "go",
		content: "",
		timestamp: Date.now(),
		description: "is a cat",
	}, 
	"DEFAULTBUCKET": {
		id: Math.random().toString(36).substring(7),
		meta: {
			name: "snippets"
		}
	}
});

app.factory("Utils", ["Config", function (Config) {

	function typeOf (obj) {
	  return {}.toString.call(obj).split(' ')[1].slice(0, -1).toLowerCase();
	}

	function setEditorOptions(obj) {
		return angular.extend({}, Config.EDITOROPTIONS, obj);
	}
	
	function getLanguageModeByExtension(name) {
	  var o = "";
	  var exs = name.split(".");
	  console.log('exs -> ' + JSON.stringify(exs));
	  for (var i = 0; i < exs.length ; i++) {
	    var ex = exs[i];
	    console.log("checking case: " + ex);
	    switch (ex) {
	      case 'js':
	        o = javascriptMode;
	        break;
	      case 'erb': 
	      case 'html':
	        o = htmlmixedMode;
	        break;
	      case 'go':
	        o = goMode;
	        break;
	      case 'md':
	      case 'mdown':
	      case 'markdown':
	        o = markdownMode;
	        break;
	      case 'py':
	        o = pythonMode;
	        break;
	      case 'rb':
	        o = rubyMode;
	        break;
	      case 'r':
	        o = rMode;
	        break;
	      case 'sh':
	      case 'zsh':
	        o = shellMode;
	        break;
	      case 'swift':
	        o = swiftMode;
	        break;
	      default: 
	        o = markdownMode;
	        break;
	    }  
	  }
	  console.log("language by extension -> " + o);
	  return o;
	}

	var goMode = {name: 'go'};
	var javascriptMode = {name: 'javascript'};
	var markdownMode = {name: 'markdown'};
	var pythonMode = {name: 'python'};
	var rubyMode = {name: 'ruby'};
	var rMode = {name: 'r'};
	var shellMode = {name: 'shell'};
	var swiftMode = {name: 'swift'};
	var htmlmixedMode = {name: 'htmlmixed', scriptTypes: [{matches: /\/x-handlebars-template|\/x-mustache/i,
	                 mode: null},
	                {matches: /(text|application)\/(x-)?vb(a|script)/i,
	                 mode: "vbscript"}]};
	var htmlembeddedMode = {name: 'htmlembedded'};

	return {
		typeOf: typeOf,
		setEditorOptions: setEditorOptions,
		getLanguageModeByExtension: getLanguageModeByExtension,
		goMode: goMode,
		javascriptMode: javascriptMode,
		markdownMode: markdownMode,
		pythonMode: pythonMode,
		rubyMode: rubyMode,
		rMode: rMode,
		shellMode: shellMode,
		swiftMode: swiftMode,
		htmlmixedMode: htmlmixedMode,
		htmlembeddedMode: htmlembeddedMode
	};
}]);


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

// SNIPPETS FACTORY.
app.factory("Snippets", ['$http', "Config", "Errors",
	function ($http, Config, Errors) {
		// var currentSnippet = {};
		var snippetsLib = {}; // {snippetId: {snippet}, snippetId: {snippet}, snippetId: [{snippet}, {snippet}, ... ]}

		function setOneToSnippetsLib(snippet) {
			if (snippet.id !== "") {
				snippetsLib[snippet.id] = snippet;
			}
		}
		function setManyToSnippetsLib(snippets) {
			console.log("Got many snippets: " + JSON.stringify(snippets));
			
			for (var i = 0; i < snippets.length; i++) {
				setOneToSnippetsLib(snippets[i]);
			}	
			console.log(JSON.stringify(snippetsLib));
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

app.factory("FS", ['$log', '$http', 'Config', function ($log, $http, Config) {

	function fetchFS() {
		var url = Config.API_URL + Config.ENDPOINTS.FS
		return $http.get(url);
	}

	function importDir(path) {
		var url = Config.API_URL + 
						  Config.ENDPOINTS.FS + 
						  Config.ENDPOINTS.BUCKETS
		var config = {
			params: {
				path: JSON.stringify(path)
			}
		}
		return $http.get(url, config); // c.JSON(200, gin.H{"b": bs, "s": ss})
	}

	function importFile(path) {
		var url = Config.API_URL + 
								  Config.ENDPOINTS.FS + 
								  Config.ENDPOINTS.SNIPPETS
		var config = {
			params: {
				path: JSON.stringify(path)
			}
		}
		return $http.get(url, config); // c.JSON(200, gin.H{"b": b, "s": s})
	}

	return {
		fetchFS: fetchFS,
		importFile: importFile,
		importDir: importDir
	};
}]);

// WEBSOCKET FACTORY.
// https://github.com/AngularClass/angular-websocket
app.factory("WS", ['$log', 'Config', '$websocket', 'Snippets', 'Utils',
	function ($log, Config, $websocket, Snippets, Utils) {
	
		var ws = $websocket(Config.WS_URL);
		var status = {
			available: false,
			error: ""
		}
		function setStatus (opts) {
			angular.extend(status, status, opts);
		}
		function getStatus() {
			return status;
		}
		function send(snippet) {
			console.log('sending ws message:\n ' + JSON.stringify(snippet));
			return ws.send(JSON.stringify(snippet));
		}

		ws.onOpen(function() {
			setStatus({available: true});
		});
		ws.onClose(function() {
			setStatus({available: false});
		});
		ws.onError(function(err) {
			setStatus({error: err});
		});

		var methods = {
		  status: status,
		  send: send,
		  ws: ws
		};
		return methods;
}]);

app.factory("Errors", [function () {
	var error = {};
	function getError() {
		return error;
	}
	function setError(err) {
		error = err;
		return err;
	}
	return {
		getError: getError,
		setError: setError
	};
}]);

app.controller("HackCtrl", ['$scope', 'WS', 'Buckets', 'Snippets', 'FS', 'Utils', '$timeout', 'Errors', 'Config', '$log',
	function ($scope, WS, Buckets, Snippets, FS, Utils, $timeout, Errors, Config, $log) {

	$scope.testes = "this is only a test"

	$scope.data = {};
	$scope.data.ws = WS.status;

	$scope.data.cs = Config.DEFAULTSNIPPET; // inits as {}
	$scope.data.cb = Config.DEFAULTBUCKET; // empty

	$scope.data.buckets = Buckets.getBuckets();
	$scope.data.snippets = Snippets.getSnippetsLib();
	
	$scope.data.error = Errors.getError();

	$scope.editorOptions = Utils.setEditorOptions({
		mode: Utils.getLanguageModeByExtension($scope.data.cs.language)
	});


	// Init.
	// 
	// Fetch all snippets [{id: "12341=", name: "boots.go" ... }]
	Snippets.getUberAll().then(function (res) {
		console.log("got snippets ->" + JSON.stringify(res.data));
		// if there are any snippets at all
		if (Utils.typeOf(res.data) !== 'null') {
			Snippets.setManyToSnippetsLib(res.data);
			$scope.data.cs = Snippets.getMostRecent(Snippets.getSnippetsLib()); 
		} 
		// return res;
	})
	.then(function () {
		return Buckets.fetchAll();
	})
	.then(function (res) {
		console.log('got buckets -> ' + JSON.stringify(res.data));
		Buckets.storeManyBuckets(res.data);

		// if lib === {}
		if (Object.keys(Snippets.getSnippetsLib()).length === 0 && Snippets.getSnippetsLib().constructor === Object) {
			$scope.data.cb = res.data[0]; // presuming bucket is snippets?
			$scope.data.cs.bucketId = $scope.data.cb.id; // set default snippet to have default bucket id
		} else {
			$scope.data.cb = Buckets.getBuckets()[$scope.data.cs.bucketId]; // set current bucket to be current snippet's bucket
		}
	})
	.then(function () {
		return FS.fetchFS();
	})
	.then(function (res) {
		$log.log('Got FS', res.data);
		$scope.data.cfs = res.data;
	})
	.catch(Errors.setError);


	// Watch for changes to current snippet and keep the library updated. 
	$scope.$watchCollection('data.cs', function (old, neww) {
		if (old !== neww) {
			Snippets.setOneToSnippetsLib($scope.data.cs);
			//check to set language
			$scope.data.cs.language = Utils.getLanguageModeByExtension($scope.data.cs.name);

			$scope.editorOptions = Utils.setEditorOptions({
				mode: Utils.getLanguageModeByExtension(neww.name)
			});
		}
		
	}, true);


	function sendUpdate() {
		if (Utils.typeOf($scope.data.cs.bucketId) === 'undefined') {
			$scope.data.cs.bucketId = Buckets.getCurrentBucket().id // make sure we're sending a snip with a  bucketId
		}
		WS.send($scope.data.cs);
		$scope.data.cs.timestamp = Date.now();
	}

	WS.ws.onMessage(function(message) {

		console.log("received ws message:\n " + message.data);
		var o = JSON.parse(message.data);

		if (Utils.typeOf(o) === 'object') {
			console.log("received msg was obj");
			if (o.bucketId !== 'undefined') { // check if is snippet
				console.log("received msg was snippet");
				if ($scope.data.cs.id === o.id) {
					console.log("is current snippet!");
					$scope.data.cs = o; // set current snippet to equal incoming
				} else {
					console.log("not current snippet. setting to lib.");
					Snippets.setOneToSnippetsLib(o); // so we have a fresh version in store
				}
			} else if (o.meta !== 'undefined') { // check if is a bucket
				console.log("shit. msg was undefined.");
			}
		} else if (Utils.typeOf(o) === 'array') {
			console.log('what the ws received an array');
		} else {
			console.log('what the fuck the ws received some shit:\n' + JSON.stringify(o));
		}
	});

	// Manual listeners. Avoid infdiggers. 
	document.getElementById("editor").addEventListener('keyup', function (e) {
	  // http://stackoverflow.com/questions/2257070/detect-numbers-or-letters-with-jquery-javascript
	  var inp = String.fromCharCode(e.keyCode);

	  if (e.metaKey || e.ctrlKey || e.altKey) { //  || /[\[cC]{1}/.test(inp) // preventing copy to clipboard?
	  	return false;
	  }

	  // if (/[VCA]/.test(inp)) {
	  // 	return false;
	  // }

  	var inp = String.fromCharCode(e.keyCode);

  	$log.log(inp);
  	
  	if (/\S/.test(inp) || e.which === 13 || e.keyCode === 8 || e.keyCode ===  9) { // 224
  	  sendUpdate(); // updateCurrentSnippetFromGUI returns currentSnippet {} var  
  	}	
	  
	});
	document.getElementById("snippetName").addEventListener('keyup', function (e) {
		
		if ($scope.data.cs.name == "") {
			$scope.data.cs.name = "_";
		}
		sendUpdate();
	});
	document.getElementById("snippetDescription").addEventListener('keyup', function (e) {
		$scope.data.cs = $scope.data.cs;
		sendUpdate();
	});


	// Create new snippet. 
	$scope.createNewSnippet = function () {
		$scope.data.cs = Snippets.buildNewSnippet();
		$scope.data.cs.bucketId = $scope.data.cb.id;
		document.getElementById("snippetName").focus();
	};	

	$scope.deleteSnippet = function (snippet) {
		var ok = window.confirm("OK to delete this snippet?")
		if (ok) {
			return Snippets.deleteSnippet(snippet)
				// c.JSON(200, gin.H{
				// 	"snippetId": snippetId,
				// 	"bucketId":  bucketId,
				// 	"snippets":  snippets,
				// })
				.then(function (res) {
					// Snippets.setManyToSnippetsLib(res.data.snippets)
					delete $scope.data.snippets[res.data.snippetId];
					$scope.data.cs = Snippets.getMostRecent(Snippets.getSnippetsLib());
				})
				.then(function (res) {
					console.log('deleted snippet');
				})
				.catch(function (err) {
					console.log("failed to delete snippet." + JSON.stringify(err.data));
				});	
		}			
	};


	$scope.selectSnippetAsCurrent = function (snippet) {
		$scope.data.cs = snippet;
	};

	$scope.createNewBucket = function() {
		Buckets.createBucket("newbucket")
		.then(function (res) {
			// TODO: websocketize. 
			Buckets.storeOneBucket(res.data);
			$scope.data.buckets = Buckets.getBuckets();
			$scope.data.cb = res.data;
		})
		.catch(function (err) {
			$log.error(err);
		})
	};

	$scope.deleteBucket = function (bucket) {
		var ok = window.confirm("Really really? This will delete '" + bucket.meta.name + "' and all its snippets. Forever.\nOK?");
		if (ok) {
			Buckets.destroyBucket(bucket)
				.then(function (res) {
					delete Buckets.getBuckets()[bucket.id];
					$scope.data.cb = Buckets.getMostRecent(Buckets.getBuckets());
					// clean snippetsLib
					angular.forEach(Snippets.getSnippetsLib(), function (val, key) { 
						if (val.bucketId === bucket.id) {
							delete this[key];
						}
					}, Snippets.getSnippetsLib());

					if ($scope.data.cs.bucketId === bucket.id) { // but don't switch cs if it wasn't effected 
						// then set to most recent existing snippet
						$scope.data.cs = Snippets.getMostRecent(Snippets.getSnippetsLib());
					}
				})
				.catch(function (err) {
					$log.log(err);
				});	
		} else {
			// nada.
		}
		
	};

	$scope.selectBucketAsCurrent = function (bucket) {
		$scope.data.cb = bucket;
	};

	$scope.saveBucket = function (bucket) {
		Buckets.putBucket(bucket)
			.success(function (res) {
				$log.log(res);
			})
			.error(function (err) {
				$log.log(err);
			});
	};

	$scope.import = function (path, isDir) {
		if (isDir) {
			FS.importDir(path)
				.success(function (res) {
					$log.log('got res:', res);
					Buckets.storeManyBuckets(res.b);
					Snippets.setManyToSnippetsLib(res.s);
					$scope.data.snippets = Snippets.getSnippetsLib();
					$scope.data.cs = Snippets.getMostRecent(Snippets.getSnippetsLib());
					$scope.data.cb = Buckets.getBuckets()[$scope.data.cs.bucketId];
				})
				.error(function (err) {
					$log.log('importDir failed', err);
				});
		} else {
			FS.importFile(path)
				.success(function (res) {
					$log.log('got res:', res);
					Buckets.storeOneBucket(res.b);
					$scope.data.cb = Buckets.getBuckets()[res.b.id];
					Snippets.setOneToSnippetsLib(res.s);
					$scope.data.snippets = Snippets.getSnippetsLib();
					$scope.data.cs = Snippets.getSnippetsLib()[res.s.id];
				})
				.error(function (err) {
					$log.log("importFile failed", err);
				});
		}
	};

}]);


