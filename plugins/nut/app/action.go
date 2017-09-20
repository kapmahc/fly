package app

import (
	"log/syslog"

	log "github.com/sirupsen/logrus"
	logrus_syslog "github.com/sirupsen/logrus/hooks/syslog"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

// Action wrapper action: load viper config first
func Action(f cli.ActionFunc) cli.ActionFunc {
	viper.SetEnvPrefix("fly")
	viper.BindEnv("env")

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	return func(c *cli.Context) error {

		if err := viper.ReadInConfig(); err != nil {
			return err
		}

		if IsProduction() {
			// ----------
			log.SetLevel(log.InfoLevel)
			if wrt, err := syslog.New(syslog.LOG_INFO, Name()); err == nil {
				log.AddHook(&logrus_syslog.SyslogHook{Writer: wrt})
			} else {
				log.Error(err)
			}
		} else {
			log.SetLevel(log.DebugLevel)
		}

		log.Infof("read config from %s", viper.ConfigFileUsed())
		return f(c)
	}
}

func init() {
	viper.SetDefault("env", "development")
}
