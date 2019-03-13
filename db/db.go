package db

import (
	"encoding/json"
	"os"
)

const (
	fileName = "auth.json"
)

// AuthData - structure which describe authentication data which user entered
// Login : login of user
// Password : password of user
// Server : IMAP server which was chosen by user
type AuthData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Server   string `json:"server"`
}

// GetAuthData open file 'auth.json' if exist
// and return data from it
// if not exit create it and put void auth data
func GetAuthData() AuthData {
	fi, err := os.Open(fileName)
	if err != nil {
		createFileWithVoidAuth()
		ad := AuthData{}
		return ad
	}
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()

	var data AuthData
	jsonParser := json.NewDecoder(fi)

	if err = jsonParser.Decode(&data); err != nil {
		panic(err.Error())
	}
	return data
}

func createFileWithVoidAuth() {
	voidAuth := map[string]string{
		"login":    "",
		"password": "",
		"server":   ""}

	fo, err := os.Create("auth.json")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	enc := json.NewEncoder(fo)
	enc.Encode(voidAuth)
}

// SaveAuthData - save authentication data to auth.json file
// data : input value
func SaveAuthData(data AuthData) {
	fo, err := os.Create("auth.json")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	enc := json.NewEncoder(fo)
	enc.Encode(data)
}
