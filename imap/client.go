package imap

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"strings"
	"time"

	db "../db"
	. "github.com/logrusorgru/aurora"
)

const (
	port = ":993"
)

// GetPostBoxState - return state of post-box
func GetPostBoxState() {
	authData := db.GetAuthData()

	conn := createConn(authData.Server)
	defer func() {
		conn.Close()
		fmt.Println("Connection successful close")
	}()

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	str := readStrings(conn, 1)
	fmt.Print(Magenta(str))
	str = login(conn, authData.Login, authData.Password)
	fmt.Print(Magenta(str))
	logout(conn)
	str = readStrings(conn, 2)
	fmt.Print(Magenta(str))
}

// Set up tls connection
func createConn(server string) *tls.Conn {
	conn, err := tls.Dial("tcp", server+port, &tls.Config{})
	if err != nil {
		fmt.Println(Red("Connection failed"))
		panic(err)
	}
	fmt.Println("Connection successful open")
	return conn
}

func login(conn *tls.Conn, login string, pass string) (response string) {
	prefix := "a1"
	fmt.Println("Client->", Green(prefix+" login "+login+"  ******"))
	fmt.Fprintf(conn, prefix+" login "+login+" "+pass+"\n")

	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		response += line
		fmt.Printf("Server-> %s", Cyan(line))
		if strings.HasPrefix(line, prefix) {
			break
		}
	}
	return
}

func examine(conn io.Writer) {
	fmt.Println("Client->", Green("a2 examine inbox"))
	fmt.Fprintf(conn, "a2 examine inbox\n")
}

func fetch(conn io.Writer) {
	fmt.Println("Client->", Green("a3 fetch 1 (body[])"))
	fmt.Fprintf(conn, "a3 fetch 1 (body[])\n")
}

func logout(conn io.Writer) {
	fmt.Println("Client->", Green("a4 logout"))
	fmt.Fprintf(conn, "a4 logout\n")
}

// 	2.3.4.  [RFC-2822] Size Message Attribute
//    The number of octets in the message, as expressed in [RFC-2822]
//    format.

// readStrings - read the specified number of lines
// conn : reading stream
// num : number of lines which should be reading
// text : result of reading
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

func readBytes() {

}
