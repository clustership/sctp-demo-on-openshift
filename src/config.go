package main

import (
  "os"
  "github.com/namsral/flag"
)

type Config struct {
    Server    bool
    DebugMode bool
    Ip        string
    Port      int
    Lport     int
    BufSize   int
    SndBuf    int
    RcvBuf    int
}

// New returns a new Config struct
func NewConfig() *Config {
  var server, debug bool
  var ip string
  var port, lport, bufsize, sndbuf, rcvbuf int

  fs := flag.NewFlagSetWithEnvPrefix(os.Args[0], "SCTP", 0)

  fs.BoolVar(&server, "server", false, "")
  fs.BoolVar(&debug, "debug", false, "")
  fs.StringVar(&ip, "ip", "0.0.0.0", "")
  fs.IntVar(&port, "port", 0, "")
  fs.IntVar(&lport, "lport", 0, "")
  fs.IntVar(&bufsize, "bufsize", 256, "")
  fs.IntVar(&sndbuf, "sndbuf", 0, "")
  fs.IntVar(&rcvbuf, "rcvbuf", 0, "")

  fs.Parse(os.Args[1:])

  return &Config{
    Server:    server,
    DebugMode: debug,
    Ip:        ip,
    Port:      port,
    Lport:     lport,
    BufSize:   bufsize,
    SndBuf:    sndbuf,
    RcvBuf:    rcvbuf,
  }

}
