package main

import (
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
)

const pinMsg = "Поблагодарить автора бота 😎"
const pinButton = "СДЕЛАТЬ ПО КАЙФУ"
const pinAnswer = "Поддержать проект: `4276 3802 1719 2553`"

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(pinButton, pinAnswer),
	),
)

func main() {
	// подключаемся к боту с помощью токена
	bot, err := tgbotapi.NewBotAPI("5941964544:AAGonU-msKpkY9N-W4E0CblToGSA9bW8SWI")
	if err != nil {
		log.Panic(err)
	}

	//bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	// Loop through each update.
	for update := range updates {
		// Check if we've gotten a message update.
		if update.Message != nil {
			// If the message was open, add a copy of our numeric keyboard.
			switch update.Message.Text {
			case "open":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, pinMsg)
				msg.ReplyMarkup = numericKeyboard

				// Send the message.
				if _, err = bot.Send(msg); err != nil {
					fmt.Println(err)
					return
				}
			}
		} else if update.CallbackQuery != nil {
			// Respond to the callback query, telling Telegram to show the user
			// a message with the data received.
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err = bot.AnswerCallbackQuery(callback); err != nil {
				fmt.Println(err)
				return
			}

			// And finally, send a message containing the data received.
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			msg.ParseMode = tgbotapi.ModeMarkdown
			if _, err = bot.Send(msg); err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
