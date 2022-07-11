module main

go 1.18

replace ktrn.com/telegram => ./telegram

require ktrn.com/telegram v0.0.0-00010101000000-000000000000

require (
	github.com/pkg/errors v0.8.1 // indirect
	gopkg.in/tucnak/telebot.v2 v2.5.0 // indirect
)
