package main

import (
	"bot/number"
	"log"
	"net/http"
	"strconv"
	"strings"

	//"github.com/bobrovka/templatebot/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var db map[int][]string = make(map[int][]string)

func main() {
	port := "8090"

	cfg := Config{
		TgToken: "1411904069:AAHr1KVp2ECCl4nolgs6qvj-cqBLDhn_mSo",
		Webhook: "https://5cc8bc7fa09f.ngrok.io",
	}

	go func() {
		_ = http.ListenAndServe(":"+port, nil)
	}()

	bot, err := tgbotapi.NewBotAPI(cfg.TgToken)
	if err != nil {
		log.Fatal("cannot new bot ", err)
	}
	log.Println("Bot created")

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(cfg.Webhook))
	if err != nil {
		log.Fatal("cannot new bot ", err)
	}
	log.Println("Webhook created")

	numClient := number.New("http://numbersapi.com")

	for update := range bot.ListenForWebhook("/") {
		msg := update.Message.Text

		userID := update.Message.From.ID
		if _, ok := db[userID]; !ok {
			db[userID] = []string{msg}
		} else {
			db[userID] = append(db[userID], msg)
		}

		if isGreeting(msg) {
			if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Hi, "+update.Message.From.FirstName)); err != nil {
				log.Println(err)
			}
		}

		if msg == "/history" {
			answer := strings.Join(db[userID], "\n")
			if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, answer)); err != nil {
				log.Println(err)
			}
		}

		num, err := strconv.Atoi(msg)
		if err != nil {
			if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Please, enter the number")); err != nil {
				log.Println(err)
			}
		} else {
			text, err := numClient.Info(num)
			if err != nil {
				log.Println(err)
			}

			if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, text)); err != nil {
				log.Println(err)
			}
		}
	}
}

var greetings []string = []string{"hi", "hello", "good morning", "good evening"}

func isGreeting(msg string) bool {
	m := strings.ToLower(msg)
	for _, greeting := range greetings {
		if m == greeting {
			return true
		}
	}

	return false
}
