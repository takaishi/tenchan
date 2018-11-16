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
	Name     string `toml:"name"`
	Match    string `toml:"match"`
	Channels []slack.Channel
}

func filterChannel(api *slack.Client, config Config) ([]CType, error) {
	channels, err := api.GetChannels(true)

	if err != nil {
		return nil, err
	}

	for i, ctype := range config.CTypes {
		r := regexp.MustCompile(ctype.Match)
		for _, channel := range channels {
			if r.Match([]byte(channel.Name)) {
				config.CTypes[i].Channels = append(config.CTypes[i].Channels, channel)
			}
		}
	}

	return config.CTypes, nil
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
		slack_channel := os.Getenv("SLACK_CHANNEL_NAME")

		var config Config
		_, err = toml.DecodeFile(c.String("config"), &config)
		if err != nil {
			return err
		}

		api := slack.New(slack_token)
		team, err := api.GetTeamInfo()
		if err != nil {
			return err
		}

		ctypes, err := filterChannel(api, config)
		if err != nil {
			return err
		}

		params := slack.PostMessageParameters{}
		_, _, err = api.PostMessage(slack_channel, "一時的に作成されたチャンネル一覧だよ〜", params)
		if err != nil {
			return err
		}
		for _, ctype := range ctypes {
			params := slack.PostMessageParameters{}
			for _, channel := range ctype.Channels {
				attachment := slack.Attachment{
					Title:     fmt.Sprintf("#%s", channel.Name),
					TitleLink: fmt.Sprintf("slack://channel?team=%s&id=%s", team.ID, channel.ID),
					Text:      channel.Topic.Value,
					Color:     "#36a64f",
				}
				params.Attachments = append(params.Attachments, attachment)
			}
			_, _, err := api.PostMessage(slack_channel, ctype.Name, params)
			if err != nil {
				return err
			}
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
