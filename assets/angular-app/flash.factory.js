'use strict';

app.factory("flash", function($rootScope) {
  // var queue = [];
  var currentMessage = {
  	message: "",
  	klass: ""
  };

  // $rootScope.$on("$routeChangeSuccess", function() {
  //   currentMessage = queue.shift() || "";
  // });

  return {
    setMessage: function(message, klass) {
     	currentMessage = {
     		message: message,
     		klass: klass
     	};
    },

    getMessage: function() {
      return currentMessage;
    }
  };
});
