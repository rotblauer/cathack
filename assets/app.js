
var chat = document.getElementById("chat");
var url = "ws://" + window.location.host + "/ws";
var ws = new WebSocket(url);
var chat = document.getElementById("chat");
var text = document.getElementById("text");

function formatDate(goTimeString) {

  var cleanString = goTimeString.replace(' UTC', '');
  var a = moment.utc(cleanString);
  var b = a.local().format(" kk:mm:ss ddd MMM D");

  if (a.isValid()) {
    return b;  
  } else {
    return goTimeString
  }
}

function getDateFromLine(line) {
  return line.split(',')[0];
}

function emphasizeHTML(string) {
  return '<em>' + string + '</em>';
}
function smallifyHTML(string) {
  return '<small>' + string + '</small>';
}
function strongifyHTML(string) {
  return '<strong>' + string + '</strong>'; 
}

function handleLineFormatting(line) {

  console.log('got line -> ' + line);

  // grab time and make it look nice
  var niceTime = formatDate(getDateFromLine(line));
  console.log('niceTime-> ' + niceTime);
  
  // replace original time with nice time
  var better = line.replace(getDateFromLine(line), smallifyHTML(niceTime));
  
  var s = line.split(','); // split from original line
  var lat = s[1];
  var lon = s[2];
  var tz = s[3];
  var subdiv = s[4];

  better = better.replace(subdiv, strongifyHTML(subdiv));
  better = better.replace(lat, '');
  better = better.replace(lon, '');
  better = better.replace(tz, '');
  better = better.replace(',,,,', ' ');
  better = better.replace(/,.*\$/, ' $'); // replace between , ... $

  // var best = '';
  // best += niceTime + ' '; 
  // best += subdiv + ' ';
  // // remove erything cept msg
  // best += better.replace(/^.*\$/, '$');

  // Remove (time_zone)
  return better; // .replace(/ *\([^)]*\) */g, " ");
  // return best;
}

function scrollChat() {
  var div = $('#chat');
  div.scrollTop(div.prop('scrollHeight'));
}

//Seems like this might not be the way to handle loadin up a text file
function LoadFile() {
    var oFrame = document.getElementById("frmFile");

    var strRawContents = oFrame.contentWindow.document.body.childNodes[0].innerHTML;
    var arrLines = strRawContents.split("\n");
    for (var i = 0; i < arrLines.length; i++) {
        var curLine = arrLines[i];
        // chat.innerText += handleLineFormatting(curLine) + "\n";
        chat.innerHTML += handleLineFormatting(curLine) + "\n";
    }
}

ws.onmessage = function (msg) {
  var line = msg.data;

  // handle line formatting
  // chat.innerText += handleLineFormatting(line) + "\n";
  chat.innerHTML += handleLineFormatting(line) + "\n";

  scrollChat(); // set to bottom nicelike
};
text.onkeydown = function (e) {
  if (e.keyCode === 13 && text.value !== "") {
    ws.send(" $ " + text.value);
    text.value = "";
  }
};

// scroll to bottom on doc ready
$(function () {
  scrollChat();
});
