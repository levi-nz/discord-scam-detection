package main

import (
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// create new session
	session, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		panic(err)
	}

	// add message listener
	session.AddHandler(receiveMessage)
	// set intents
	session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages
	// attempt to open session
	if err = session.Open(); err != nil {
		panic(err)
	}

	// block execution until Ctrl+C (or another signal) is received
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// close session
	if err = session.Close(); err != nil {
		panic(err)
	}
}
