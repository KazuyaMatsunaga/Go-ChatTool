package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

type Client struct {
	name  []byte
	conn  net.Conn
	color int
}

var clientList []*Client
var colorList = [5]int{32, 33, 34, 35, 36}

func send(msg []byte) {
	for _, cl := range clientList {
		_, err := cl.conn.Write(msg)
		if err != nil {
			continue
		}
	}
}

func receiver(cl *Client) {
	buf := make([]byte, 560)
	for {
		n, err := cl.conn.Read(buf)
		if err != nil {
			go send(makeMsgForQuit(string(cl.name)))
			break
		}

		go send(makeMsg(buf[:n], cl))
		buf = make([]byte, 560)
	}
}

func createClient(conn net.Conn) {
	name := getName(conn)
	color := getColor()
	cl := Client{
		name:  name,
		conn:  conn,
		color: color,
	}
	clientList = append(clientList, &cl)
	send(makeMsgForJoin(string(name)))
	go receiver(&cl)
}

func getName(conn net.Conn) []byte {
	buf := make([]byte, 560)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Fail get name... m(_ _)m")
		conn.Close()
		os.Exit(1)
	}
	return buf[:n]
}

func getColor() int {
	rand.Seed(time.Now().UnixNano())
	return colorList[rand.Intn(5)]
}

func getTime() string {
	return time.Now().Format("2006/01/02(15:04)")
}

func SprintColor(msg string, color int) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color, msg)
}

func makeMsg(msg []byte, cl *Client) []byte {
	newTime := fmt.Sprintf("%s", getTime())
	newData := fmt.Sprintf("[%s] %s", cl.name, string(msg))
	newData = SprintColor(newData, cl.color)
	newMsg := newTime + newData
	return []byte(newMsg)
}

func makeMsgForJoin(msg string) []byte {
	newMsg := fmt.Sprintf("[HERE COMES DAREDEVIL!!] %s Joined!! (∩´∀`∩) at %s", msg, getTime())
	return []byte(SprintColor(newMsg, 31))
}

func makeMsgForQuit(msg string) []byte {
	newMsg := fmt.Sprintf("%s Quit. (^ ^)ﾉｼ at %s", msg, getTime())
	return []byte(SprintColor(newMsg, 31))
}

func main() {
	// TCP://127.0.0.1:8888でListen(hostを省略する場合は':8888'のように指定)
	listenr, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		log.Fatal("Fail to listen...")
	}
	defer listenr.Close()
	fmt.Println("Ready Accept with 127.0.0.1:8888")

	// listner.Acceptは1回受け付けるとcloseしてしまうため、何度もAcceptを呼ぶ
	for {
		//　コネクションを確率する
		conn, err := listenr.Accept()
		if err != nil {
			log.Fatal("Fail to connect... m(_ _)m")
			continue
		}
		// コネクションを切断する
		defer conn.Close()

		createClient(conn)
	}
}
