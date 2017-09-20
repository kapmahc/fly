package app

import "github.com/urfave/cli"

var commands []cli.Command

// Register register command
func Register(args ...cli.Command) {
	commands = append(commands, args...)
}
