package imap

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	db "../db"
	. "github.com/logrusorgru/aurora"
)

const (
	port = ":993"
)

// TryToLogin - return result of login
func TryToLogin() (successful bool) {
	authData := db.GetAuthData()

	conn := createConn(authData.Server)
	defer func() {
		conn.Close()
		fmt.Println("Connection successful close")
	}()

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	successful = false
	str := login(conn, authData.Login, authData.Password)
	if len(str) > 0 {
		successful = responseParser(str[len(str)-1])
	} else {
		fmt.Println("Error->", Red("NO completed"))
	}

	str = logout(conn)

	return
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

func login(conn *tls.Conn, login string, pass string) []string {
	prefix := "a0001"
	fmt.Println("Client->", Green(prefix+" login "+login+"  ******"))
	fmt.Fprintf(conn, prefix+" login "+login+" "+pass+"\n")

	return readBeforPrefixLine(conn, prefix)
}

func examine(conn io.Writer) {
	fmt.Println("Client->", Green("a2 examine inbox"))
	fmt.Fprintf(conn, "a2 examine inbox\n")
}

func fetch(conn io.Writer) {
	fmt.Println("Client->", Green("a3 fetch 1 (body[])"))
	fmt.Fprintf(conn, "a3 fetch 1 (body[])\n")
}

func logout(conn *tls.Conn) []string {
	prefix := "a0004"
	fmt.Println("Client->", Green(prefix+" logout"))
	fmt.Fprintf(conn, prefix+" logout\n")

	return readBeforPrefixLine(conn, prefix)
}

// 	2.3.4.  [RFC-2822] Size Message Attribute
//    The number of octets in the message, as expressed in [RFC-2822]
//    format.

// readStrings - read the specified number of lines
// conn : reading stream
// num : number of lines which should be reading
// text : result of reading
func readStrings(conn io.Reader, num int) (text []string) {
	reader := bufio.NewReader(conn)
	for index := 0; index < num; index++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error->", Red(err))
			break
		}
		text = append(text, line)
		fmt.Printf("Server-> %s", Cyan(line))
	}
	return
}

func readBytes() {

}

func readBeforPrefixLine(conn *tls.Conn, prefix string) (response []string) {
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error->", Red(err))
			break
		}
		response = append(response, line)
		fmt.Printf("Server-> %s", Cyan(line))
		if strings.HasPrefix(line, prefix) {
			break
		}
	}
	return
}

func responseParser(response string) (successful bool) {
	var valid = regexp.MustCompile(`^.{5} OK `)
	successful = valid.MatchString(response)
	return
}
