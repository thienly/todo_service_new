package main

import (
	"bufio"
	"fmt"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"
	"log"
	"os"
	"syscall"
)

var token string

func main() {
	app := &cli.App{
		Name:  "todo",
		Usage: "todo tracking",
		Action: func(context *cli.Context) error {
			fmt.Println("boom! I say!")
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "login",
				Usage: "login",
				Action: func(context *cli.Context) error {
					buf := bufio.NewReader(os.Stdin)
					fmt.Println("user name:")
					username, err := buf.ReadBytes('\n')
					fmt.Println("password:")
					readPassword, err := term.ReadPassword(syscall.Stdin)
					if err != nil {
						fmt.Println(err)
					} else {
						fmt.Println(string(username))
						fmt.Println(string(readPassword))
						token = "abc"
					}
					return err
				},
				UsageText: "using for login to system",
				OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
					_, _ = fmt.Fprintf(c.App.Writer, "for shame\n")
					return err
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}
