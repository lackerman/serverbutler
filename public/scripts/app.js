function getData(url, callback) {
  const request = new XMLHttpRequest();
  request.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200) {
      callback(request.responseText);
    }
  };
  request.open("GET", url, true);
  request.send();
}

const wsUri = "ws:///";
let websocket;
let output;

function init() {
  output = document.getElementById("output");
  websocket = new WebSocket(wsUri);
  websocket.onopen = evt => onOpen(evt);
  websocket.onclose = evt => onClose(evt);
  websocket.onmessage = evt => onMessage(evt);
  websocket.onerror = evt => onError(evt);
}

function onOpen(evt) {
  writeToScreen("CONNECTED");
  doSend("WebSocket rocks");
}

function onClose(evt) {
  writeToScreen("DISCONNECTED");
}

function onMessage(evt) {
  writeToScreen('<span style="color: blue;">RESPONSE: ' + evt.data + "</span>");
}

function onError(evt) {
  writeToScreen('<span style="color: red;">ERROR:</span> ' + evt.data);
}

function doSend(message) {
  writeToScreen("SENT: " + message);
  websocket.send(message);
}

function writeToScreen(message) {
  const span = document.createElement("span");
  span.style.display = "block";
  span.innerHTML = message;
  output.appendChild(span);
}

window.addEventListener("load", init, false);
