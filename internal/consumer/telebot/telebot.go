package telebot

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"strconv"
	"time"
)

var (
	Bot  *tb.Bot
	Chat *tb.Chat
)

func CreateBot(botToken string, chatID int64)  {
	var err error

	Chat = &tb.Chat{ID: chatID}
	Bot, err = tb.NewBot(tb.Settings{
		Token: botToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Println(err)
		return
	}
	Bot.Handle("/chatID", func(m *tb.Message) {
		Bot.Send(m.Chat, "chatID: " + strconv.FormatInt(m.Chat.ID, 10))
	})
}

func StartBot() {
	go Bot.Start()
}

func Send(msg string) {
	Bot.Send(Chat, msg)
}