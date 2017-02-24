package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/afrolovskiy/udoit/store"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq" // for postgres
)

const (
	startCmd = "start" // default startup command
	addCmd   = "add"
	listCmd  = "list"
	pingCmd  = "ping" // special cmd for dev mode
)

func getUpdatesChan(bot *tgbotapi.BotAPI) (tgbotapi.UpdatesChannel, error) {
	webhook := os.Getenv("UDOIT_WEBHOOK") != ""

	if webhook {
		webhookPath := "/webhook/" + bot.Token
		url := strings.TrimSuffix(os.Getenv("UDOIT_BASE_URL"), "/")

		_, err := bot.SetWebhook(tgbotapi.NewWebhook(url + webhookPath))
		if err != nil {
			log.Panicf("failed to set webhook: %s", err)
		}

		updates := bot.ListenForWebhook(webhookPath)

		go http.ListenAndServe(":"+os.Getenv("PORT"), nil)

		return updates, nil
	}

	return bot.GetUpdatesChan(tgbotapi.UpdateConfig{Timeout: 60})
}

func main() {
	dbc, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}
	dbc.SetMaxOpenConns(20) // because heroku limits

	bot, err := tgbotapi.NewBotAPI(os.Getenv("UDOIT_API_TOKEN"))
	if err != nil {
		log.Fatalf("failed to init bot api: %s", err)
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
		message := update.Message
		log.Printf("[%s] %s", message.From.UserName, message.Text)

		var msg *tgbotapi.MessageConfig

		switch cmd := message.Command(); cmd {
		case addCmd:
			descr := strings.TrimSpace(message.CommandArguments())

			t, err := store.CreateTask(dbc, descr, message.From.ID, message.Chat.ID)
			if err != nil {
				log.Fatalf("failed to add task: %s", err)
			}
			log.Printf("created task: %#v", t)

		case listCmd:
			tasks, err := store.ListTasks(dbc, message.Chat.ID)
			if err != nil {
				log.Fatalf("failed to get tasks: %s", err)
			}

			descrs := make([]string, 0, len(tasks))
			for _, t := range tasks {
				taskStr := fmt.Sprintf("%d", t.IDinchat) + " " + t.Description
				descrs = append(descrs, taskStr)
			}

			tmp := tgbotapi.NewMessage(message.Chat.ID, strings.Join(descrs, "\n"))
			msg = &tmp

		case startCmd:
			tmp := tgbotapi.NewMessage(message.Chat.ID, "Hello! I am \"U do it\" bot")
			msg = &tmp

		case pingCmd:
			tmp := tgbotapi.NewMessage(message.Chat.ID, "I'm OK!")
			msg = &tmp

		default:
			continue
		}

		// send msg at the end

		if msg != nil {
			if _, err = bot.Send(msg); err != nil {
				log.Printf("failed to send message: %s", err)
			}
		}

	}
}
