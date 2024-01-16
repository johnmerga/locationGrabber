package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func main() {
	debug := flag.Bool("debug", false, "debug")
	flag.Parse()
	fmt.Printf("\n\n\nDocker time: %s\n\n\n", humanDate(time.Now().UTC()))
	token := os.Getenv("TELEGRAM_API_KEY")
	gkey := os.Getenv("GOOGLE_API_JSON")
	spreadsheetId := os.Getenv("SPREADSHEET_ID")
	if token == "" || gkey == "" || spreadsheetId == "" {
		if token == "" {
			fmt.Println("TELEGRAM_API_KEY is missing")
		}
		if gkey == "" {
			fmt.Println("GOOGLE_API_KEY is missing")
		}
		if spreadsheetId == "" {
			fmt.Println("SPREADSHEET_ID is missing")
		}
		os.Exit(1)
	}
	chooseLng := getLang()

	ctx := context.Background()
	data, err := os.ReadFile(gkey)
	if err != nil {
		log.Fatal(err)
	}
	conf, err := google.JWTConfigFromJSON(data, sheets.SpreadsheetsScope)
	if err != nil {
		log.Fatal(err)
	}
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(conf.Client(ctx)))
	if err != nil {
		log.Fatal(err)
	}

	rng := "Sheet1!A1:B2"

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, rng).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) > 0 {
		fmt.Println("Sheet data:")
		fmt.Println(resp.Values)
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = *debug

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
		// test
		if update.Message != nil { // If we got a message
			if update.Message.Chat.Type == "private" { // If it's not a group
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello, please Use me in a group chat")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
				continue
			}
			// chech if the message is Replyied to a message
			if (update.Message.Chat.Type == "group" || update.Message.Chat.Type == "supergroup") && update.Message.ReplyToMessage != nil {
				location := update.Message.ReplyToMessage.Location
				isSameUser := update.Message.From.ID == update.Message.ReplyToMessage.From.ID
				if location != nil {
					isInEthiopia := isEthiopia(location)
					isWorkingHrs, currentTime, err := isWorkingHours()
					if err != nil {
						log.Fatal(err)
						continue
					}
					humanTime := humanDate(currentTime)
					fmt.Printf("\n\n\nCurrent time: %s\n\n\n", humanTime)
					if strings.Contains(humanTime, "Sun") {
						rpl := fmt.Sprintf("%s\n\n%s\n\n%s", chooseLng.holyday.eng, chooseLng.holyday.orm, chooseLng.holyday.amh)
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, rpl)
						msg.ReplyToMessageID = update.Message.MessageID
						bot.Send(msg)
						continue
					}

					if !isWorkingHrs {
						rpl := fmt.Sprintf("%s\n\n%s\n\n%s", chooseLng.wrongTime.eng, chooseLng.wrongTime.orm, chooseLng.wrongTime.amh)
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, rpl)
						msg.ReplyToMessageID = update.Message.MessageID
						bot.Send(msg)
						continue
					}
					isGroupAdmin := isGroupAdmin(bot, update.Message.Chat.ID, update.Message.From.ID)

					if isSameUser && isInEthiopia || (isGroupAdmin && isInEthiopia) {
						userMessage := update.Message.Text
						username := update.Message.From.UserName
						userFullName := fmt.Sprintf("%s %s", update.Message.From.FirstName, update.Message.From.LastName)

						// add data to sheet
						var vr sheets.ValueRange
						mapCoordination := fmt.Sprintf("%f,%f", location.Latitude, location.Longitude)
						branch := fmt.Sprintf("Siinqee Bank %s", userMessage)
						vr.Values = append(vr.Values, []interface{}{branch, mapCoordination, userFullName, username})
						_, err := srv.Spreadsheets.Values.Append(spreadsheetId, rng, &vr).ValueInputOption("RAW").Do()
						if err != nil {
							log.Fatalf("Unable to retrieve data from sheet: %v", err)
						}
						botReply := fmt.Sprintf("Approved✅\n\nLatitude: %f\nLongitude: %f\nBranch: %s\n\n\nThank You %s.", location.Latitude, location.Longitude, userMessage, userFullName)
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, botReply)
						msg.ReplyToMessageID = update.Message.MessageID
						bot.Send(msg)
					} else if (isSameUser && !isInEthiopia) || (isGroupAdmin && !isInEthiopia) {
						botReply := fmt.Sprintf("%s\n\n%s\n\n%s", chooseLng.wrongLocation.eng, chooseLng.wrongLocation.orm, chooseLng.wrongLocation.amh)
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, botReply)
						msg.ReplyToMessageID = update.Message.MessageID
						bot.Send(msg)
					}
				}
			}
		}
	}
}
