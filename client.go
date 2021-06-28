package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	ADDR      = ":8080"
	END_BYTES = "\001\003\003\007"
)

func main() {
	conn, err := net.Dial("tcp6", ADDR)
	if err != nil {
		panic("error to connect")
	}
	defer conn.Close()
	sendConn(conn, input("your nick: "))
	for {
		go clientOutput(conn)
		clientInput(conn)
	}
}

func clientInput(conn net.Conn) {
	for {
		sendConn(conn, input(""))
	}
}

func clientOutput(conn net.Conn) {
	for {
		msg, err := recvConn(conn)
		if err != nil {
			fmt.Println("Disconnected from server")
			os.Exit(0)
		}
		fmt.Println(msg)
	}
}

func input(msg string) string {
	fmt.Print(msg)
	reader := bufio.NewReader(os.Stdin)
	result, _ := reader.ReadString('\n')

	return strings.Trim(result, "\n")

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

func sendConn(conn net.Conn, message string) (int, error) {
	return conn.Write([]byte(message + END_BYTES))
}
