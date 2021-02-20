package main

import (
	"fmt"
	"log"
	"net"
	"redict/server/namespace"
	"redict/server/persistence"
	"strconv"
	"strings"
)

//Access Privileges
const (
	user Role = iota
	Manager
)
type Role int

const (
	PORT = "5000"
	BUFFERSIZE = 4096
	ManagerSecret = "redict"
)

var clientAttached = 0
var pool *namespace.Container


func main(){
	pool = namespace.NewContainer()
	var rdb *persistence.RDB = nil //persistence.RDBInit()


	server, err := net.Listen("tcp", ":"+PORT)


	defer server.Close()

	if err != nil{
		log.Fatal(err)
	}
	log.Println("TCP server is UP @ localhost: ", PORT)

	for{
		connection, err := server.Accept()
		if err != nil{
			log.Println("Client Connection failed")
		}
		newNamespace := namespace.New(clientAttached, rdb)

		pool.Push(newNamespace)

		go handleClient(connection, newNamespace)
		clientAttached++
	}
}

func handleClient(conn net.Conn, namespace *namespace.Namespace){
	defer conn.Close()

	roleAssigned := user
	log.Println("Client ", clientAttached)
	buffer := make([]byte, BUFFERSIZE)
	for {
		n,_ := conn.Read(buffer)
		command := string(buffer[:n])
		commandArr := strings.Split(command," ")
		fmt.Println(command)

		switch strings.ToLower(commandArr[0]) {
		case "put":
			if roleAssigned == user {
				//user put into own namespace
				namespace.Put(commandArr[1], commandArr[2])
			}else{
				//manager update all namespaces
				pool.Put(commandArr[1], commandArr[2])
			}
		case "get":
			var value string
			if roleAssigned == user {
				//get the value from current Namespace
				 value = namespace.Get(commandArr[1])
			}else{
				value = pool.Get(commandArr[1])
			}
			conn.Write([]byte(value))

		case "upgrade":
			var message string
			if commandArr[1] == ManagerSecret {

				roleAssigned = Manager
				message = "Role Upgraded Successfully"
			}else{
				message = "Incorrect Password"
			}
			conn.Write([]byte(message))

		case "downgrade":
			if roleAssigned == user {
				conn.Write([]byte("Role Switch is only meant for Managers!!"))
				continue
			}

			uid,  err := strconv.Atoi(commandArr[1])
			if err !=nil{
				conn.Write([]byte(err.Error()))
				continue
			}
			var message string
			namespace = pool.GetNamespace(uid)
			if  namespace == nil {
				message = "Invalid Namespace uid"
			}else{
				message = fmt.Sprintf("Switched to User %d successfully.", uid)
				roleAssigned = user
			}
			conn.Write([]byte(message))

		case "stats":
			var message string

			if roleAssigned == Manager {
				message = pool.GetStats()
			}else{
				message = "whole containers Stats can only be accessed by a Manager"
			}
			conn.Write([]byte(message))

		case "close":
			pool.Delete(namespace)

			log.Println("Client Switched off: ", namespace.GetUID())
			return

		}
		
	}

}
