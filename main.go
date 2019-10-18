package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"golang.org/x/xerrors"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/shitposting/discord-random/utility"
	"gitlab.com/shitposting/memesapi/rest/client"
)

func main() {

	if err := envSetup(); err != nil {
		log.Fatal(xerrors.Errorf("cannot detect env: %w", err))
	}

	// Initialize memes api client
	mClient = client.New(memesEndpoint)

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(handleMessages)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Printf("Bot is now running.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func handleMessages(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.ToLower(m.Content) == "random" {

		meme, err := mClient.Random()
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "error")
		}

		err = utility.DownloadFile(meme.Data.Filename, meme.Data.URL)
		if err != nil {
			log.Fatal(err)
		}

		file, err := os.Open(meme.Data.Filename)
		if err != nil {
			log.Fatal(err)
		}

		s.ChannelFileSend(m.ChannelID, meme.Data.Filename, file)

		os.Remove(meme.Data.Filename)
	}
}
