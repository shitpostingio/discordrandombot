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

const (
	//MaxRetries is the maximum amount of attempts to
	//forward that the bot should make
	MaxRetries = 3
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

		for i := 0; i < MaxRetries; i++ {

			meme, err := mClient.Random("", "", "", m.Author.ID)
			if err != nil {
				continue
			}

			var memeFile *os.File

			if strings.HasPrefix(meme.Meme.URL, "https") { // Download meme and open it
				err := downloadFile(meme.Meme.Filename, meme.Meme.URL)
				if err != nil {
					continue
				}
				memeFile, err = os.Open(meme.Meme.Filename)
				if err != nil {
					continue
				}

				defer os.Remove(meme.Meme.Filename)
			} else { // if no https prefix we have a local path
				memeFile, err = os.Open(meme.Meme.URL)
				if err != nil {
					continue
				}
			}

			mFile := discordgo.File{
				Name:   meme.Meme.Filename,
				Reader: memeFile,
			}

			toSend := discordgo.MessageSend{
				Content: fmt.Sprintf("Source: <https://t.me/shitpost/%d>", meme.Meme.MessageID),
				File:    &mFile,
			}

			s.ChannelMessageSendComplex(m.ChannelID, &toSend)

			defer memeFile.Close()
			if err == nil {
				return
			}
		}

		s.ChannelMessageSend(m.ChannelID, "Unable to send meme, try again later")
	}
}
