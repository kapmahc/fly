package cache

import (
	"fmt"

	"github.com/kapmahc/fly/plugins/nut/app"
	"github.com/urfave/cli"
)

func init() {
	app.RegisterCommand(cli.Command{
		Name:    "cache",
		Aliases: []string{"c"},
		Usage:   "cache operations",
		Subcommands: []cli.Command{
			{
				Name:    "list",
				Usage:   "list all cache keys",
				Aliases: []string{"l"},
				Action: app.Action(func(_ *cli.Context) error {
					if err := app.Open(); err != nil {
						return err
					}
					keys, err := List()
					if err != nil {
						return err
					}
					for _, v := range keys {
						fmt.Println(v)
					}
					return nil
				}),
			},
			{
				Name:    "clear",
				Usage:   "clear cache items",
				Aliases: []string{"c"},
				Action: app.Action(func(_ *cli.Context) error {
					return Flush()
				}),
			},
		},
	})
}
