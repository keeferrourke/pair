package cmd

import (
	"fmt"
	"os"

	"gopkg.in/urfave/cli.v1"
)

const version = "0.0.1"

var (
	// With provides the `pair with` command. Modifies the VCS author to reflect
	// the invoker and the other specified authors.
	With = cli.Command{
		Name:  "with",
		Usage: "Pair with another author.",
		Action: func(cx *cli.Context) {
			// TODO
			//vcs.SetAuthor(cfg.Read().Vsc, cfg.Read().Author)
		},
	}
	// Self provides the `pair self` command. Modifies the VCS author to reflect
	// just the invoker.
	Self = cli.Command{
		Name:    "self",
		Aliases: []string{"me"},
		Usage:   "It's just you.",
		Action: func(cx *cli.Context) {
			// TODO
			//authors := []string{}
			//vsc.SetAuthor(cfg.Read().Vsc, authors)
		},
	}
	// WhoAmI provides the `pair whoami` command. Lists who the current author
	// or set of authors is.
	WhoAmI = cli.Command{
		Name:  "whoami",
		Usage: "Who are you anyway?",
		Action: func(cx *cli.Context) {
			// TODO
		},
	}

	// Branch provides the `pair branch` command. Changes the VCS branch.
	// If provided branch name exists, changes to that branch. Otherwise,
	// a new branch is created prefixed with the author names.
	Branch = cli.Command{
		Name:    "branch",
		Aliases: []string{"b"},
		Usage:   "Checkout branch.",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:   "no-prefix",
				Usage:  "Do not prefix new branch with usernames.",
				EnvVar: "PAIR_NO_BRANCH_PREFIX",
			},
		},
		Action: func(cx *cli.Command) {
			// TODO
		},
	}
	// Config provides the `pair config` command.
	Config = cli.Command{
		Name:  "config",
		Usage: "View and create pairing configurations.",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "global, g",
				Usage: "Use global configuration.",
			},
		},
		Subcommands: []cli.Command{
			{
				Name:  "dump",
				Usage: "Dump the current config.",
				Action: func(cx *cli.Context) {
					// TODO
				},
			},
			{
				Name:  "new",
				Usage: "Interactively create new config.",
				Action: func(cx *cli.Context) {
					// TODO
				},
			},
		},
	}
)

func main() {
	cli.VersionPrinter = func(cx *cli.Context) {
		fmt.Fprintf(cx.App.Writer, "%s %s - %s",
			cx.App.Name, cx.App.Version, cx.App.Description)
	}
	app := cli.NewApp()

	app.Name = "pair"
	app.Description = `Pair programming utility.
Configures your VCS (default: git) author name to reflect multiple authors.
Based on Square's pair utility.`
	app.Version = version

	app.Commands = []cli.Command{
		With,
		Self,
		WhoAmI,
		Branch,
		Config,
	}
	app.CommandNotFound = func(c *cli.Context, command string) {
		fmt.Fprintf(c.App.Writer, "Did you read the manual? %s isn't in it.\n", command)
	}

	app.Run(os.Args)
}
