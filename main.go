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
	"gitlab.com/shitposting/memesapi/rest/client"
)

func main() {

	if err := envSetup(); err != nil {
		log.Fatal(xerrors.Errorf("cannot detect env: %w", err))
	}

	// Initialize memes api client
	mClient = client.New(apiEndpoint, apiPlatform)

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
	fmt.Println("Bot is now running.")

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

		meme, err := mClient.Random("", "", "", m.Author.ID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "could not get random meme")
		}

		var memeFile *os.File

		if strings.HasPrefix(meme.Data.URL, "https") { // Download meme and open it
			err := downloadFile(meme.Data.Filename, meme.Data.URL)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "unable to save meme")
				log.Fatal(err)
			}
			memeFile, err = os.Open(meme.Data.Filename)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "unable to open meme")
			}

			defer os.Remove(meme.Data.Filename)
		} else { // if no https prefix we have a local path
			memeFile, err = os.Open(meme.Data.URL)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "unable to open meme")
			}
		}

		_, err = s.ChannelFileSend(m.ChannelID, meme.Data.Filename, memeFile)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "can't send meme")
		}

		defer memeFile.Close()
	}
}
