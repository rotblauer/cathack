'use strict';

var app = angular.module("cathack", [
	'ngWebSocket',
	'ui.codemirror'
	]);

// CONFIG. 
app.constant('Config', {
  // gulp environment: injects environment vars
  "WS_URL": "ws://" + window.location.host + "/hack/ws",
  "API_URL": "http://" + window.location.host + "/hack",
  "ENDPOINTS": {
  	"BUCKETS": "/b",
  	"SNIPPETS": "/s"
  },
  "EDITOROPTIONS": {
		mode: {name: 'markdown'},
		lineNumbers: true,
		tabSize: 2,
		inputStyle: "contenteditable",
		styleSelectedText: true,
		matchBrackets: true,
		autoCloseBrackets: true,
		showHint: true
	},
	"DEFAULTSNIPPET": {
		id: Math.random().toString(36).substring(7),
		bucketId: "c25pcHBldHM=", // This will be set by the controller pending either the currentBucket (or later any given bucket).
		name: "boots.go",
		language: "go",
		content: "well hello",
		timestamp: Date.now(),
		description: "is a cat",
	}
});

// app.config(function ($stateProvider, $urlRouterProvider) {

//   // ROUTING with ui.router
//   $urlRouterProvider.otherwise('/');
//   $stateProvider
//     // this state is placed in the <ion-nav-view> in the index.html
//     .state('buckets', {
//       url: '/',
//       abstract: true,
//       templateUrl: 'main/templates/menu.html',
//       controller: 'MenuCtrl as menu'
//     })
//       .state('main.home', {
//         url: '/home',
//         views: {
//           'main-home': {
//             templateUrl: 'main/templates/home.html',
//             controller: 'HomeCtrl'
//           }
//         }
//       })
//       .state('main.listDetail', {
//         url: '/list/:poopId',
//         views: {
//           'main-list': {
//             templateUrl: 'main/templates/list-detail.html',
//             controller: 'ListDetailCtrl',
//             resolve: {
//               poop: function ($stateParams, Tallyrally) {
//                 return Tallyrally.get({id: $stateParams.poopId}, function (res) {
//                   return res;
//                 });
//               }
//               // poopId: function ($stateParams) {
//               //   return $stateParams.poopId;
//               // }
//             }
//           }
//         }
//       });
// });

// BUCKETS FACTORY.
app.factory("Buckets", ['$http', 'Config', "Errors", "Snippets", 'Utils', function ($http, Config, Errors, Snippets, Utils) {
	var buckets = {}; // {bucketId: {id: "234234932-=", meta: {name: "snippets", timestamp: 1232354234}, }
	var currentBucket = {};

	function getBuckets() {
		return buckets;
	}
	function storeBucket(bucket) {
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
				storeBucket(buckets[i]);
			}	
		}
		console.log('BUCKETS: ' + JSON.stringify(buckets));
	}
	function getCurrentBucket() {
		return currentBucket;
	}
	function fetchAll() {
		return $http({
			method: "GET",
			url: Config.API_URL + Config.ENDPOINTS.BUCKETS
			// error: Errors.setError
			// success: function (res) {
			// 	var o = JSON.parse(res.data); // []
			// 	for (i = 0; i < o.length; i++) {
			// 		storeBucket(o[i]);
			// 	}
			// }
		});
	}
	function fetchSnippetsFor(id) {
		return $http({
			method: "GET",
			url: Config.API_URL + Config.ENDPOINTS.BUCKETS + "/" + id
			// error: Errors.setError,
			// success: function (res) {
				
			// }
		});	
	}
	function setCurrentBucketAs(bucket) {
		currentBucket = bucket;
		return currentBucket;
	}
	return {
		storeBucket: storeBucket,
		storeManyBuckets: storeManyBuckets,
		buckets: buckets,
		// currentBucket: currentBucket,
		setCurrentBucketAs: setCurrentBucketAs,
		getCurrentBucket: getCurrentBucket,
		fetchAll: fetchAll,
		getBuckets: getBuckets
	};
}]);

