package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

type User struct {
	Name string
	conn net.Conn
}

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "9988"
	SERVER_TYPE = "tcp"
)

var user_count = 0
var users = make(map[string]User, 1)

func main() {

	// duummy users
	// users["adeel"] = User{
	// 	Name: "adeel",
	// 	conn: nil,
	// }
	// users["wahaj"] = User{
	// 	Name: "wahaj",
	// 	conn: nil,
	// }

	fmt.Println("Chatserver running... :)")
	server, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer server.Close()
	fmt.Println("listening on " + SERVER_HOST + ":" + SERVER_PORT)
	fmt.Println("Waiting for users...")

	for {
		connection, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			os.Exit(1)
		}

		go processUser(connection)
		user_count++
	}
}
func sendAllUsers(connection net.Conn) {
	fmt.Println("Sending all users...")
	list := make([]string, 0)
	for _, v := range users {
		list = append(list, v.Name)
	}
	all_users, err := json.Marshal(&list)
	_, err = connection.Write([]byte(string(all_users)))
	if err != nil {
		fmt.Println("Cant able to send the list to client")
	}
}

func registerUser(connection net.Conn, user string) string {
	_, ok := users[user]
	if !ok {
		users[user] = User{Name: user, conn: connection}
	}
	return user
}
func sendToAll(conn net.Conn, myId string, msg string) {
	myName := users[myId].Name
	// broadcasting to other clients too
	for userId, user := range users {
		if myId != userId {
			conn := user.conn
			_, err := conn.Write([]byte(string(myName + ":" + msg)))
			if err != nil {
				fmt.Println("Error writing:", err.Error())
			}
		}
	}
}
func sentTo(conn net.Conn, myId string, sendTo string, msg string) {
	user, ok := users[sendTo]
	// fmt.Println(user, ok)
	fmt.Println("Sending to... ", user.Name)
	if ok {
		conn := user.conn
		_, err := conn.Write([]byte(string(myId + ":" + msg)))
		if err != nil {
			fmt.Println("Error writing:", err.Error())
		}
	}
}
func processUser(conn net.Conn) {
	var genUserId string
	for {
		buffer := make([]byte, 1024)
		mLen, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}
		command := strings.Split(string(buffer[:mLen]), ":")
		if command[0] == "@reg" {
			genUserId = command[1]
			registerUser(conn, genUserId)
			fmt.Println(users)
		} else if command[0] == "@all" {
			msg := command[1]
			sendToAll(conn, genUserId, msg)
		} else if command[0] == "@disc" {
			// send message to all clients to update
			// _, err := conn.Write([]byte(string(myName + ":" + msg)))
			// if err != nil {
			// 	fmt.Println("Error writing:", err.Error())
			// }
			// connection.Close()
		} else if command[0] == "@getall" {
			sendAllUsers(conn)
		} else {
			fmt.Println(command[0], command[1], command[2])
			to, msg := command[1], command[2]
			sentTo(conn, genUserId, to, msg)
		}
	}

}
