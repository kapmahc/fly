package nut

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-pg/pg"
	"github.com/kapmahc/fly/plugins/nut/app"
	"github.com/spf13/viper"
)

var orm *pg.DB

// DB database
func DB() *pg.DB {
	return orm
}

func openDb() error {
	args := viper.GetStringMap("postgresql")
	opt := pg.Options{
		Addr:     fmt.Sprintf("%s:%d", args["host"], args["port"]),
		Database: args["dbname"].(string),
		User:     args["user"].(string),
		Password: args["password"].(string),
	}
	switch args["sslmode"].(string) {
	case "allow":
		fallthrough
	case "prefer":
		opt.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	case "disable":
		opt.TLSConfig = nil
	default:
		return errors.New("pg sslmode is not supported")
	}

	db := pg.Connect(&opt)
	if !app.IsProduction() {
		db.OnQueryProcessed(func(evt *pg.QueryProcessedEvent) {
			qry, err := evt.FormattedQuery()
			if err != nil {
				panic(err)
			}
			log.Printf("%s %s", time.Since(evt.StartTime), qry)
		})
	}
	if _, err := db.Exec("SELECT NOW()"); err != nil {
		return err
	}
	orm = db
	return nil
}

func init() {
	app.RegisterResource(openDb)
}
