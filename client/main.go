package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)
const (
	PORT = "5000"
	HOST = "localhost"
	BUFFERSIZE = 4096
)

func CatchSignal(conn net.Conn) {

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-sigc
		//Tell server to close the thread assigned to this particular client
		conn.Write([]byte("close"))

		os.Exit(0)
	}()
}

func main(){
	server, err := net.Dial("tcp", HOST+":"+PORT)

	defer server.Close()

	if err != nil{
		log.Fatal(err)
	}
	log.Println("Client is Connected to Server @ ",HOST,":", PORT)

	CatchSignal(server)
	handleConnection(server)
}

func handleConnection(conn net.Conn){
	stdreader := bufio.NewReader(os.Stdin)
	buffer := make([]byte, BUFFERSIZE)

	for {

		fmt.Printf("127.0.0.1:%s>> ", PORT)
		cmd, _  := stdreader.ReadString('\n')
		cmd = strings.Trim(cmd, "\n")
		cmdArr := strings.Split(cmd, " ")

		switch strings.ToLower(cmdArr[0]){

		case "get":
			if len(cmdArr) <2 {
				fmt.Println("Get requires a key to fetch value")
				continue
			}
			conn.Write([]byte(cmd))
			n, _ := conn.Read(buffer)
			fmt.Println(string(buffer[:n]))

		case "put":
			if len(cmdArr) <3 {
				fmt.Println("Put requires a key & Value pair store")
				continue
			}

			conn.Write([]byte(cmd))
		case "upgrade":
			if len(cmdArr) <2 {
				fmt.Println("Get requires a key to fetch value")
				continue
			}
			conn.Write([]byte(cmd))
			n, _ := conn.Read(buffer)
			fmt.Println(string(buffer[:n]))
		case "stats":
			conn.Write([]byte(cmd))
			n, _ := conn.Read(buffer)
			fmt.Println(string(buffer[:n]))

		case "downgrade":
			if len(cmdArr) <2 {
				fmt.Println("Downgrade requires a UID of the specific user")
				continue
			}
			conn.Write([]byte(cmd))
			n, _ := conn.Read(buffer)
			fmt.Println(string(buffer[:n]))
		}

	}
}