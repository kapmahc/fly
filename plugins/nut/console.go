package nut

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/csrf"
	"github.com/kapmahc/fly/plugins/nut/app"
	"github.com/kapmahc/fly/plugins/nut/health"
	"github.com/kapmahc/fly/plugins/nut/job"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

func _startBackgroundJob() {
	go func() {
		host, _ := os.Hostname()
		for {
			if err := job.Receive(host); err != nil {
				log.Error(err)
			}
		}
	}()
}

func _startHelathCheck() {
	ticker := time.NewTicker(5 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				health.Check()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func _startHTTPListen() error {
	port := viper.GetInt("server.port")
	addr := fmt.Sprintf(":%d", port)
	log.Infof(
		"application starting on http://localhost:%d",
		port,
	)
	rt, err := app.Router()
	if err != nil {
		return err
	}

	srv := &http.Server{
		Addr: addr,
		Handler: csrf.Protect(
			[]byte(viper.GetString("secret")),
			csrf.CookieName("csrf"),
			csrf.RequestHeader("Authenticity-Token"),
			csrf.FieldName("authenticity_token"),
			csrf.Secure(viper.GetBool("server.ssl")),
		)(rt),
	}

	if !app.IsProduction() {
		return srv.ListenAndServe()
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			log.Error(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("shutdown server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return err
	}
	log.Println("server exiting")
	return nil
}

func startServer(_ *cli.Context) error {
	if err := app.Open(); err != nil {
		return err
	}

	_startBackgroundJob()
	_startHelathCheck()

	return _startHTTPListen()
}

func init() {
	app.RegisterCommand(cli.Command{
		Name:    "server",
		Aliases: []string{"s"},
		Usage:   "start the app server",
		Action:  app.Action(startServer),
	})

}
