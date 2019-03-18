package imap

import (
	"bytes"
	"fmt"
	"mime"
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

func parseLetter1(letter []string) (date, subject, from, to string) {
	var b strings.Builder
	for _, line := range letter {
		if !strings.HasSuffix(line, "\n") {
			fmt.Println("lol")
		}

		b.WriteString(line)
	}
	strLetter := b.String()

	var dateRegexp = regexp.MustCompile(`Date: [\w!&=/, (\-)(\|)(\.)(\?)(\+):#@]*`)
	var usefulDate = regexp.MustCompile(`[a-zA-Z]{3},  ?(\d+) [a-zA-Z]{3} \d{4} \d+:\d+`)
	dateInText := dateRegexp.FindString(strLetter)
	date = usefulDate.FindString(dateInText)

	var subjectRegexp = regexp.MustCompile(`Subject: [\w!&=/, "<>(\-)(\|)(\.)(\?)(\+):#@]*`)
	subjectInText := subjectRegexp.FindString(strLetter)
	if len(subjectInText) > 9 {
		subject = utf8Decoder(subjectInText[9:])
	}

	var fromRegexp = regexp.MustCompile(`From: [\w!&=/, "<>(\-)(\|)(\.)(\?)(\+):#@]*`)
	fromInText := fromRegexp.FindString(strLetter)
	if len(fromInText) > 6 {
		from = utf8Decoder(fromInText[6:])
	}

	var toRegexp = regexp.MustCompile(`To: [\w!&=/, "<>(\-)(\|)(\.)(\?)(\+):#@]*`)
	toInText := toRegexp.FindString(strLetter)
	if len(toInText) > 4 {
		to = utf8Decoder(toInText[4:])
	}

	var contentTransferEncodingRegexp = regexp.MustCompile(`Content-Transfer-Encoding: [\w!&=/, "<>(\-)(\|)(\.)(\?)(\+):#@]*`)
	encodingInText := contentTransferEncodingRegexp.FindString(strLetter)
	encoding := ""
	if len(encodingInText) > 27 {
		encoding = utf8Decoder(encodingInText[27:])
	}

	fmt.Println(Brown(strLetter))
	fmt.Println(Brown("-----------------------------------------------------"))
	fmt.Println(Brown(date))
	fmt.Println(Brown(subject))
	fmt.Println(Brown(from))
	fmt.Println(Brown(to))
	fmt.Println(Brown(encoding))

	return
}

func parseLetter(letter []string) (date, subject, from, to string) {
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

	fmt.Println(Brown(msg.Header.Get("To")))

	// fmt.Println(Brown(strLetter))
	return
}
