package main

import (
	"os"
)

const (
	Nick = "build-bot"
	Chan = "#dashboard"
)

var (
	Username    = os.Getenv("SASL_USER")
	Password    = os.Getenv("SASL_PASSWORD")
	Server      = os.Getenv("SERVER")
	VerifyTLS   = os.Getenv("VERIFY_TLS") == "true"
	StorageFile = os.Getenv("STORAGE_FILE")
	Timezone    = os.Getenv("TZ")
)

func main() {
	c, err := New(Username, Password, Server, VerifyTLS, StorageFile, Timezone)
	if err != nil {
		panic(err)
	}

	panic(c.bottom.Client.Connect())
}
