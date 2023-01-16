package main

import (
	"fmt"
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
			CallbackQueryHandler(update.CallbackQuery)

			log.Printf("[%s, %s] %d %s", UserName, UserNameFirstName, ChatID, Text)

			if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
				switch {
				case update.CallbackQuery != nil:
					CallbackQueryHandler(update.CallbackQuery)
					fmt.Println("Я тут")
					continue
				case Text == "/start":

				//bot.me(chat_id = message.chat.id, message_id = to_pin)
				default:
					fmt.Println(Text)
					msgText := "Поддержите меня 😎"
					//switch len(UserName) {
					//case 0:
					//	msgText += UserNameFirstName
					//default:
					//	msgText += UserName
					//}

					msg := tgbotapi.NewMessage(
						ChatID,
						msgText)

					v := "/поддержать"
					button := []tgbotapi.InlineKeyboardButton{{
						Text: "Поддержать 👍",
						URL:  &v,
					},
					}
					var buttons [][]tgbotapi.InlineKeyboardButton
					buttons = append(buttons, button)
					//markup := tgbotapi.InlineKeyboardMarkup{InlineKeyboard: buttons}
					//msg.ReplyMarkup = tgbotapi.NewInlineKeyboardButtonURL("Поддержать 👍", "/поддержать")
					//_ = []tgbotapi.InlineKeyboardButton{}
					//msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(button)

					var rows []tgbotapi.InlineKeyboardButton
					rows = append(rows, tgbotapi.NewInlineKeyboardButtonData("Next", "/поддержать"))
					//rows = append(rows, tgbotapi.NewInlineKeyboardButtonURL("Поддержать", "/поддержать"))
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

func CallbackQueryHandler(query *tgbotapi.CallbackQuery) {
	split := query.Data
	if split == "pager" {
		HandleNavigationCallbackQuery(query.Message.MessageID, split)
		return
	}
}

func HandleNavigationCallbackQuery(messageId int, data string) {
	pagerType := data
	_ = messageId

	fmt.Println(pagerType)
}
