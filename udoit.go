package main

import (
	"log"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func getUpdatesChan(bot *tgbotapi.BotAPI) (tgbotapi.UpdatesChannel, error) {
	webhook := os.Getenv("UDOIT_WEBHOOK") != ""

	if webhook {
		_, err := bot.SetWebhook(tgbotapi.NewWebhook("")) // todo url
		if err != nil {
			log.Panicf("failed to set webhook: %s", err)
		}

		return bot.ListenForWebhook("/" + bot.Token), nil
	}

	return bot.GetUpdatesChan(tgbotapi.UpdateConfig{Timeout: 60})
}

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("UDOIT_API_TOKEN"))
	if err != nil {
		log.Panicf("failed to init bot api: %s", err)
	}

	log.Printf("authorized on account %s", bot.Self.UserName)
	bot.Debug = true

	updates, err := getUpdatesChan(bot)
	if err != nil {
		log.Panicf("failed to get updates: %s", err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID
		_, err = bot.Send(msg)
		if err != nil {
			log.Printf("failed to send message: %s", err)
		}
	}
}
