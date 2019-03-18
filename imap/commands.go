package imap

import (
	"crypto/tls"
	"fmt"
	"strconv"

	. "github.com/logrusorgru/aurora"
)

func login(conn *tls.Conn, login string, pass string) []string {
	prefix := "a0001"
	fmt.Println("Client->", Green(prefix+" login "+login+"  ******"))
	fmt.Fprintf(conn, prefix+" login "+login+" "+pass+"\n")

	return readBeforPrefixLine(conn, prefix)
}

func examineInbox(conn *tls.Conn) []string {
	prefix := "a0002"
	fmt.Println("Client->", Green(prefix+" examine inbox"))
	fmt.Fprintf(conn, prefix+" examine inbox\n")

	return readBeforPrefixLine(conn, prefix)
}

func selectInbox(conn *tls.Conn) []string {
	prefix := "a0003"
	fmt.Println("Client->", Green(prefix+" select inbox"))
	fmt.Fprintf(conn, prefix+" select inbox\n")

	return readBeforPrefixLine(conn, prefix)
}

func fetchHeader(conn *tls.Conn) []string {
	prefix := "a0004"
	fmt.Println("Client->", Green(prefix+" fetch 1:* (ENVELOPE UID) "))
	fmt.Fprintf(conn, prefix+" fetch 1:* (ENVELOPE UID) \n")

	return readBeforPrefixLine(conn, prefix)
}

func searchUnseen(conn *tls.Conn) []string {
	prefix := "a0005"
	fmt.Println("Client->", Green(prefix+" search unseen"))
	fmt.Fprintf(conn, prefix+" search unseen\n")

	return readBeforPrefixLine(conn, prefix)
}

func findLetter(conn *tls.Conn, uid string) []string {
	prefix := "a0006"
	fmt.Println("Client->", Green(prefix+" search UID "+uid))
	fmt.Fprintf(conn, "%s search UID %s\n", prefix, uid)

	return readBeforPrefixLine(conn, prefix)
}

func fetchLetter(conn *tls.Conn, num int) []string {
	prefix := "a0007"
	fmt.Println("Client->", Green(prefix+" fetch "+strconv.Itoa(num)+" (body[])"))
	fmt.Fprintf(conn, "%s fetch %d (body[]) \n", prefix, num)

	return readBeforPrefixLine(conn, prefix)
}

// func setFlagSeen(conn *tls.Conn, num int) []string {
// 	prefix := "a0008"
// 	fmt.Println("Client->", Green(prefix+" store "+strconv.Itoa(num)+" +FLAGS \\Seen"))
// 	fmt.Fprintf(conn, "%s store %d +FLAGS \\Seen\n", prefix, num)

// 	return readBeforPrefixLine(conn, prefix)
// }

func logout(conn *tls.Conn) []string {
	prefix := "a00010"
	fmt.Println("Client->", Green(prefix+" logout"))
	fmt.Fprintf(conn, prefix+" logout\n")

	return readBeforPrefixLine(conn, prefix)
}
