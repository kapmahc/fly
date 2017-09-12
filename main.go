package main

import (
	"log"
	"os"

	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/astaxie/beego/session/redis"
	"github.com/kapmahc/fly/plugins/nut"
	_ "github.com/kapmahc/fly/routers"
	_ "github.com/lib/pq"
	_ "github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
)

func main() {
	if err := nut.Main(os.Args...); err != nil {
		log.Fatal(err)
	}
}
