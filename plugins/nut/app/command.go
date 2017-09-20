package app

import "github.com/urfave/cli"

var commands []cli.Command

// RegisterCommand register command
func RegisterCommand(args ...cli.Command) {
	commands = append(commands, args...)
}
