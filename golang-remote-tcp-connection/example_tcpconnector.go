package main

import (
	"net"
	"fmt"
	"time"
)

var stopChan = make(chan bool)
var mRemoteAddress = ":8080"

func stopConnection() {
	fmt.Println("Stopping send")
	stopChan <- true
}

func dialConnection() (net.Conn) {
	conn, err := net.Dial("tcp", mRemoteAddress)
	if err != nil {
		fmt.Println("Error: while dial: ", err)
		time.Sleep(time.Second)
		return dialConnection()
	}
	return conn
}

func getData() {
	buf := make([]byte, 256)

	for {
		buf = nil

		conn := dialConnection()
		time.Sleep(time.Second)
		fmt.Fprintf(conn, "GET /about HTTP/1.0\r\n\r\n")

		select {
		case <-stopChan:
			fmt.Println("Closing connection")
			conn.Close()
			return
		default:
			n, err := conn.Read(buf[:])
			if err != nil {
				fmt.Println("READ err: ", err, " bytes: ", n)
				continue
			}
			//Get the text, and process it
			fmt.Println( "got back: ", buf, " of ", n, " bytes")
		}
	}
}


func main() {
	go getData()

	time.Sleep(time.Second * 100)
	stopConnection()

	fmt.Println("Now you can disconnect")
}
