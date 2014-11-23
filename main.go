package main

import (
  "fmt"
  "io"
  "log"
  "net/http"

  "code.google.com/p/go.net/websocket"
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

import "html/template"

func rootHandler(w http.ResponseWriter, *http.Request) {
  rootTemplate.Execute(w, listenAddr)
}

var rootTemplate = template.Must(template.New("root").Parse(`
<!DOCTYPE html>
<htmll>
<head>
<meta charset="utf-8" />
<script>
</script>
</html>
  `))