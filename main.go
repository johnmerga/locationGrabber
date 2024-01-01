package main

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	token := "tokens"
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		shareLocationBtn := []tgbotapi.KeyboardButton{
			tgbotapi.NewKeyboardButtonLocation("Share Location"),
		}
		keyboard := tgbotapi.NewReplyKeyboard(shareLocationBtn)
		keyboard.OneTimeKeyboard = true
		keyboard.ResizeKeyboard = true
		keyboard.Selective = true
		if update.Message != nil { // If we got a message
			if update.Message.IsCommand() && update.Message.Command() == "location" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please share your location")
				msg.ReplyToMessageID = update.Message.MessageID
				msg.ReplyMarkup = keyboard
				bot.Send(msg)
			}
			if update.Message.Chat.Type == "private" { // If it's not a group
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello, please Use me in a group chat")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			}
			// chech if the message is Replyied to a message
			if (update.Message.Chat.Type == "group" || update.Message.Chat.Type == "supergroup") && update.Message.ReplyToMessage != nil {
				location := update.Message.ReplyToMessage.Location
				isSameUser := update.Message.From.ID == update.Message.ReplyToMessage.From.ID
				isInEthiopia := isEthiopia(location)
				if location != nil {
					if isSameUser && isInEthiopia {
						userMessage := update.Message.Text
						botReply := fmt.Sprintf("Approved✅\n\nLatitude: %f\nLongitude: %f\nBranch: %s", location.Latitude, location.Longitude, userMessage)
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, botReply)
						msg.ReplyToMessageID = update.Message.MessageID
						bot.Send(msg)
					} else {
						botReply := "❌Rejected❌\n This location is not in Ethiopia. Please send a location in Ethiopia"
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, botReply)
						msg.ReplyToMessageID = update.Message.MessageID
						bot.Send(msg)
					}
				}
			}
		}
	}
}

// chech location range
func isEthiopia(location *tgbotapi.Location) bool {
	// Latitude: 3.4227 to 14.882
	// Longitude: 32.9986 to 47.9824
	lat := location.Latitude
	lon := location.Longitude
	if lat >= 3.4227 && lat <= 14.882 && lon >= 32.9986 && lon <= 47.9824 {
		return true
	}
	return false
}
