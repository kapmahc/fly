package nut

import (
	"crypto/x509/pkix"
	"fmt"
	"os"
	"path"
	"text/template"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/kapmahc/fly/plugins/nut/app"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
	"golang.org/x/text/language"
)

func generateConfig(c *cli.Context) error {
	const fn = "config.toml"
	if _, err := os.Stat(fn); err == nil {
		return fmt.Errorf("file %s already exists", fn)
	}
	fmt.Printf("generate file %s\n", fn)

	viper.Set("env", c.String("environment"))
	args := viper.AllSettings()
	fd, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer fd.Close()
	end := toml.NewEncoder(fd)
	err = end.Encode(args)

	return err

}

func generateMigration(c *cli.Context) error {
	name := c.String("name")
	if len(name) == 0 {
		cli.ShowCommandHelp(c, "migration")
		return nil
	}
	const pkg = "migrations"
	version := time.Now().Format("20060102150405000")
	root := path.Join("db", pkg)
	if err := os.MkdirAll(root, 0700); err != nil {
		return err
	}
	file := path.Join(root, fmt.Sprintf("%s_%s.go", version, name))
	fmt.Printf("generate file %s\n", file)
	fd, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return err
	}
	defer fd.Close()

	tpl, err := template.ParseFiles(path.Join("templates", "nginx.conf"))
	if err != nil {
		return err
	}

	return tpl.Execute(fd, struct {
		Name    string
		Version string
	}{
		Name:    name,
		Version: version,
	})
}

func generateLocale(c *cli.Context) error {
	name := c.String("name")
	if len(name) == 0 {
		cli.ShowCommandHelp(c, "locale")
		return nil
	}
	lng, err := language.Parse(name)
	if err != nil {
		return err
	}
	const root = "locales"
	if err = os.MkdirAll(root, 0700); err != nil {
		return err
	}
	file := path.Join(root, fmt.Sprintf("%s.ini", lng.String()))
	fmt.Printf("generate file %s\n", file)
	fd, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer fd.Close()
	return err
}

func generateNginxConf(*cli.Context) error {
	if err := app.Open(); err != nil {
		return err
	}
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	name := app.Name()
	fn := path.Join("tmp", "nginx.conf")
	if err = os.MkdirAll(path.Dir(fn), 0700); err != nil {
		return err
	}
	fmt.Printf("generate file %s\n", fn)
	fd, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer fd.Close()

	tpl, err := template.ParseFiles(path.Join("templates", "nginx.conf"))
	if err != nil {
		return err
	}

	return tpl.Execute(fd, struct {
		Port int
		Root string
		Name string
		Ssl  bool
	}{
		Name: name,
		Port: viper.GetInt("server.port"),
		Root: pwd,
		Ssl:  viper.GetBool("server.ssl"),
	})
}

func generateSsl(c *cli.Context) error {
	name := c.String("name")
	if len(name) == 0 {
		cli.ShowCommandHelp(c, "openssl")
		return nil
	}
	root := path.Join("etc", "ssl", name)

	key, crt, err := CreateCertificate(
		true,
		pkix.Name{
			Country:      []string{c.String("country")},
			Organization: []string{c.String("organization")},
		},
		c.Int("years"),
	)
	if err != nil {
		return err
	}

	fnk := path.Join(root, "key.pem")
	fnc := path.Join(root, "crt.pem")

	fmt.Printf("generate pem file %s\n", fnk)
	err = WritePemFile(fnk, "RSA PRIVATE KEY", key, 0600)
	fmt.Printf("test: openssl rsa -noout -text -in %s\n", fnk)

	if err == nil {
		fmt.Printf("generate pem file %s\n", fnc)
		err = WritePemFile(fnc, "CERTIFICATE", crt, 0444)
		fmt.Printf("test: openssl x509 -noout -text -in %s\n", fnc)
	}
	if err == nil {
		fmt.Printf(
			"verify: diff <(openssl rsa -noout -modulus -in %s) <(openssl x509 -noout -modulus -in %s)",
			fnk,
			fnc,
		)
	}
	fmt.Println()
	return err
}

func init() {
	app.RegisterCommand(cli.Command{
		Name:    "generate",
		Aliases: []string{"g"},
		Usage:   "generate file template",
		Subcommands: []cli.Command{
			{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "generate config file",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "environment, e",
						Value: "development",
						Usage: "environment, like: development, production, stage, test...",
					},
				},
				Action: generateConfig,
			},
			{
				Name:    "nginx",
				Aliases: []string{"ng"},
				Usage:   "generate nginx.conf",
				Action:  app.Action(generateNginxConf),
			},
			{
				Name:    "migration",
				Usage:   "generate migration file",
				Aliases: []string{"m"},
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "name, n",
						Usage: "name",
					},
				},
				Action: generateMigration,
			},
			{
				Name:    "openssl",
				Aliases: []string{"ssl"},
				Usage:   "generate ssl certificates",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "name, n",
						Usage: "name",
					},
					cli.StringFlag{
						Name:  "country, c",
						Value: "Earth",
						Usage: "country",
					},
					cli.StringFlag{
						Name:  "organization, o",
						Value: "Mother Nature",
						Usage: "organization",
					},
					cli.IntFlag{
						Name:  "years, y",
						Value: 1,
						Usage: "years",
					},
				},
				Action: generateSsl,
			},
			{
				Name:    "locale",
				Usage:   "generate locale file",
				Aliases: []string{"l"},
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "name, n",
						Usage: "locale name",
					},
				},
				Action: generateLocale,
			},
		},
	})
}
