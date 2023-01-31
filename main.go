package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"schoperation/schopyatch/bot"
	"syscall"
)

func main() {
	log.SetPrefix("SY|")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config, err := bot.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading the config: %v", err)
	}

	/*
		TODO cmds:

		- help (reads descriptions, usages, and aliases)
		about

		- join voice channel (to keep functionality in one place)
		- leave voice channel

		play
			support for
				- search
				- playlists
		- pause
		- resume
		seek
		- nowplaying
		loop track/queue

		- queue
		- skip
		- skipto
		- clear
		- shuffle
	*/
	schopYatch := bot.NewSchopYatchBot(config)
	err = schopYatch.SetupClient()
	if err != nil {
		log.Fatalf("Error building SchopYatch: %v", err)
	}

	err = schopYatch.SetupLavalink()
	if err != nil {
		log.Fatalf("Error setting up lavalink: %v", err)
	}

	defer schopYatch.Lavalink.Close()
	defer schopYatch.Client.Close(context.TODO())

	err = schopYatch.Client.OpenGateway(context.TODO())
	if err != nil {
		log.Fatalf("Error connecting to Discord gateway: %v", err)
	}

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s
}
