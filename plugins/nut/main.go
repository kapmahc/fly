package nut

import (
	"crypto/x509/pkix"
	"database/sql"
	"fmt"
	"html/template"
	"os"
	"path"
	"time"

	"golang.org/x/text/language"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/toolbox"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	"github.com/urfave/cli"
)

// Main entry
func Main(args ...string) error {
	ap := cli.NewApp()
	ap.Name = args[0]
	ap.Version = fmt.Sprintf("%s (%s)", Version, BuildTime)
	ap.Authors = []cli.Author{
		cli.Author{
			Name:  AuthorName,
			Email: AuthorEmail,
		},
	}

	ts, err := time.Parse(time.RFC1123Z, BuildTime)
	if err != nil {
		return err
	}
	ap.Compiled = ts

	ap.Copyright = Copyright
	ap.Usage = Usage
	ap.EnableBashCompletion = true
	ap.Commands = []cli.Command{
		{
			Name:    "database",
			Aliases: []string{"db"},
			Usage:   "database operations",
			Subcommands: []cli.Command{
				{
					Name:    "migrate",
					Usage:   "migrate the DB to the most recent version available",
					Aliases: []string{"m"},
					Action: migrateAction(func(_ *cli.Context, mig *migrate.Migrate) error {
						return mig.Up()
					}),
				},
				{
					Name:    "rollback",
					Usage:   "roll back the version by 1",
					Aliases: []string{"r"},
					Action: migrateAction(func(_ *cli.Context, mig *migrate.Migrate) error {
						return mig.Steps(-1)
					}),
				},
				{
					Name:    "version",
					Usage:   "print current migration version",
					Aliases: []string{"v"},
					Action: migrateAction(func(_ *cli.Context, mig *migrate.Migrate) error {
						ver, dirty, err := mig.Version()
						if err != nil {
							return err
						}
						if dirty {
							fmt.Printf("%d(dirty)\n", ver)
						} else {
							fmt.Println(ver)
						}
						return nil
					}),
				},
				{
					Name:    "drop",
					Usage:   "drop everyting inside database",
					Aliases: []string{"c"},
					Action: migrateAction(func(_ *cli.Context, mig *migrate.Migrate) error {
						return mig.Drop()
					}),
				},
			},
		},
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "generate file template",
			Subcommands: []cli.Command{
				{
					Name:    "nginx",
					Aliases: []string{"ng"},
					Usage:   "generate nginx.conf",
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "ssl",
							Usage: "https?",
						},
					},
					Action: generateNginxConf,
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
					Name:    "migration",
					Usage:   "generate migration file",
					Aliases: []string{"m"},
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name, n",
							Usage: "name",
						},
					},
					Action: func(c *cli.Context) error {
						name := c.String("name")
						if len(name) == 0 {
							cli.ShowCommandHelp(c, "migration")
							return nil
						}
						now := time.Now()
						for _, act := range []string{"up", "down"} {
							if err := generateMigration(now.Unix(), name, act); err != nil {
								return err
							}
						}
						return nil
					},
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
		},
	}
	ap.Action = mainAction

	return ap.Run(args)
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
func generateMigration(version int64, name, act string) error {
	file := path.Join("db", "migrate", fmt.Sprintf("%d_%s.%s.sql", version, name, act))
	fmt.Printf("generate file %s\n", file)
	fd, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer fd.Close()
	return nil
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
func generateNginxConf(c *cli.Context) error {
	tpl, err := template.ParseFiles(path.Join("templates", "nginx.conf"))
	if err != nil {
		return err
	}
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	fn := path.Join("tmp", "nginx.conf")
	fmt.Printf("generate file %s\n", fn)
	fd, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer fd.Close()

	cfg := beego.BConfig
	return tpl.Execute(
		fd,
		map[string]interface{}{
			"name":  cfg.ServerName,
			"port":  cfg.Listen.HTTPPort,
			"theme": cfg.AppName,
			"root":  pwd,
			"ssl":   c.Bool("ssl"),
		})
}

func migrateAction(fun func(*cli.Context, *migrate.Migrate) error) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		drv := beego.AppConfig.String("databasedriver")
		src := beego.AppConfig.String("databasesource")
		db, err := sql.Open(drv, src)
		if err != nil {
			return err
		}
		defer db.Close()

		ins, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			return err
		}

		pwd, err := os.Getwd()
		if err != nil {
			return err
		}
		mig, err := migrate.NewWithDatabaseInstance(
			"file://"+path.Join(pwd, "db", "migrate"),
			drv, ins)
		if err != nil {
			return err
		}
		if err := fun(ctx, mig); err != nil {
			return err
		}

		fmt.Println("Done!")
		return nil
	}
}

func mainAction(_ *cli.Context) error {
	logs.SetLogger(logs.AdapterConsole)
	logs.SetLogger(logs.AdapterFile, `{"filename":"`+path.Join("tmp", "www.log")+`"}`)

	orm.Debug = beego.BConfig.RunMode != beego.PROD
	orm.RegisterDataBase(
		"default",
		beego.AppConfig.String("databasedriver"),
		beego.AppConfig.String("databasesource"),
	)

	if err := Open(); err != nil {
		return err
	}

	toolbox.StartTask()
	defer toolbox.StopTask()

	go func() {
		host, _ := os.Hostname()
		JOBBER().Receive(host)
	}()

	beego.ErrorController(&ErrorController{})
	beego.Run()
	return nil
}
