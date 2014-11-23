package main

import (
  "fmt"
  "io"
  "log"
  "net"
)

const LISTEN_ADDR = "localhost:4000"

var partner = make(chan io.ReadWriteCloser)

func main() {
  l, err := net.Listen("tcp", LISTEN_ADDR)
  
  if err != nil {
    log.Fatal(err)
  }

  for {
    c, err := l.Accept()
    
    if err != nil {
      log.Fatal(err)
    }

    go match(c)
  }
}

func match(c io.ReadWriteCloser) {
  fmt.Fprint(c, "Waiting for a partner...")
  select {
  case partner <- c:
    // now handled by the other goroutine
  case p := <- partner:
    chat(p, c)
  }
}

func chat(a, b io.ReadWriteCloser) {
  fmt.Fprintln(a, "Found one! Say hi.")
  fmt.Fprintln(b, "Found one! Say hi.")
  go io.Copy(a, b)
  io.Copy(b, a)
}