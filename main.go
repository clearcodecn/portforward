package main

import (
	"flag"
	"io"
	"log"
	"net"
)

var (
	listenAddr string
	destAddr   string
)

func init() {
	flag.StringVar(&listenAddr, "l", "", "listen address")
	flag.StringVar(&destAddr, "a", "", "forward address")
}

func main() {
	flag.Parse()
	if listenAddr == "" || destAddr == "" {
		log.Fatal("listen address or forward address can't be empty")
	}

	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go func() {
			remoteConn, err := net.Dial("tcp", destAddr)
			if err != nil {
				log.Println("[ERRO] failed to dial remote", err)
				return
			}

			log.Printf("[INFO] accepting %v to %v", conn.RemoteAddr().String(), remoteConn.RemoteAddr())

			defer remoteConn.Close()
			defer conn.Close()

			go io.Copy(remoteConn, conn)
			io.Copy(conn, remoteConn)
		}()
	}
}
