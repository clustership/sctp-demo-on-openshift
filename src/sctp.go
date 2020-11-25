package main

import (
	"log"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/ishidawataru/sctp"
)

func serveClient(conn net.Conn, bufsize int) error {
	for {
		buf := make([]byte, bufsize+128) // add overhead of SCTPSndRcvInfoWrappedConn
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("read failed: %v", err)
			return err
		}
		log.Printf("read: %d", n)
		n, err = conn.Write(buf[:n])
		if err != nil {
			log.Printf("write failed: %v", err)
			return err
		}
		log.Printf("write: %d", n)
	}
}

func main() {
  log.Println("Loading config")
  conf := NewConfig()

  log.Printf("Configured port: %d\n", conf.Port)

  var rcvbuf, sndbuf int
	ips := []net.IPAddr{}

	for _, i := range strings.Split(conf.Ip, ",") {
		if a, err := net.ResolveIPAddr("ip", i); err == nil {
			log.Printf("Resolved address '%s' to %s", i, a)
			ips = append(ips, *a)
		} else {
			log.Printf("Error resolving address '%s': %v", i, err)
		}
	}

	addr := &sctp.SCTPAddr{
		IPAddrs: ips,
		Port:    conf.Port,
	}
	log.Printf("raw addr: %+v\n", addr.ToRawSockAddrBuf())

	if conf.Server {
		ln, err := sctp.ListenSCTP("sctp", addr)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		log.Printf("Listen on %s", ln.Addr())

		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Fatalf("failed to accept: %v", err)
			}
			log.Printf("Accepted Connection from RemoteAddr: %s", conn.RemoteAddr())
			wconn := sctp.NewSCTPSndRcvInfoWrappedConn(conn.(*sctp.SCTPConn))
			if conf.SndBuf != 0 {
				err = wconn.SetWriteBuffer(conf.SndBuf)
				if err != nil {
					log.Fatalf("failed to set write buf: %v", err)
				}
			}
			if conf.RcvBuf != 0 {
				err = wconn.SetReadBuffer(conf.RcvBuf)
				if err != nil {
					log.Fatalf("failed to set read buf: %v", err)
				}
			}
			sndbuf, err = wconn.GetWriteBuffer()
			if err != nil {
				log.Fatalf("failed to get write buf: %v", err)
			}
			rcvbuf, err = wconn.GetWriteBuffer()
			if err != nil {
				log.Fatalf("failed to get read buf: %v", err)
			}
			log.Printf("SndBufSize: %d, RcvBufSize: %d", sndbuf, rcvbuf)

			go serveClient(wconn, conf.BufSize)
		}

	} else {
		var laddr *sctp.SCTPAddr
		if conf.Lport != 0 {
			laddr = &sctp.SCTPAddr{
				Port: conf.Lport,
			}
		}
		conn, err := sctp.DialSCTP("sctp", laddr, addr)
		if err != nil {
			log.Fatalf("failed to dial: %v", err)
		}
    log.Println("conn established")

		log.Printf("Dial LocalAddr: %s; RemoteAddr: %s", conn.LocalAddr(), conn.RemoteAddr())

		if conf.SndBuf != 0 {
			err = conn.SetWriteBuffer(conf.SndBuf)
			if err != nil {
				log.Fatalf("failed to set write buf: %v", err)
			}
		}
		if conf.RcvBuf != 0 {
			err = conn.SetReadBuffer(conf.RcvBuf)
			if err != nil {
				log.Fatalf("failed to set read buf: %v", err)
			}
		}

		sndbuf, err = conn.GetWriteBuffer()
		if err != nil {
			log.Fatalf("failed to get write buf: %v", err)
		}
		rcvbuf, err = conn.GetReadBuffer()
		if err != nil {
			log.Fatalf("failed to get read buf: %v", err)
		}
		log.Printf("SndBufSize: %d, RcvBufSize: %d", sndbuf, rcvbuf)

		ppid := 0
		for {
			info := &sctp.SndRcvInfo{
				Stream: uint16(ppid),
				PPID:   uint32(ppid),
			}
			ppid += 1
			conn.SubscribeEvents(sctp.SCTP_EVENT_DATA_IO)
			buf := make([]byte, conf.BufSize)
      log.Printf("Generating random string (size: %d)\n", conf.BufSize)

			n, err := rand.Read(buf)
			if n != conf.BufSize {
				log.Fatalf("failed to generate random string len: %d", conf.BufSize)
			}
			n, err = conn.SCTPWrite(buf, info)
			if err != nil {
				log.Fatalf("failed to write: %v", err)
			}
			log.Printf("write: len %d", n)
			n, info, err = conn.SCTPRead(buf)
			if err != nil {
				log.Fatalf("failed to read: %v", err)
			}
			log.Printf("read: len %d, info: %+v", n, info)
			time.Sleep(time.Second)
		}
	}
}
