package main

import (
	"log"

	_ "github.com/kapmahc/fly/plugins/erp"
	_ "github.com/kapmahc/fly/plugins/forum"
	_ "github.com/kapmahc/fly/plugins/mall"
	_ "github.com/kapmahc/fly/plugins/nut"
	"github.com/kapmahc/fly/plugins/nut/app"
	_ "github.com/kapmahc/fly/plugins/ops/mail"
	_ "github.com/kapmahc/fly/plugins/ops/vpn"
	_ "github.com/kapmahc/fly/plugins/pos"
	_ "github.com/kapmahc/fly/plugins/reading"
	_ "github.com/kapmahc/fly/plugins/survey"
)

func main() {
	if err := app.Main(); err != nil {
		log.Panic(err)
	}
}
