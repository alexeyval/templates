package main

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
	"reflect"
)

func main() {
	// подключаемся к боту с помощью токена
	bot, err := tgbotapi.NewBotAPI("5941964544:AAGonU-msKpkY9N-W4E0CblToGSA9bW8SWI")
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// инициализируем канал, куда будут прилетать обновления от API
	var ucfg = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	upd, _ := bot.GetUpdatesChan(ucfg)

	// читаем обновления из канала
	for {
		select {
		case update := <-upd:
			if update.Message == nil {
				continue
			}

			// Пользователь, который написал боту
			UserName := update.Message.From.UserName
			UserNameFirstName := update.Message.From.FirstName

			// ID чата/диалога.
			// Может быть идентификатором как чата с пользователем
			// (тогда он равен UserID) так и публичного чата/канала
			ChatID := update.Message.Chat.ID

			// Текст сообщения
			Text := update.Message.Text

			log.Printf("[%s, %s] %d %s", UserName, UserNameFirstName, ChatID, Text)

			if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
				switch {
				case Text == "/start":
					msgText := "Привет, "
					switch len(UserName) {
					case 0:
						msgText += UserNameFirstName
					default:
						msgText += UserName
					}

					msg := tgbotapi.NewMessage(
						ChatID,
						msgText)
					_, _ = bot.Send(msg)
				}
			}
		}
	}
}
