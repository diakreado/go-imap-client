package imap

import (
	"crypto/tls"
	"fmt"
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

// Letter - struct which discribe content of email
type Letter struct {
	Date,
	Subject,
	From,
	To,
	Body string
}

// Boxes - name of boxes and one selected
type Boxes struct {
	BoxesNames  []string
	SelectedBox int
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
func GetListOfMails(selectedBox string) (envelopes []Envelope, listOfBoxes Boxes) {
	authData := db.GetAuthData()
	conn := createConn(authData.Server)
	defer func() {
		conn.Close()
		fmt.Println("Connection successful close")
	}()
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	login(conn, authData.Login, authData.Password)
	examineBox(conn, selectedBox)
	responseFetch := fetchHeader(conn)
	responseSearch := searchUnseen(conn)
	if len(responseFetch) > 0 && len(responseSearch) > 0 {
		messages := mergeStringsToMessages(responseFetch)
		unseen := parseSearch(responseSearch)
		envelopes = extractUsefulData(messages, unseen)
	} else {
		fmt.Println("Error->", Red("NO completed"))
	}
	listOfBoxes.BoxesNames = parseBoxList(getListOfBoxes(conn))
	for i, name := range listOfBoxes.BoxesNames {
		if name == selectedBox {
			listOfBoxes.SelectedBox = i
		}
	}

	logout(conn)
	return
}

// GetLetter - find by UID letter and return it
func GetLetter(uid string, selectedBox string) (letter Letter) {
	authData := db.GetAuthData()
	conn := createConn(authData.Server)
	defer func() {
		conn.Close()
		fmt.Println("Connection successful close")
	}()
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	login(conn, authData.Login, authData.Password)
	selectInbox(conn, selectedBox)

	letterNum := parseSearch(findLetter(conn, uid))[0]
	letterText := fetchLetter(conn, letterNum)

	letter.Date, letter.Subject, letter.From, letter.To, letter.Body = parseLetter(letterText)

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
