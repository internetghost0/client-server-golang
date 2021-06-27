package main

import (
	"net"
	"strings"
)

var (
	Connections = make(map[net.Conn]string)
)

const (
	ADDR      = ":8080"
	END_BYTES = "\001\003\003\007"
)

func main() {
	listener, err := net.Listen("tcp6", ADDR)
	if err != nil {
		panic("server error")
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go handleConnection(conn)

	}

}

func handleConnection(conn net.Conn) {

	nick, err := ConnRead(conn)
	if err != nil {
		return
	}
	Connections[conn] = nick
	for {
		msg, err := ConnRead(conn)
		if err != nil {
			break
		}
		for c := range Connections {
			if c != conn {
				ConnWrite(c, Connections[c]+": "+msg)
			}
		}
	}
	delete(Connections, conn)

}

func ConnRead(conn net.Conn) (string, error) {
	var (
		bytes   = make([]byte, 1024)
		message = ""
	)
	for {
		length, err := conn.Read(bytes)
		if err != nil {
			return "nil", err
		}
		message += string(bytes[:length])
		if strings.HasSuffix(message, END_BYTES) {
			message = strings.TrimSuffix(message, END_BYTES)
			break
		}
	}
	return message, nil
}

func ConnWrite(conn net.Conn, message string) (int, error) {
	return conn.Write([]byte(message + END_BYTES))
}
