'use strict';

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
			// console.log('sending ws message:\n ' + JSON.stringify(snippet));
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

