package main

import (
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
)

const pinMsg = "Hello"
const pinButton = "➡️"
const pinAnswer = "Thank you ❤️"

var sendAnswer []string

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(pinButton, pinAnswer),
	),
)

func main() {
	sendAnswer = append(sendAnswer, "Message 1")
	sendAnswer = append(sendAnswer, "Message 2")
	// подключаемся к боту с помощью токена
	bot, err := tgbotapi.NewBotAPI("TOKEN")
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
			case "/start":
				fallthrough
			case "open":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, pinMsg)
				msg.ReplyMarkup = numericKeyboard

				// Send the message.
				message, err := bot.Send(msg)
				if err != nil {
					fmt.Println(err)
					return
				}

				c := tgbotapi.PinChatMessageConfig{
					ChatID:              update.Message.Chat.ID,
					MessageID:           message.MessageID,
					DisableNotification: false,
				}
				_, err = bot.PinChatMessage(c)
				if err != nil {
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
			return

			// And finally, send a message containing the data received.
			for _, answer := range sendAnswer {
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, answer)
				msg.ParseMode = tgbotapi.ModeMarkdown
				if _, err = bot.Send(msg); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}
