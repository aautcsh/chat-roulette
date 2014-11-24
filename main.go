package main

import (
  "fmt"
  "io"
  "log"
  "net/http"

  "code.google.com/p/go.net/websocket"

  "html/template"
)

const LISTEN_ADDR = "localhost:4000"

func main() {
  http.HandleFunc("/", rootHandler)
  http.Handle("/socket", websocket.Handler(socketHandler))

  err := http.ListenAndServe(LISTEN_ADDR, nil)

  if err != nil {
    log.Fatal(err)
  }
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
  rootTemplate.Execute(w, LISTEN_ADDR)
}

var rootTemplate = template.Must(template.New("root").Parse(`
<!-- FROM: https://www.websocket.org/echo.html -->
<!DOCTYPE html>
<htmll>
<head>
<meta charset="utf-8" />
<script language="javascript" type="text/javascript">
//var wsUri = "ws://echo.websocket.org/"; 
var wsUri = "ws://{{.}}/socket"; 
var output;
var input;
var send;
function init() {
  output = document.getElementById("output");
  input = document.getElementById("input");
  send = document.getElementById("send");
  send.onclick = sendClickHandler;
  input.onkeydown = function(event) { if (event.keyCode == 13) send.click(); };
  testWebSocket();
}
function sendClickHandler() {
  doSend(input.value);
  input.value = '';
}
function testWebSocket() {
  websocket = new WebSocket(wsUri);
  websocket.onopen = function(evt) { onOpen(evt) };
  websocket.onclose = function(evt) { onClose(evt) };
  websocket.onmessage = function(evt) { onMessage(evt) };
  websocket.onerror = function(evt) { onError(evt) }; 
}
function onOpen(evt) {
  writeToScreen("CONNECTED");
  //doSend("WebSocket rocks");
}
function onClose(evt) {
  writeToScreen("DISCONNECTED");
}
function onMessage(evt) {
  writeToScreen('<span style="color: blue;">RESPONSE: ' + evt.data+'</span>');
  //websocket.close();
}
function onError(evt) {
  writeToScreen('<span style="color: red;">ERROR:</span> ' + evt.data); 
}
function doSend(message) {
  writeToScreen("SENT: " + message);
  websocket.send(message);
}
function writeToScreen(message) {
  var pre = document.createElement("p");
  pre.style.wordWrap = "break-word";
  pre.innerHTML = message;
  output.appendChild(pre);
}
window.addEventListener("load", init, false);
</script>
<h2>WebSocket Test</h2>
<input type="text" id="input" /><input type="button" id="send" value="Send" />
<div id="output"></div>
</html>
  `))

type socket struct {
  io.ReadWriter
  done chan bool
}

func (s socket) Close() error {
  s.done <- true
  return nil
}

func socketHandler(ws * websocket.Conn) {
  s := socket{ws, make(chan bool)}
  go match(s)
  <- s.done
}

var partner = make(chan io.ReadWriteCloser)

func match(c io.ReadWriteCloser) {
    fmt.Fprint(c, "Waiting for a partner...")
    
    select {
    case partner <- c:
        // now handled by the other goroutine
    case p := <-partner:
        chat(p, c)
    }
}

func chat(a, b io.ReadWriteCloser) {
    fmt.Fprintln(a, "Found one! Say hi.")
    fmt.Fprintln(b, "Found one! Say hi.")
    
    errc := make(chan error, 1)
    
    go cp(a, b, errc)
    go cp(b, a, errc)
    
    if err := <-errc; err != nil {
        log.Println(err)
    }
    
    a.Close()
    b.Close()
}

func cp(w io.Writer, r io.Reader, errc chan<- error) {
    _, err := io.Copy(w, r)
    errc <- err
}
