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
app.factory("Buckets", ['$http', 'Config', function ($http, Config) {
	var buckets = [];
	var currentBucket = {};

	function getAll() {
		return $http({
			method: "GET",
			url: Config.API_URL + Config.ENDPOINTS.BUCKETS
		});
	}

	function getOne(id) {
		return $http({
			method: "GET",
			url: Config.API_URL + Config.ENDPOINTS.BUCKETS + "/" + id
		});
	}

	function setBucketAsCurrent(bucket) {
		currentBucket = bucket;
		return currentBucket;
	}

	return {
		buckets: buckets,
		currentBucket: currentBucket,
		setBucketAsCurrent: setBucketAsCurrent,
		getAll: getAll,
		getOne: getOne
	};
}]);

// SNIPPETS FACTORY.
app.factory("Snippets", ['$http', function ($http) {
	var currentSnippet = {};

	function getCurrentSnippet() {
		return currentSnippet;
	}

	function setCurrentSnippetAs(snippet) {
		currentSnippet = snippet;
		return currentSnippet;
	}

	return {
		currentSnippet: currentSnippet,
		setCurrentSnippetAs: setCurrentSnippetAs,
		getCurrentSnippet: getCurrentSnippet
	};
}]);

// WEBSOCKET FACTORY.
// https://github.com/AngularClass/angular-websocket
app.factory("WS", ['$log', 'Config', '$websocket', 'Snippets',
	function ($log, Config, $websocket, Snippets) {
	
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
		return ws.send(snippet);
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

		console.log("received ws message: " + message.data);
		if (message.data !== "Connected.") {
			var snippet = JSON.parse(message.data);
			if (Snippets.currentSnippet.id === snippet.id) {
				Snippets.setCurrentSnippetAs(snippet);
			}
		} else {
			console.log(message.data);
		}
	});

	var methods = {
	  status: status,
	  send: send
	};
	return methods;
}]);

app.factory("Utils", function () {
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
});

app.controller("HackCtrl", ['$scope', 'WS', 'Buckets', 'Snippets', 'Utils', '$timeout',
	function ($scope, WS, Buckets, Snippets, Utils, $timeout) {

	$scope.testes = "this is only a test"

	$scope.data = {};
	$scope.data.msgs = WS.collection;
	$scope.data.ws = WS.status;
	$scope.data.cs = Snippets.getCurrentSnippet();

	Buckets.getAll().then(function (res) {
		$scope.data.buckets = res.data;	
		$scope.data.currentBucket = $scope.data.buckets[0]; // Set current as first for now.
		Buckets.getOne($scope.data.currentBucket.id).then(function (res) {
			
			if (res.data !== null) {
				console.log(JSON.stringify(res));
				// $scope.data.currentBucketSnippets = Buckets.setBucketAsCurrent(res.data);	

				// Snippets.setCurrentSnippetAs(res); // Set current as first for now.
				// $scope.data.cs = Snippets.getCurrentSnippet();	
			} else {
				// nil in bucket
				console.log("nil in bucket")
			}
			
		});
	});

	function setEditorOpts(snip) {
		$scope.editorOptions = {
			mode:  Utils.getLanguageModeByExtension(snip.name || 'boots.go'),
			lineNumbers: true,
			tabSize: 2,
			inputStyle: "contenteditable",
			styleSelectedText: true,
			matchBrackets: true,
			autoCloseBrackets: true,
			showHint: true
		};	
	}
	setEditorOpts($scope.data.cs);
	
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