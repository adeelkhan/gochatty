// client
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "9988"
	SERVER_TYPE = "tcp"
)

func main() {
	// establish connection
	connection, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Enter your name:")
	var reader = bufio.NewReader(os.Stdin)
	name, _ := reader.ReadString('\n')
	go processServerMsg(connection)

	// send some data
	for {
		fmt.Printf("%s:", name)
		var reader = bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')
		_, err = connection.Write([]byte(message))
	}
}
func processServerMsg(serverConn net.Conn) {
	for {
		buffer := make([]byte, 1024)
		mLen, err := serverConn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading;", err.Error())
		}
		fmt.Println(string(buffer[:mLen]))
	}
}
