package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
	"github.com/nlopes/slack"
	"github.com/urfave/cli"
	"log"
	"os"
	"regexp"
)

type Config struct {
	CTypes []CType `toml:"ctype"`
}

type CType struct {
	Name  string `toml:"name"`
	Match string `toml:"match"`
}

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

		var config Config
		_, err = toml.DecodeFile(c.String("config"), &config)
		if err != nil {
			return err
		}

		api := slack.New(slack_token)
		channels, err := api.GetChannels(true)
		if err != nil {
			return err
		}
		for _, ctype := range config.CTypes {
			for _, channel := range channels {
				r := regexp.MustCompile(ctype.Match)
				if r.Match([]byte(channel.Name)) {
					fmt.Printf("%s\n", channel.Name)
				}
			}
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
