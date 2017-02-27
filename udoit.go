package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"strconv"

	"github.com/afrolovskiy/udoit/store"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq" // for postgres
)

const (
	addCmd    = "add"
	deleteCmd = "delete"
	listCmd   = "list"
	pingCmd   = "ping"  // special cmd for dev mode
	startCmd  = "start" // default startup command
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
	defer dbc.Close()
	dbc.SetMaxOpenConns(20) // because heroku limits

	bot, err := tgbotapi.NewBotAPI(os.Getenv("UDOIT_API_TOKEN"))
	if err != nil {
		log.Fatalf("failed to init bot api: %s", err)
	}
	log.Printf("authorized on account %s", bot.Self.UserName)

	go func() {
		updates, err := getUpdatesChan(bot)
		if err != nil {
			log.Panicf("failed to get updates: %s", err)
		}

		log.Print("Ready to get updates from Telegram")

		for update := range updates {
			if update.Message == nil {
				continue
			}
			message := update.Message
			log.Printf("[%s] %s", message.From.UserName, message.Text)

			var msg *tgbotapi.MessageConfig

			switch cmd := message.Command(); cmd {
			case addCmd:
				desc := strings.TrimSpace(message.CommandArguments())

				t, err := store.CreateTask(dbc, desc, message.From.ID, message.Chat.ID)
				if err != nil {
					log.Fatalf("failed to add task: %s", err)
				}
				log.Printf("created task: %#v", t)
				tmp := tgbotapi.NewMessage(message.Chat.ID, "Task "+fmt.Sprintf("%d", t.IDInChat)+" created")
				msg = &tmp

			case listCmd:
				var tasks []store.Task
				var err error
				if message.Chat.IsPrivate() {
					tasks, err = store.UserTasks(dbc, message.From.ID)
				} else {
					tasks, err = store.ChatTasks(dbc, message.Chat.ID)
				}
				if err != nil {
					log.Fatalf("failed to get tasks: %s", err)
				}

				var msgText string

				if len(tasks) > 0 {
					descs := make([]string, 0, len(tasks))
					for _, t := range tasks {
						taskStr := "#" + fmt.Sprintf("%d", t.IDInChat) + " " + t.Description
						descs = append(descs, taskStr)
					}
					msgText = strings.Join(descs, "\n")
				} else {
					msgText = "No current tasks"
				}

				tmp := tgbotapi.NewMessage(message.Chat.ID, msgText)
				msg = &tmp

			case deleteCmd:
				// todo use regexp to extract number
				arg := message.CommandArguments()

				if id, err := strconv.Atoi(arg); err != nil {
					log.Print("/delete no number")
				} else {
					if err = store.DeleteTask(dbc, message.Chat.ID, id); err != nil {
						log.Printf("failed to delete task %s", err)
					}
					tmp := tgbotapi.NewMessage(message.Chat.ID, "Task "+fmt.Sprintf("%d", id)+" deleted")
					msg = &tmp
				}

			case startCmd:
				tmp := tgbotapi.NewMessage(message.Chat.ID, "Hello! I am \"U do it\" bot")
				msg = &tmp

			case pingCmd:
				tmp := tgbotapi.NewMessage(message.Chat.ID, "I'm OK!")
				msg = &tmp

			default:
				// todo make task from any personal message
				continue
			}

			// send msg at the end

			if msg != nil {
				if _, err = bot.Send(msg); err != nil {
					log.Printf("failed to send message: %s", err)
				}
			}
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)
	sig := <-sigs
	log.Printf("Signal comming %s", sig)

	if os.Getenv("UDOIT_WEBHOOK") != "" {
		if _, err = bot.RemoveWebhook(); err != nil {
			log.Printf("failed to remove webhook %s", err)
		}
	}
}
