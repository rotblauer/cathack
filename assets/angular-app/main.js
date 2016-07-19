'use strict';

var app = angular.module("cathack", [
	'ngWebSocket',
	'ui.codemirror',
	'contenteditable',
	'ui.router',
	'yaru22.angular-timeago'
	])
.config(function (timeAgoSettings) {
	timeAgoSettings.breakpoints = {
	  secondsToMinute: 45, // in seconds
	  secondsToMinutes: 90, // in seconds
	  minutesToHour: 45, // in minutes
	  minutesToHours: 90, // in minutes
	  hoursToDay: 24, // in hours
	  hoursToDays: 42, // in hours
	  daysToMonth: 30, // in days
	  daysToMonths: 45, // in days
	  daysToYear: 365, // in days
	  yearToYears: 1.5 // in year
	}
});
