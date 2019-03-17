package imap

import (
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	db "../db"
	. "github.com/logrusorgru/aurora"
)

const (
	port = ":993"
)

// Envelope - struct for data from envelope of mails
type Envelope struct {
	Date    string
	Subject string
	Sender  string
	Email   string
	Seen    bool
	UID     int
}

type Letter struct {
	Text string
}

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
func GetListOfMails() (envelopes []Envelope) {
	authData := db.GetAuthData()
	conn := createConn(authData.Server)
	defer func() {
		conn.Close()
		fmt.Println("Connection successful close")
	}()
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	login(conn, authData.Login, authData.Password)
	examine(conn)
	responseFetch := fetchHeader(conn)
	responseSearch := searchUnseen(conn)
	if len(responseFetch) > 0 && len(responseSearch) > 0 {
		messages := mergeStringsToMessages(responseFetch)
		unseen := parseSearch(responseSearch)
		envelopes = extractUsefulData(messages, unseen)
	} else {
		fmt.Println("Error->", Red("NO completed"))
	}
	logout(conn)
	return
}

// GetLetter - find by UID letter and return it
func GetLetter(uid string) (letter Letter) {
	authData := db.GetAuthData()
	conn := createConn(authData.Server)
	defer func() {
		conn.Close()
		fmt.Println("Connection successful close")
	}()
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	login(conn, authData.Login, authData.Password)
	examine(conn)

	letterNum := parseSearch(findLetter(conn, uid))[0]
	letterText := fetchLetter(conn, letterNum)

	var b strings.Builder
	for _, line := range letterText {
		b.WriteString(line)
	}

	letter.Text = b.String()

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