// SNIPPETS FACTORY.
app.factory("Snippets", ['$http', "Config", "Errors",
	function ($http, Config, Errors) {
		var currentSnippet = {};
		var snippetsLib = {}; // {snippetId: [{snippet}, {snippet}], snippetId: [{snippet}, {snippet}, ... ]}

		function setOneToSnippetsLib(snippet) {
			if (snippet.id !== "") {
				snippetsLib[snippet.id] = snippet;
			}
			console.log('SNIPPETSLIB: ' + JSON.stringify(snippetsLib))
		}
		function setManyToSnippetsLib(snippets) {
			console.log("Got many snippets: " + JSON.stringify(snippets));
			
			for (var i = 0; i < snippets.length; i++) {
				setOneToSnippetsLib(snippets[i]);
			}	
			console.log(JSON.stringify(snippetsLib));
		}
		function getSnippetsLib() {
			return snippetsLib;
		}
		function getDefaultSnippet(bucketId) {
			return angular.extend({}, Config.DEFAULTSNIPPET, {bucketId: bucketId});
		}
		function getCurrentSnippet() {
			return currentSnippet;
		}
		function setCurrentSnippetAs(snippet) {
			currentSnippet = snippet;
			return currentSnippet;
		}
		function getUberAll() {
			return $http({
				method: "GET",
				url: Config.API_URL + Config.ENDPOINTS.SNIPPETS
			});
		}
		return {
			setOneToSnippetsLib: setOneToSnippetsLib,
			setManyToSnippetsLib: setManyToSnippetsLib,
			// currentSnippet: currentSnippet,
			setCurrentSnippetAs: setCurrentSnippetAs,
			getCurrentSnippet: getCurrentSnippet,
			getSnippetsLib: getSnippetsLib,
			getUberAll: getUberAll
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
			console.log('sending ws message');
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

		ws.onMessage(function(message) {

			console.log("received ws message:\n " + message.data);
			var o = JSON.parse(message.data);

			if (Utils.typeOf(o) === 'object') {
				console.log("received msg was obj");
				if (o.bucketId !== 'undefined') { // check if is snippet
					console.log("received msg was snippet");
					if (Snippets.getCurrentSnippet().id === o.id) {
						console.log("is current snippet!");
						Snippets.setCurrentSnippetAs(o);
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

		var methods = {
		  status: status,
		  send: send
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

app.controller("HackCtrl", ['$scope', 'WS', 'Buckets', 'Snippets', 'Utils', '$timeout', 'Errors', 'Config',
	function ($scope, WS, Buckets, Snippets, Utils, $timeout, Errors, Config) {

	$scope.testes = "this is only a test"

	$scope.data = {};
	$scope.data.ws = WS.status;

	$scope.data.cs = Config.DEFAULTSNIPPET; // inits as {}
	$scope.data.cb = Buckets.getCurrentBucket();

	$scope.data.buckets = Buckets.getBuckets();
	$scope.data.snippets = Snippets.getSnippetsLib();
	$scope.data.error = Errors.getError();

	$scope.editorOptions = Utils.setEditorOptions({
		mode: Utils.getLanguageModeByExtension($scope.data.cs.language)
	});

	Snippets.getUberAll().then(function (res) {
		console.log(res);
		Snippets.setManyToSnippetsLib(res.data);
	}).then(function (res) {
		if (Snippets.getSnippetsLib.length > 0) {
			$scope.data.cs = Snippets.getCurrentSnippet[0];
		}
	}).catch(function (err) {
		console.log('errr!');
		Errors.setError(err);
	});

	Buckets.fetchAll().then(function (res) {
		console.log(res.data);
		Buckets.storeManyBuckets(res.data);
		// $scope.data.buckets = Buckets.getBuckets();
	});

	$scope.$watchCollection('data.cs', function (old, neww) {
		console.log("data.cs changed:\n old: " + JSON.stringify(old) + "\n new: " + JSON.stringify(neww));
		WS.send(neww);
	}, true);

	// Snippets.getUberAll().then(function (res) {
	// 	console.log("Got snippets uberall ctrl: " + JSON.stringify(res.data));
	// 	$scope.data.snippets = Snippets.getSnippetsLib;
	// 	// $scope.data.snippets = Snippets.getSnippetsLib;
	// 	// Snippets.setManyToSnippetsLib(res.data);
	// 	// $scope.data.snippets = Snippets.getSnippetsLib;
	// });

	// Buckets.fetchAll().then(function (res) {
	// 	// guaranteed at least one bucket
	// 	// get snippets for arbitrarily first bucket
	// 	Snip.fetchSnippetsFor(res.data[0].id);
	// });

	// Buckets.fetchAll().then(function (res) {
	// 	$scope.data.buckets = res.data;	
	// 	$scope.data.currentBucket = $scope.data.buckets[0]; // Set current as first for now.
	// 	Buckets.getOne($scope.data.currentBucket.id).then(function (res) {
			
	// 		if (res.data !== null) {
	// 			console.log(JSON.stringify(res));
	// 			// $scope.data.currentBucketSnippets = Buckets.setCurrentBucketAs(res.data);	

	// 			// Snippets.setCurrentSnippetAs(res); // Set current as first for now.
	// 			// $scope.data.cs = Snippets.getCurrentSnippet();	
	// 		} else {
	// 			// nil in bucket
	// 			console.log("nil in bucket")
	// 		}
			
	// 	});
	// });

	// function setEditorOpts(snip) {

	// }
	// setEditorOptions($scope.data.cs);
	

	// Utils.getLanguageModeByExtension(snip.name || 'boots.go'),

	// function setEditorOpts(snip) {
	// 	$scope.editorOptions = {
	// 		mode:  
	// 		lineNumbers: true,
	// 		tabSize: 2,
	// 		inputStyle: "contenteditable",
	// 		styleSelectedText: true,
	// 		matchBrackets: true,
	// 		autoCloseBrackets: true,
	// 		showHint: true
	// 	};	
	// }
	// setEditorOpts($scope.data.cs);
	
	// $scope.chooseSnippet = function(snippet) {
	// 	$scope.data.cs = Snippets.setCurrentSnippetAs(snippet);
	// };

	// Now we can just watch the whole current snippet in one listener! 
	// Pretty nifty!
	
		// $scope.$watch('data.cs', function (old, neww, scope) {
		// 	console.log("o: " +JSON.stringify(old) + "\n" + "n: " + JSON.stringify(neww));
			
		// 	if (neww !== old) {
		// 		console.log('stuff changed!');
		// 		var sendMe = {};
		// 		sendMe.id = neww.id;
		// 		// Update timestamp.
		// 		sendMe.timestamp = Date.now();
		// 		$scope.data.cs.timestamp = sendMe.timestamp
		// 		sendMe.name = neww.name;
		// 		sendMe.meta = neww.meta;
		// 		sendMe.content = neww.content;
		// 		// Check for language change.
		// 		if (old.language !== neww.language) {
		// 			sendMe.language = Utils.getLanguageModeByExtension(sendMe.name || 'boots.go');
		// 			setEditorOpts(sendMe.language);	
		// 		}	
		// 		WS.send(JSON.stringify(sendMe));
		// 	}
		// }, true);	
	
	

}]);

// http://fdietz.github.io/recipes-with-angular-js/common-user-interface-patterns/editing-text-in-place-using-html5-content-editable.html
app.directive("contenteditable", function() {
  return {
    restrict: "A",
    require: "ngModel",
    link: function(scope, element, attrs, ngModel) {

      function read() {
        ngModel.$setViewValue(element.html());
      }

      ngModel.$render = function() {
        element.html(ngModel.$viewValue || "");
      };

      element.bind("blur keyup change", function() {
        scope.$apply(read);
      });
    }
  };
});