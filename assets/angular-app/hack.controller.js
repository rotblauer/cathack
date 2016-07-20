'use strict';

app.controller("HackCtrl", ['$scope', '$location', 'WS', 'IP', 'Buckets', 'Snippets', 'FS', 'Utils', '$timeout', 'Errors', 'Config', '$log', 'flash',
	function ($scope, $location, WS, IP, Buckets, Snippets, FS, Utils, $timeout, Errors, Config, $log, flash) {

	$scope.testes = "this is only a test"

	$scope.data = {};
	$scope.data.ws = WS.status;

	$scope.data.cs = Config.DEFAULTSNIPPET; // inits as {}
	$scope.data.cb = Config.DEFAULTBUCKET; // empty

	$scope.data.buckets = Buckets.getBuckets();
	$scope.data.snippets = Snippets.getSnippetsLib();
	$scope.data.cfs = FS.getFS();
	
	$scope.data.ip = IP.getIp();

	$scope.data.error = Errors.getError();

	$scope.flash = flash;

	$scope.data.editorOptions = Utils.setEditorOptions({
		mode: Utils.getLanguageModeByExtension($scope.data.cs.language)
	});

	function flashAlert(text, classs) {
		flash.setMessage(text, classs);
		$timeout(function() {
			flash.setMessage({});
		}, 3000);
	}

	function setEditOpts(snippet) {
		$scope.data.editorOptions = Utils.setEditorOptions({
			mode: Utils.getLanguageModeByExtension(snippet.name)
		});
	}

	function sendUpdate() {
		if (Utils.typeOf($scope.data.cs.bucketId) === 'undefined') {
			$scope.data.cs.bucketId = Buckets.getCurrentBucket().id // make sure we're sending a snip with a  bucketId
		}
		$scope.data.cs.ipCity = $scope.data.ip['city'];
		$scope.data.cs.ip = $scope.data.ip['ip'];
		// $log.log('sendup snip:', $scope.data.cs);
		WS.send($scope.data.cs);
		$scope.data.cs.timestamp = Date.now();
	}

	IP.fetchIp()
		.then(function (res) {
			$log.log('got ip', res);
			IP.storeIp(res.data);
		})
		.catch(function (err) {
			$log.log('fail getting ip', err);
		});

	// TODO uirouter
	// $scope.data.cb = bucket;
	// $log.log("$scope.data.cb: ", $scope.data.cb);
	// $scope.data.cs = snippet;
	// $log.log("$scope.data.cs: ", $scope.data.cs);
	// if ($scope.data.cs.bucketId === null) {
	// 	$scope.data.cs.bucketId = bucket.id;
	// }
	// $scope.data.cfs = FS.getFS();
	// 
	// var hashSnippetId = $location.hash();
	// $log.log('hashSnippetId: ', hashSnippetId);
	// var winpath = $location.path();
	// $log.log('$locations.. ', $location.path(), $location.hash(), $location.url());
	
	function setSnippetParam(snippet) {
		return $location.hash(snippet.id);
	}
	function getSnippetParam() {
		return $location.hash().replace('/','');
	}

	// Init.
	// 
	// Fetch all snippets [{id: "12341=", name: "boots.go" ... }]
	Snippets.getUberAll().then(function (res) {
		
		// if there are any snippets at all
		if (Utils.typeOf(res.data) !== 'null') {
			Snippets.setManyToSnippetsLib(res.data);
			if (getSnippetParam().length > 0) {
				
				// Get snippet by hash param. 
				var s = Snippets.getSnippetsLib()[getSnippetParam()];
				// If snippet exists, great. Make it current. 
				if (s) {
					$scope.data.cs = s	
				
				// ... redirect in case snippet has been deleted. 
				} else {
					$scope.data.cs = Snippets.getMostRecent(Snippets.getSnippetsLib());
					// setSnippetParam($scope.data.cs);
				}
				
			} else {
				// There is no snippet Id in the params. Get most recent. 
				$scope.data.cs = Snippets.getMostRecent(Snippets.getSnippetsLib()); 	
				// ... then set the hash param. 
				// setSnippetParam($scope.data.cs);
			}
		} 
	})
	.then(function () {
		return Buckets.fetchAll();
	})
	.then(function (res) {
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
		FS.storeFS(res.data);
		
	})
	.catch(Errors.setError);


	// Watch for changes to current snippet and keep the library updated. 
	$scope.$watchCollection('data.cs', function (old, neww) {
		if (old !== neww) {
			Snippets.setOneToSnippetsLib($scope.data.cs);
			//check to set language
			$scope.data.cs.language = Utils.getLanguageModeByExtension($scope.data.cs.name);

			setEditOpts(neww);
			setSnippetParam(neww);
		}
		
	}, true);


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
	// document.getElementById("editor").addEventListener('change', function (from, to, text, removed, origin) {
	// 	sendUpdate();
	// });
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

  	// $log.log(inp);
  	
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
		// setSnippetParam($scope.data.cs);
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
					flashAlert("Successfully deleted snippet.", 'info');
				})
				.catch(function (err) {
					console.log("failed to delete snippet." + JSON.stringify(err.data));
					flashAlert("Failed to delete snippet.", 'danger');
				});	
		}			
	};

	$scope.selectSnippetAsCurrent = function (snippet) {
		$scope.data.cs = snippet;
		$scope.selectBucketAsCurrent(Buckets.getBuckets()[$scope.data.cs.bucketId]);
		// setSnippetParam(snippet);
		setEditOpts(snippet);
	};

	$scope.createNewBucket = function() {
		Buckets.createBucket("newbucket")
		.then(function (res) {
			$log.log(res.data);
			

			// TODO: websocketize. 
			Buckets.storeOneBucket(res.data);
			$scope.data.buckets = Buckets.getBuckets();
			$scope.data.cb = res.data;
			flashAlert('Created ' + res.data.meta.name + '.', 'success');
			
			$timeout(function () {
				$('#myModal-' + res.data.id).modal('show');		
			}, 500);
			
		})
		.catch(function (err) {
			$log.error(err);
			flashAlert('Failed to create bucket. Error: ' + err + '.', 'danger');
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
					flashAlert('Successfully destroyed ' + bucket.meta.name + '.', 'info');
				})
				.catch(function (err) {
					$log.log(err);
					flashAlert('Failed to destroy ' + bucket.meta.name + '. Error: ' + err + '.', 'danger');
				});	
		} else {
			// nada.
		}
		
	};

	$scope.selectBucketAsCurrent = function (bucket) {
		$scope.data.cb = bucket;
	};

	// renaming bucket
	$scope.saveBucket = function (bucket) {
		Buckets.putBucket(bucket)
			.success(function (res) {
				// $log.log(res);
				flashAlert('Successful renaming.', 'success');
				FS.fetchFS()
					.then(function (res) {
						FS.storeFS(res.data);
					})
					.catch(function (err) {
						$log.log("Error: ", err);
					});
			})
			.error(function (err) {
				$log.log(err);
				flashAlert('Renaming failed. Sry.', 'danger');
			});
	};

	$scope.import = function (path, isDir) {
		var path = path;
		if (isDir) {
			FS.importDir(path)
				.success(function (res) {
					// TODO set current bucket to the imported one
					$log.log('got res:', res);
					Buckets.storeManyBuckets(res.b);
					Snippets.setManyToSnippetsLib(res.s);
					$scope.data.snippets = Snippets.getSnippetsLib();
					$scope.data.cs = Snippets.getMostRecent(Snippets.getSnippetsLib());
					$scope.data.cb = Buckets.getBuckets()[$scope.data.cs.bucketId];
					flashAlert('Successfully imported ' + path.path + '.', 'success');
				})
				.error(function (err) {
					$log.log('importDir failed', err);
					flashAlert('Import failed for ' + path.path + '.', 'danger');
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
					flashAlert('Successfully imported ' + path.path + '.', 'success');
				})
				.error(function (err) {
					$log.log("importFile failed", err);
					flashAlert('Import failed for ' + path.path  + '.', 'danger');
				});
		}
	};

	$scope.writeSnippetToFile = function (snippet) {
		FS.writeSnippetToFile(snippet)
			.then(function (res) {
				$log.log('success! wrote snippet to file', res);
				flashAlert('Successfully saved ' + snippet.name + ' to the FS!', 'success');
				
			})
			.then(function () {
				return FS.fetchFS();
			})
			.then(function (res) {
				FS.storeFS(res.data);
			})
			.catch(function (err) {
				$log.log('error!', err);
				flashAlert('Failed to save ' + snippet.name + ' to the FS.', 'danger');
			});
	};

	$scope.writeBucketToFS = function (bucket) {
		FS.writeBucketToFS(bucket)
			.then(function (res) {
				flashAlert('Successfully wrote all snippets in ' + bucket.meta.name + ' to the FS.', 'success');
				
			})
			.then(function () {
				return FS.fetchFS();
			})
			.then(function (res) {
				FS.storeFS(res.data);
			})
			.catch(function (err) {
				flashAlert('Failed to write ' + bucket.meta.name + 'to the FS. Error: ' + err.data, 'danger');
			});
	};

	// Get associated FS path (if exists).
	$scope.getFSPathForSnippet = function (snippet) {
		// Get path name. 
		var snippetPath = snippet.name;
		var paths = FS.getFS();
		for (var i = 0; i < paths.length; i++) {
			var path = paths[i];
			if (!path.fileInfo.isDir && path.path.indexOf(Buckets.getBuckets()[snippet.bucketId].meta.name) > 0 && path.path.indexOf(snippetPath) > 0) {
				return path;
			}
		}
		return null;
	};

	$scope.comparableTimeInt = function (date) {
		return Date.parse(date);
	};

	$scope.toInt = function(timeString) {
		return parseInt(timeString, 10);
	}



}]);