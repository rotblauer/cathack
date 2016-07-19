'use strict';

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