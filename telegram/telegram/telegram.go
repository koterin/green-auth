package telegram

import (
	"context"
	"log"
	"time"
    "strconv"
    "os"

	telebot "gopkg.in/tucnak/telebot.v2"
)

type Recipient struct {
    ID int
}

func (id Recipient) Recipient() string {
    return strconv.Itoa(id.ID)
}

func StartTelegramBot(ctx context.Context) {
	settings := telebot.Settings {
		Token: os.Getenv("TG_BOT_KEY"),
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
			return
		}
        log.Println("User started bot: ", m.Sender.Username)

        var userChat Recipient
        userChat.ID = int(m.Chat.ID)

	message := "Сообщите ваш ID админу для авторизации: " + userChat.Recipient()
        bot.Send(userChat, message)
	})


	go func() {
		bot.Start()
	} ()

    log.Println("[Telegram] Bot started")

    <-ctx.Done()
	bot.Stop()
}
