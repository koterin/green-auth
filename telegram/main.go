package main

import (
    "context"
    "os"
    "os/signal"

    "github.com/alexflint/go-arg"
    "telegram/utils"
)

var args struct {
    TG_BOT_KEY string `arg:"env"`
}

func main() {
    p := arg.MustParse(&args)
    if args.TG_BOT_KEY == "" {
        p.Fail("Error: TG_BOT_KEY not set")
    }

    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)

    ctx, cancel := context.WithCancel(context.Background())

    go utils.StartTelegramBot(ctx, args.TG_BOT_KEY)

    <-c
    cancel()
}

