package main

import (
	"bufio"
	"fmt"
	"github.com/thienly/todo_service_new/cli/internal/auth"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"
	"log"
	"os"
	"syscall"
)

var (
	tokeGetter = auth.NewTokeGetter()
)
func main() {
	var tokenData *auth.TokenResponse
	app := &cli.App{Commands: []*cli.Command{{
		Name:  "login",
		Usage: "Use to log in the system",
		After: func(ctx *cli.Context) error {
			fmt.Println("Please provide your user name:")
			reader := bufio.NewReader(os.Stdin)
			userName, _:= reader.ReadString('\n')
			fmt.Println("Please provide your password:")
			password, _ := term.ReadPassword(int(syscall.Stdin))
			data, err := tokeGetter.Get(userName, string(password))
			if err != nil {
				return err
			}
			tokenData = data
			fmt.Println("You've successfully logged in!!")
			return nil
		},
	}}}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
