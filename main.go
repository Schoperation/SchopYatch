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

	schopYatch, err := bot.NewSchopYatchBot(config, bot.SchopYatchVersion)
	if err != nil {
		log.Fatalf("Error building SchopYatch: %v", err)
	}

	defer schopYatch.LavalinkClient.Close()
	defer schopYatch.Client.Close(context.TODO())

	err = schopYatch.Client.OpenGateway(context.TODO())
	if err != nil {
		log.Fatalf("Error connecting to Discord gateway: %v", err)
	}

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s
}
