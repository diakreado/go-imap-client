package imap

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"mime"
	"mime/quotedprintable"
	"net/mail"
	"regexp"
	"strconv"
	"strings"

	. "github.com/logrusorgru/aurora"
)

func isOK(response string) (successful bool) {
	var valid = regexp.MustCompile(`^.{5} OK `)
	successful = valid.MatchString(response)
	return
}

// convert strings to messages
// some messages can contain from many strings
func mergeStringsToMessages(response []string) (result []string) {
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

func extractUsefulData(response []string, unseen []int) (result []Envelope) {
	var dateRegexp = regexp.MustCompile(`[a-zA-Z]{3},  ?(\d+) [a-zA-Z]{3} \d{4}`)
	var subjectRegexp = regexp.MustCompile(`" "(.)*" \(\(`)
	var contactRegexp = regexp.MustCompile(`" \(\(("[\w!&=/, (\-)(\|)(\.)(\?)(\+):#@]*"|NIL) ("[\w!&=/, (\-)(\|)(\.)(\?)(\+):#@]*"|NIL) ("[\w!&=/, (\-)(\|)(\.)(\?)(\+):#@]*"|NIL) ("[\w!&=/, (\-)(\|)(\.)(\?)(\+):#@]*"|NIL)\)\)`)
	var uidRegexp = regexp.MustCompile(`UID \d+`)
	for i, line := range response {
		date := dateRegexp.FindString(line)

		subject := ""
		matchesSubject := subjectRegexp.FindString(line)
		a1 := []rune(matchesSubject)
		if len(matchesSubject) > 7 {
			subject = utf8Decoder(string(a1[3 : len(matchesSubject)-4]))
		}

		matchesContact := contactRegexp.FindString(line)
		sender := ""
		email := ""
		a2 := []rune(matchesContact)
		if len(matchesContact) > 6 {
			sender, email = contactParser(string(a2[4 : len(matchesContact)-2]))
		}

		uid, err := strconv.Atoi(strings.Split(uidRegexp.FindString(line), " ")[1])
		if err != nil {
			fmt.Println("UID Atoi :", Red(err))
		}

		envelope := Envelope{date, subject, sender, email, !contains(unseen, i+1), uid}
		result = append(result, envelope)
	}
	return
}

func contactParser(contact string) (sender, email string) {
	parts := strings.Split(contact, "\"")
	var codeNil = regexp.MustCompile(`NIL`)
	var partOfContact []string
	for _, part := range parts {
		if !codeNil.MatchString(part) && part != "" && part != " " {
			partOfContact = append(partOfContact, utf8Decoder(part))
		}
	}
	switch len(partOfContact) {
	case 1:
		{
			sender = fmt.Sprintf("%s", partOfContact[0])
		}
	case 2:
		{
			sender = fmt.Sprintf("%s", partOfContact[0])
			email = fmt.Sprintf("<%s@%s>", partOfContact[0], partOfContact[1])
		}
	case 3:
		{
			sender = fmt.Sprintf("%s", partOfContact[0])
			email = fmt.Sprintf("<%s@%s>", partOfContact[1], partOfContact[2])
		}
	case 4:
		{
			sender = fmt.Sprintf("%s %s", partOfContact[0], partOfContact[1])
			email = fmt.Sprintf("<%s@%s>", partOfContact[2], partOfContact[3])
		}
	}
	return
}

func utf8Decoder(text string) (result string) {
	var codeUTF8 = regexp.MustCompile(`=\?(utf-8|UTF-8)\?[\w!&=/, (\-)(\|)(\.)(\?)(\+):#@]*\?=`)
	dec := new(mime.WordDecoder)
	if codeUTF8.MatchString(text) {
		var b strings.Builder
		parts := strings.Split(text, " ")
		for _, part := range parts {
			if part != "" && part != " " {
				if codeUTF8.MatchString(part) {
					decodedText, _ := dec.Decode(part)
					b.WriteString(decodedText)
				} else {
					b.WriteString(part)
				}
			}
		}
		result = b.String()
	} else {
		result = text
	}
	return
}

func parseSearch(response []string) (result []int) {
	if len(response) == 2 {
		var num = regexp.MustCompile(`\d+`)
		numArr := num.FindAllString(response[0], -1)
		for _, part := range numArr {
			if newNum, err := strconv.Atoi(part); err == nil {
				result = append(result, newNum)
			}
		}
	}
	return
}

func parseLetter(letter []string) (date, subject, from, to, body string) {
	var b strings.Builder
	numOfStr := len(letter)
	for i, line := range letter {
		if i > 0 && i < numOfStr-2 {
			b.WriteString(line)
		}
	}
	strLetter := b.String()

	msg, err := mail.ReadMessage(bytes.NewBuffer([]byte(strLetter)))
	if err != nil {
		panic(err)
	}

	date = msg.Header.Get("Date")
	subject = utf8Decoder(msg.Header.Get("Subject"))
	from = utf8Decoder(msg.Header.Get("From"))
	to = utf8Decoder(msg.Header.Get("To"))

	encoding := msg.Header.Get("Content-Transfer-Encoding")

	if encoding == "base64" {
		buf := new(bytes.Buffer)
		buf.ReadFrom(msg.Body)
		s := buf.String()
		decodedB64, err := base64.StdEncoding.DecodeString(s)
		if err != nil {
			fmt.Println("decodeB64 error:", err)
			return
		}
		body = string(decodedB64)
	} else {
		r := quotedprintable.NewReader(msg.Body)
		decodedQI, err := ioutil.ReadAll(r)
		if err != nil {
			fmt.Println("decodeQI error:", err)
			return
		}
		body = string(decodedQI)
	}

	return
}

func parseBoxList(response []string) (nameOfBoxes []string) {
	var b strings.Builder
	for _, line := range response {
		b.WriteString(line)
	}
	strResponse := b.String()

	var boxNameRegexp = regexp.MustCompile(`"/" [\w]+`)
	boxNames := boxNameRegexp.FindAllString(strResponse, -1)

	for _, nameInText := range boxNames {
		nameOfBoxes = append(nameOfBoxes, nameInText[4:])
	}

	return
}
