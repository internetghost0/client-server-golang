package server

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
		panic("the ip-address is already taken")
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
	nick, err := recvConn(conn)
	if err != nil {
		return
	}
	if strings.ToLower(nick) == "server" {
		conn.Write([]byte("Server: you cannot use this nickname\n" + END_BYTES))
		conn.Close()
		return
	}
	if len(nick) > 64 {
		nick = nick[:64]
	}
	nick = strings.ReplaceAll(nick, " ", "_")
	nick = strings.ReplaceAll(nick, "`", "_")

	for c := range Connections {
		if nick == Connections[c] {
			conn.Write([]byte("Server: this nickname is already taken" + END_BYTES))
			conn.Close()
			return
		}
	}

	Connections[conn] = nick
	sendEveryoneExcept("Hello to server!\n", nil)
	conn.Write([]byte("Server: your nickname is `" + Connections[conn] + "`\n"))

	conn.Write([]byte("Online users:\n"))
	for c := range Connections {
		conn.Write([]byte("   `" + Connections[c] + "`\n"))
	}
	conn.Write([]byte(END_BYTES))
	sendEveryoneExcept("Server: new user `"+Connections[conn]+"` has connected", conn)
	for {
		msg, err := recvConn(conn)
		if err != nil {
			sendEveryoneExcept("Server: user `"+Connections[conn]+"` has disconnected ", conn)
			break
		}
		sendEveryoneExcept(Connections[conn]+": "+msg, conn)
	}
	delete(Connections, conn)

}

func recvConn(conn net.Conn) (string, error) {
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

func sendEveryoneExcept(msg string, conn net.Conn) {
	for c := range Connections {
		if c != conn {
			c.Write([]byte(msg + END_BYTES))
		}
	}
}
