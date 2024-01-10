package main

import (
	"context"
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

type languages struct {
	eng string
	orm string
	amh string
}
type LanguageOptions struct {
	wrongLocation languages
	wrongTime     languages
	holyday       languages
}

// main
func main() {

	fmt.Printf("\n\n\nDocker time: %s\n\n\n", time.Now().UTC())
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

					if isSameUser && isInEthiopia {
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
						botReply := fmt.Sprintf("Approved✅\n\nLatitude: %f\nLongitude: %f\nBranch: %s", location.Latitude, location.Longitude, userMessage)
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, botReply)
						msg.ReplyToMessageID = update.Message.MessageID
						bot.Send(msg)
					} else if isSameUser && !isInEthiopia {
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

// chech working hours
func isWorkingHours() (bool, time.Time, error) {
	currentTime, err := convertFranceToEastAfricaTime()
	if err != nil {
		return false, time.Time{}, err
	}
	if currentTime.Hour() >= 8 && currentTime.Hour() <= 17 {
		return true, currentTime, nil
	} else {
		return false, currentTime, nil
	}
}

// convert france time to east africa time
func convertFranceToEastAfricaTime() (time.Time, error) {
	// Example UTC date time
	now := time.Now()

	// Get EAT location
	loc, err := time.LoadLocation("Africa/Nairobi")
	if err != nil {
		return time.Time{}, err
	}

	// Convert UTC to EAT
	eat := now.In(loc)
	return eat, nil

}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	// Convert the time to UTC before formatting it
	return t.Format("Tue 02 Jan 2006 - 15:04")
}

func getLang() *LanguageOptions {
	chooseLng := LanguageOptions{
		wrongLocation: struct {
			eng string
			orm string
			amh string
		}{
			eng: "❌Rejected!❌\n This location is not in Ethiopia. only locations in Ethiopia are accepted. Please send again",
			orm: "❌Hin fudhatamne!❌\nBakki kun Itoophiyaa keessa hin jiru. bakkeewwan Itoophiyaa keessa jiran qofatu fudhatama qaba. Irra deebi’uun nuuf ergaa",
			amh: "❌ውድቅ ተደርጓል!❌\nይህ ቦታ ኢትዮጵያ ውስጥ አይደለም፣ ኢትዮጵያ ውስጥ ያሉ ቦታዎች ብቻ ተቀባይነት አላቸው። እባኮትን በድጋሚ ላኩ",
		},
		wrongTime: struct {
			eng string
			orm string
			amh string
		}{
			eng: "Sorry, You can only send your location between 8:00 AM and 5:00 PM",
			orm: "Dhiifama, ganama sa'aatii 2:00 hanga galgala sa'aatii 11:00 gidduutti qofa Locationi erguu dandeessu",
			amh: "ይቅርታ፣ Location መላክ የሚችሉት ከጠዋቱ 2፡00 እስከ ምሽቱ 11፡00 ሰዓት ብቻ ነው።",
		},
		holyday: struct {
			eng string
			orm string
			amh string
		}{
			eng: "Sorry, you can't send your location today. Please share your location from Monday to Saturday between 8:00 AM and 5:00 PM by being at your branch",
			orm: "Dhiifama, har'a Locationi erguu hin dandeessan. Maaloo bakka jirtan Isniina hanga Dilbataatti sa'aatii 2:00 AM hanga 1:00 PM gidduutti damee keessan irratti argamuun nuuf qoodaa",
			amh: "ይቅርታ፣ ዛሬ Location መላክ አትችልም። እባኮትን ከሰኞ እስከ ቅዳሜ ከጠዋቱ 2፡00 እስከ 11፡00 ሰአት ባለው ጊዜ ውስጥ በቅርንጫፍዎ በመገኘት ያካፍሉ።",
		},
	}
	return &chooseLng
}
