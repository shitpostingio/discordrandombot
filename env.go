package main

import (
	"log"
	"os"

	"gitlab.com/shitposting/memesapi/rest/client"
)

var (
	discordToken string

	memesEndpoint string

	//err is declared here for functions that return an error as the second value
	err error

	mClient *client.Client
)

func envSetup() error {
	var ok bool

	discordToken, ok = os.LookupEnv("DISCORD_TOKEN")
	if discordToken == "" || !ok {
		log.Fatalf("Discord token bot is not optional!")
	}

	memesEndpoint, ok = os.LookupEnv("API_ENDPOINT")
	if memesEndpoint == "" || !ok {
		memesEndpoint = "http://localhost:34378"
	}

	return nil
}
