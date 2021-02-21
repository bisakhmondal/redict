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

	if len(os.Args) > 1 {
		parseArgs(conn)
	}
	
	stdreader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("127.0.0.1:%s>> ", PORT)
		cmd, _  := stdreader.ReadString('\n')
		cmd = strings.Trim(cmd, "\n")
		cmdArr := strings.Split(cmd, " ")

		if Communicate(conn, cmdArr, cmd) {
			return
		}

	}
}

func parseArgs(conn net.Conn){
	i := 1
	for i< len(os.Args) {
		switch os.Args[i] {
		case "get":
			cmd := os.Args[i] +" "+ os.Args[i+1]
			Communicate(conn, os.Args[i:i+2], cmd)
			i+=2
		case "put":
			cmd := strings.Join(os.Args[i:i+3]," ")
			Communicate(conn, os.Args[i:i+3], cmd)
			i+=3
		default:
			i++
			fmt.Println("unknown initializer: ", os.Args[i])
		}
	}
}

func Communicate(conn net.Conn, cmdArr []string, cmd string) bool {
	buffer := make([]byte, BUFFERSIZE)
	switch strings.ToLower(cmdArr[0]) {

	case "get":
		if len(cmdArr) < 2 {
			fmt.Println("Get requires a key to fetch value")
			return false
		}
		conn.Write([]byte(cmd))
		n, _ := conn.Read(buffer)
		fmt.Println(string(buffer[:n]))

	case "put":
		if len(cmdArr) < 3 {
			fmt.Println("Put requires a key & Value pair store")
			return false
		}

		conn.Write([]byte(cmd))

	case "upgrade":
		if len(cmdArr) < 2 {
			fmt.Println("Get requires a key to fetch value")
			return false
		}
		conn.Write([]byte(cmd))
		n, _ := conn.Read(buffer)
		fmt.Println(string(buffer[:n]))

	case "stats":
		conn.Write([]byte(cmd))
		n, _ := conn.Read(buffer)
		fmt.Println(string(buffer[:n]))

	case "downgrade":
		if len(cmdArr) < 2 {
			fmt.Println("Downgrade requires a UID of the specific user")
			return false
		}
		conn.Write([]byte(cmd))
		n, _ := conn.Read(buffer)
		fmt.Println(string(buffer[:n]))

	case "close":
		conn.Write([]byte("close"))
		return true

	default:
		fmt.Println("Unknown Command!! Try again")
		fmt.Println("Here is the List...")
		fmt.Println("1. get <key>")
		fmt.Println("2. put <key> <value>")
		fmt.Println("3. upgrade <secret>")
		fmt.Println("4. downgrade <uid>")
		fmt.Println("5. close\n")
	}
	return false
}