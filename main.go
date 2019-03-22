package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"gitlab.com/shitposting/discord-random/utility"

	conf "gitlab.com/shitposting/discord-random/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	//configFilePath is the path to the config file
	configFilePath string

	//config is a struct that contains all of the informations
	//that have been loaded from the config file
	config conf.Config

	//Version represents the current admin-bot version, a compile-time value
	Version string

	//Build is the git tag for the current version
	Build string

	//err is declared here for functions that return an error as the second value
	err error

	//db is a pointer to our GORM connection to the database
	db *gorm.DB

	bot *tgbotapi.BotAPI
)

func main() {

	//Loading command line flags into their variables
	setCLIParams()

	//Reads the config files and returns an appropriate struct
	config, err = conf.ReadConfig(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	//telegram bot
	bot, err = tgbotapi.NewBotAPI(config.TelegramTokenBot)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = false

	// setting up database connection
	db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4,utf8&parseTime=True", config.DatabaseUsername, config.DatabasePassword, config.DatabaseAddress, config.DatabaseName))
	if err != nil {
		log.Fatal(err)
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + config.DiscordTokenBot)
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

		fileid := utility.GetRandomFileID(db)

		path, err := utility.GetFile(bot, fileid)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Failed, rip")
			return
		}

		pic, err := os.Open(path)
		if err != nil {
			return
		}
		s.ChannelFileSend(m.ChannelID, path, pic)
		os.Remove(path)
	}
}

func setCLIParams() {
	flag.StringVar(&configFilePath, "config", "./config.toml", "configuration file path")
	flag.Parse()
}
