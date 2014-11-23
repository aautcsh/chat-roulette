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
<!DOCTYPE html>
<htmll>
<head>
<meta charset="utf-8" />
<script>
</script>
</html>
  `))

type socket struct {
  conn *websocket.Conn
  done chan bool
}

func (s socket) Read(b []byte) (int, error) { return s.conn.Read(b) }
func (s socket) Write(b []byte) (int, error) { return s.conn.Write(b) }

func (s socket) Close() error {
  s.done <- true
  return nil
}

func socketHandler(ws * websocket.Conn) {
  s := socket{conn: ws, done: make(chan bool)}
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
