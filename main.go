package main

import (
	"os"
)

const (
	Nick = "build-bot"
	Chan = "#dashboard"
)

var (
	Username   = os.Getenv("SASL_USER")
	Password   = os.Getenv("SASL_PASSWORD")
	Server     = os.Getenv("SERVER")
	VerifyTLS  = os.Getenv("VERIFY_TLS") == "true"
	StorageDir = os.Getenv("STORAGE_DIR")
	Timezone   = os.Getenv("TZ")
)

func main() {
	c, err := New(Username, Password, Server, VerifyTLS, StorageDir, Timezone)
	if err != nil {
		panic(err)
	}

	c.bottom.Client.Connect()
}
