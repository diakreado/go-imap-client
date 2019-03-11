package client

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

const (
	url  = "imap.rambler.ru"
	port = "993"
)

func main() {
	log.SetFlags(log.Lshortfile)
	conf := &tls.Config{}

	// Set up tls connection
	conn, err := tls.Dial("tcp", url+":"+port, conf)
	if err != nil {
		log.Println("Connection failed")
		log.Println(err)
		return
	}
	defer conn.Close()

	// Read data from Server
	go func(conn net.Conn) {
		reader := bufio.NewReader(conn)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			fmt.Printf("Server-> %s", line)
		}
	}(conn)

	// for {
	// 	reader := bufio.NewReader(os.Stdin)
	// 	fmt.Print("Text to send : ")
	// 	text, _ := reader.ReadString('\n')
	// 	commandHandler(conn, text)
	// 	// fmt.Fprintf(conn, text)
	// }

	// login(conn, "trashBag@ro.ru", "123123123")
	// time.Sleep(200 * time.Millisecond)

	// examine(conn)
	// fetch(conn)

	time.Sleep(2 * time.Second)
}

func login(conn io.Writer, login string, pass string) {
	fmt.Fprintf(conn, "a1 login "+login+" "+pass+"\n")
}

func examine(conn io.Writer) {
	fmt.Fprintf(conn, "a2 examine inbox\n")
}

func fetch(conn io.Writer) {
	fmt.Fprintf(conn, "a3 fetch 1 (body[])\n")
}

func commandHandler(conn io.Writer, command string) {
	switch command {
	case "lol\n":
		// fmt.Fprintf(conn, "lol")
		fmt.Println("lol")
	}
}
