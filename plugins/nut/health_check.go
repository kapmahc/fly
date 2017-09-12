package nut

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/session"
	"github.com/astaxie/beego/session/redis"
	"github.com/astaxie/beego/toolbox"
	"github.com/streadway/amqp"
)

type databaseHealthCheck struct {
}

func (p *databaseHealthCheck) Check() error {
	o, e := orm.GetDB()
	if e != nil {
		return e
	}
	return o.Ping()
}

type sessionHealthCheck struct {
}

func (p *sessionHealthCheck) Check() error {
	switch beego.AppConfig.String("sessionprovider") {
	case "redis":
		var prv session.Provider = &redis.Provider{}
		return prv.SessionInit(3600, beego.AppConfig.String("sessionproviderconfig"))
	}
	return nil
}

type cacheHealthCheck struct {
}

func (p *cacheHealthCheck) Check() error {
	cm, err := cache.NewCache(
		beego.AppConfig.String("cacheprovider"),
		beego.AppConfig.String("cacheproviderconfig"),
	)
	if err != nil {
		return err
	}
	return cm.Put("ping", "pong", 5*time.Second)
}

type amqpHealthCheck struct {
}

func (p *amqpHealthCheck) Check() error {
	return JOBBER().open(func(ch *amqp.Channel) error {
		return nil
	})
}

func init() {
	toolbox.AddHealthCheck("database", &databaseHealthCheck{})
	toolbox.AddHealthCheck("session", &sessionHealthCheck{})
	toolbox.AddHealthCheck("cache", &cacheHealthCheck{})
	toolbox.AddHealthCheck("amqp", &amqpHealthCheck{})
}
