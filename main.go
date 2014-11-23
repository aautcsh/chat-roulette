package main

import (
  "io"
  "log"
  "net"
)

const LISTEN_ADDR = "localhost:4000"

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

    go io.Copy(c, c)
  }
}