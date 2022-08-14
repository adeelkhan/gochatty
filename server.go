package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "9988"
	SERVER_TYPE = "tcp"
)

var user_count = 0
var users = make(map[string]net.Conn, 1)

func main() {
	fmt.Println("Chatserver running... :)")
	server, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer server.Close()
	fmt.Println("listening on " + SERVER_HOST + ":" + SERVER_PORT)
	fmt.Println("Waiting for client...")

	for {
		connection, err := server.Accept()
		// registerting user
		userId := "client" + strconv.Itoa(user_count)
		users[userId] = connection
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			os.Exit(1)
		}
		fmt.Printf("%s connected...\n", userId)
		go processClient(userId, connection)
		user_count++
	}
}
func processClient(userId string, connection net.Conn) {

	for {
		buffer := make([]byte, 1024)
		mLen, err := connection.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}
		fmt.Printf("%s: %s", userId, string(buffer[:mLen]))
		_, err = connection.Write([]byte("server: " + string(buffer[:mLen])))

		// broadcasting to other clients too
		for Id, conn := range users {
			if Id != userId {
				_, err = conn.Write([]byte("server: " + string(buffer[:mLen])))
				if err != nil {
					fmt.Println("Error writing:", err.Error())
				}
			}
		}
	}
	// connection.Close()
}
