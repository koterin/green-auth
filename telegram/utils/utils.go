package utils

import (
	"context"
	"log"
	"time"
    "strconv"

	telebot "gopkg.in/tucnak/telebot.v2"
)

type Recipient struct {
    ID int
}

func (user Recipient) Recipient() string {
    return strconv.Itoa(user.ID)
}

func StartTelegramBot(ctx context.Context, TG_BOT_KEY string) {
	settings := telebot.Settings {
		Token: TG_BOT_KEY,
		Poller: &telebot.LongPoller {
			Timeout: 1 * time.Second,
		},
	}

	bot, err := telebot.NewBot(settings)
	if err != nil {
		log.Fatal(err)
	}

	bot.Handle("/start", func(m *telebot.Message) {
		if !m.Private() {
            log.Println("Error: chat is not private")
            return
        }

        log.Println("User started bot: ", m.Sender.Username)

        var userChat Recipient
        userChat.ID = int(m.Chat.ID)

        message := "Сообщи этот ID админу (Кате) для авторизации: " + userChat.Recipient()
        bot.Send(userChat, message)
	})


	go func() {
		bot.Start()
	} ()

    log.Println("Telegram Bot started")

    <-ctx.Done()

    log.Println("Telegram Bot stopped")
	bot.Stop()
}
