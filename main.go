package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/nlopes/slack"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "config, c",
		},
	}

	app.Action = func(c *cli.Context) error {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		slack_token := os.Getenv("SLACK_TOKEN")

		api := slack.New(slack_token)
		channels, err := api.GetChannels(true)
		if err != nil {
			return err
		}
		for _, channel := range channels {
			fmt.Printf("%s\n", channel.Name)
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
