package main

import (
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
	"reflect"
)

func main() {
	// подключаемся к боту с помощью токена
	bot, err := tgbotapi.NewBotAPI("TOKEN")
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
			//CallbackQueryHandler(update.CallbackQuery)

			log.Printf("[%s, %s] %d %s", UserName, UserNameFirstName, ChatID, Text)

			if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
				switch {
				case Text == "/start":

				default:
					msgText := "Поддержите меня 😎"
					msg := tgbotapi.NewMessage(
						ChatID,
						msgText)

					var rows []tgbotapi.InlineKeyboardButton
					rows = append(rows, tgbotapi.NewInlineKeyboardButtonData("text", "data"))
					msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows)

					message, err := bot.Send(msg)
					fmt.Println(err)

					c := tgbotapi.PinChatMessageConfig{
						ChatID:              ChatID,
						MessageID:           message.MessageID,
						DisableNotification: false,
					}
					_, _ = bot.PinChatMessage(c)
				}
			}
		}
	}
}
