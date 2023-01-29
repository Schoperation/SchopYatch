package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"schoperation/schopyatch/bot"
	"schoperation/schopyatch/command"
	"syscall"
)

func main() {
	log.SetPrefix("SY|")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config, err := bot.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading the config: %v", err)
	}

	commands := []command.Command{
		command.NewPingCmd(),
		command.NewPlayCmd(),
	}

	schopYatch := bot.NewSchopYatchBot(config, commands)
	err = schopYatch.SetupClient()
	if err != nil {
		log.Fatalf("Error building SchopYatch: %v", err)
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
