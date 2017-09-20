package main

import (
	"log"

	"github.com/kapmahc/fly/plugins/nut/app"
)

func main() {
	if err := app.Main(); err != nil {
		log.Panic(err)
	}
}
