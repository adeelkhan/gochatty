// client
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "9988"
	SERVER_TYPE = "tcp"
)

var users = make([]string, 1)

func fetchAllUsers(conn net.Conn) {
	buffer := make([]byte, 1024)
	mLen, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading;", err.Error())
	}
	if err = json.Unmarshal(buffer[:mLen], &users); err != nil {
		panic(err)
	}
}

func processInput() string {
	var reader = bufio.NewReader(os.Stdin)
	userInput, _ := reader.ReadString('\n')
	userInput = strings.Trim(userInput, "\n")
	return userInput
}
func main() {
	// establish connection
	connection, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		panic(err)
	}

	//registering
	fmt.Printf("Enter your name:")
	userName := processInput()
	_, err = connection.Write([]byte("@reg:" + userName))
	go processServerMsg(connection)
	processClientMsg(connection)
}
func processClientMsg(conn net.Conn) {

	for {
		fmt.Printf("%s", "=>:")
		message := processInput()
		command := strings.Split(message, ":")
		if command[0] != "@all" &&
			command[0] != "@to" &&
			command[0] != "@getall" {
			continue
		}
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending message to server")
		}
	}
}
func processServerMsg(serverConn net.Conn) {
	for {
		buffer := make([]byte, 1024)
		mLen, err := serverConn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading;", err.Error())
		}
		fmt.Println("\n" + string(buffer[:mLen]))
		fmt.Printf("%s", "=>:")
	}
}
