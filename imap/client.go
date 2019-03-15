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
		successful = isOK(str[len(str)-1])
	} else {
		fmt.Println("Error->", Red("NO completed"))
	}
	logout(conn)
	return
}

// GetListOfMails - return array of string
func GetListOfMails() (messages []string) {
	authData := db.GetAuthData()
	conn := createConn(authData.Server)
	defer func() {
		conn.Close()
		fmt.Println("Connection successful close")
	}()
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	login(conn, authData.Login, authData.Password)
	examine(conn)
	response := fetchHeader(conn)
	if len(response) > 0 {
		messages = mergeResponse(response)
	} else {
		fmt.Println("Error->", Red("NO completed"))
	}
	logout(conn)
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

func examine(conn *tls.Conn) []string {
	prefix := "a0002"
	fmt.Println("Client->", Green(prefix+" examine inbox"))
	fmt.Fprintf(conn, prefix+" examine inbox\n")

	return readBeforPrefixLine(conn, prefix)
}

func fetchHeader(conn *tls.Conn) []string {
	prefix := "a0003"
	fmt.Println("Client->", Green(prefix+" fetch 1:* (ENVELOPE) "))
	fmt.Fprintf(conn, prefix+" fetch 1:* (ENVELOPE) \n")

	return readBeforPrefixLine(conn, prefix)
}

func logout(conn *tls.Conn) []string {
	prefix := "a0005"
	fmt.Println("Client->", Green(prefix+" logout"))
	fmt.Fprintf(conn, prefix+" logout\n")

	return readBeforPrefixLine(conn, prefix)
}

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

// 	2.3.4.  [RFC-2822] Size Message Attribute
//    The number of octets in the message, as expressed in [RFC-2822]
//    format.
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

func isOK(response string) (successful bool) {
	var valid = regexp.MustCompile(`^.{5} OK `)
	successful = valid.MatchString(response)
	return
}

// convert number of string to number of message
func mergeResponse(response []string) (result []string) {
	var valid = regexp.MustCompile(`^\* ([0-9]+) FETCH `)
	var lenght = len(response)
	for i, line := range response {
		if valid.MatchString(line) {
			result = append(result, line)
		} else {
			if i != lenght-1 {
				result[len(result)-1] += line
			}
		}
	}
	return
}
