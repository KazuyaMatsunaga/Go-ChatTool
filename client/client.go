package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

var running = true

func sender(conn net.Conn, name string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, _, _ := reader.ReadLine()
		if string(input) == "\\q" {
			running = false
			break
		}
		_, err := conn.Write(input)
		if err != nil {
			fmt.Println("sender write")
		}
	}
}

func receiver(conn net.Conn, name string) {
	buf := make([]byte, 560)
	for running == true {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Receiver read")
			os.Exit(1)
		}
		fmt.Println(string(buf[:n]))
		buf = make([]byte, 560)
	}
}

func main() {
	fmt.Print("Please input your name: ")
	reader := bufio.NewReader(os.Stdin)
	name, _, err := reader.ReadLine()

	// TCP://127.0.0.1:8888に接続する
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		log.Fatal("Fail to connect TCP://127.0.0.1:8888... m(_ _)m")
	}
	defer conn.Close()

	// 名前を送信する
	_, err = conn.Write(name)
	if err != nil {
		fmt.Println("Write your name.")
		os.Exit(1)
	}

	go receiver(conn, string(name))
	go sender(conn, string(name))

	for running {
		time.Sleep(1 * 1e9)
	}
}
