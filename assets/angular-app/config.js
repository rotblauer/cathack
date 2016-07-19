'use strict';


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
		],
		readOnly: false,
		lineWrapping: true
	},
	"DEFAULTSNIPPET": {
		id: Math.random().toString(36).substring(9),
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