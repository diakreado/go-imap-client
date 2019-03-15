package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"time"

	. "github.com/logrusorgru/aurora"
)

const (
	url  = "imap.yandex.ru"
	port = ":993"
)

// GetPostBoxState - return state of post-box
func main() {
	conn, err := tls.Dial("tcp", url+port, &tls.Config{})
	if err != nil {
		fmt.Println(Red("Connection failed"))
		panic(err)
	}
	defer func() {
		conn.Close()
		fmt.Println("Connection successful close")
	}()
	fmt.Println("Connection successful open")

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	str := readStrings(conn, 1)
	fmt.Print(Magenta(str))

	fmt.Fprintf(conn, "a1 login some text\n")

	// str = readStrings(conn, 2)
	// fmt.Print(Magenta(str))

	time.Sleep(1 * time.Second)
}

func readStrings(conn io.Reader, num int) (text string) {
	reader := bufio.NewReader(conn)
	for index := 0; index < num; index++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		text += line
		fmt.Printf("Server-> %s", Cyan(line))
	}
	return
}
