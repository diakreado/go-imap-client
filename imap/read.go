package imap

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"strings"

	. "github.com/logrusorgru/aurora"
)

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
