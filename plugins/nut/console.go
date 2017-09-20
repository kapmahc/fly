package nut

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

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
	if !app.IsProduction() {
		return Router().Run(addr)
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: Router(),
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
	app.RegisterCommand(cli.Command{
		Name:    "routes",
		Aliases: []string{"rt"},
		Usage:   "print out all defined routes",
		Action: func(*cli.Context) error {
			if err := openRouter(); err != nil {
				return err
			}

			tpl := "%-7s %s\n"
			fmt.Printf(tpl, "METHOD", "PATH")
			for _, r := range Router().Routes() {
				fmt.Printf(tpl, r.Method, r.Path)
			}
			return nil
		},
	})
}
