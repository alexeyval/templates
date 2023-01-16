package main

import (
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
	"strconv"
	"strings"
)

var count = 2
var maxPages = len(data) / count // = 5

var data = []string{"DummyData1", "DummyData2", "DummyData3", "DummyData4", "DummyData5", "DummyData6", "DummyData7", "DummyData8", "DummyData9", "DummyData10"}

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

	var chatId = int64(746200068) // <--- Place Chat Id Here

	SendDummyData(bot, chatId, 0, 2, nil) // Send initial data

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

			/*if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" */
			{
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
				case update.CallbackQuery != nil:
					CallbackQueryHandler(bot, update.CallbackQuery, ChatID)
					continue
				}
			}
		}
	}
}

func CallbackQueryHandler(bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, ChatID int64) {
	split := strings.Split(query.Data, ":")
	if split[0] == "pager" {
		HandleNavigationCallbackQuery(bot, query.Message.MessageID, ChatID, split[1:]...)
		return
	}
}

func HandleNavigationCallbackQuery(bot *tgbotapi.BotAPI, messageId int, ChatID int64, data ...string) {
	pagerType := data[0]
	currentPage, _ := strconv.Atoi(data[1])
	itemsPerPage, _ := strconv.Atoi(data[2])

	if pagerType == "next" {
		nextPage := currentPage + 1
		if nextPage < maxPages {
			SendDummyData(bot, ChatID, nextPage, itemsPerPage, &messageId)
		}
	}
	if pagerType == "prev" {
		previousPage := currentPage - 1
		if previousPage >= 0 {
			SendDummyData(bot, ChatID, previousPage, itemsPerPage, &messageId)
		}
	}
}

func SendDummyData(bot *tgbotapi.BotAPI, chatId int64, currentPage, count int, messageId *int) {
	text, keyboard := DummyDataTextMarkup(currentPage, count)

	var cfg tgbotapi.Chattable
	if messageId == nil {
		msg := tgbotapi.NewMessage(chatId, text)
		msg.ReplyMarkup = keyboard
		cfg = msg
	} else {
		msg := tgbotapi.NewEditMessageText(chatId, *messageId, text)
		msg.ReplyMarkup = &keyboard
		cfg = msg
	}

	bot.Send(cfg)
}

func DummyDataTextMarkup(currentPage, count int) (text string, markup tgbotapi.InlineKeyboardMarkup) {
	text = strings.Join(data[currentPage*count:currentPage*count+count], "\n")

	var rows []tgbotapi.InlineKeyboardButton
	if currentPage > 0 {
		rows = append(rows, tgbotapi.NewInlineKeyboardButtonData("Previous", fmt.Sprintf("pager:prev:%d:%d", currentPage, count)))
	}
	if currentPage < maxPages-1 {
		rows = append(rows, tgbotapi.NewInlineKeyboardButtonData("Next", fmt.Sprintf("pager:next:%d:%d", currentPage, count)))
	}

	markup = tgbotapi.NewInlineKeyboardMarkup(rows)
	return
}
