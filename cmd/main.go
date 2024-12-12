package main

import (
	"log"

	"github.com/cutlery47/email-service/internal/app"
)

func main() {
	log.Fatal("error: ", app.Run())
}
