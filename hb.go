package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "hb"
	app.Usage = "Perform tasks in relation to HobbyFarm"
	app.Commands = []*cli.Command{
		dlCommand(),
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
