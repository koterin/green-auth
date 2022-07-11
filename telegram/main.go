package main

import (
    "context"
	"os"
	"os/signal"

    "ktrn.com/telegram"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go telegram.StartTelegramBot(ctx)

	<-c
	cancel()
}

