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

func examineBox(conn *tls.Conn, selectedBox string) []string {
	prefix := "a0002"
	fmt.Println("Client->", Green(prefix+" examine "+selectedBox))
	fmt.Fprintf(conn, "%s examine %s\n", prefix, selectedBox)

	return readBeforPrefixLine(conn, prefix)
}

func selectInbox(conn *tls.Conn, selectedBox string) []string {
	prefix := "a0003"
	fmt.Println("Client->", Green(prefix+" select "+selectedBox))
	fmt.Fprintf(conn, "%s select %s\n", prefix, selectedBox)

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

func getListOfBoxes(conn *tls.Conn) []string {
	prefix := "a0008"
	fmt.Println("Client->", Green(prefix+` LIST "" "%" `))
	fmt.Fprintf(conn, "%s LIST \"\" \"%%\" \n", prefix)

	return readBeforPrefixLine(conn, prefix)
}

func logout(conn *tls.Conn) []string {
	prefix := "a00010"
	fmt.Println("Client->", Green(prefix+" logout"))
	fmt.Fprintf(conn, prefix+" logout\n")

	return readBeforPrefixLine(conn, prefix)
}
