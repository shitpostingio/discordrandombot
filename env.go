package main

import (
	"log"
	"os"

	"github.com/shitpostingio/randomapi/rest/client"
)

var (
	discordToken string

	apiEndpoint string

	apiPlatform string

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

	apiEndpoint, ok = os.LookupEnv("API_ENDPOINT")
	if apiEndpoint == "" || !ok {
		apiEndpoint = "http://localhost:34378"
	}

	apiPlatform, ok = os.LookupEnv("API_PLATFORM")
	if apiPlatform == "" || !ok {
		apiPlatform = "discordrandom"
	}

	return nil
}
