package main

import (
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/api/sheets/v4"
)

type PreviousMessage struct {
	chatID     int64
	messageID  int
	UserID     int64
	Message    *tgbotapi.Message
	IsLocation bool
}

type languages struct {
	eng string
	orm string
	amh string
}
type LanguageOptions struct {
	wrongLocation languages
	wrongTime     languages
	holyday       languages
	alreadyExist  languages
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
	// Get EAT location
	loc, err := time.LoadLocation("Africa/Nairobi")
	if err != nil {
		return time.Time{}, err
	}
	eat := time.Now().In(loc)
	return eat, nil
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	format := "Mon 02 Jan 2006 - 15:04:05"
	return t.Format(format)
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
			eng: "Sorry, you can only send Location only by being present at your branch from 8:00 AM to 5:00 PM.",
			orm: "Dhiifama, Location erguu kan dandeessan ganama sa'aatii 2:00 hanga galgala sa'aatii 11:00tti damee keessan irratti argamuun qofa.",
			amh: "ይቅርታ፣ Location መላክ የሚችሉት ከጠዋቱ 2፡00 እስከ ምሽቱ 11፡00 በቅርንጫፍዎ በመገኘት ብቻ ነው።",
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
		alreadyExist: struct {
			eng string
			orm string
			amh string
		}{
			eng: "Sorry, this location already exists in the database. you don't need to register it again.",
			orm: "Dhiifama, bakki kun duraan kuusdeetaa keessa jira. irra deebitee galmeessuun si hin barbaachisu.",
			amh: "ይቅርታ፣ ይህ አካባቢ አስቀድሞ በመረጃ ቋቱ ውስጥ አለ። እንደገና መመዝገብ አያስፈልግዎትም።",
		},
	}
	return &chooseLng
}

// is group admin
func isGroupAdmin(bot *tgbotapi.BotAPI, chatID int64, userID int64) bool {
	isAdmin := false
	admins, err := bot.GetChatAdministrators(tgbotapi.ChatAdministratorsConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: chatID,
		},
	})
	if err != nil {
		return false
	}
	for _, admin := range admins {
		if admin.User.ID == userID {
			isAdmin = true
		}
	}
	return isAdmin
}

// set chat previous log
func setPreviousMessage(p *PreviousMessage, chatID int64, messageID int, userID int64, isLocation bool, tg *tgbotapi.Message) {
	p.chatID = chatID
	p.messageID = messageID
	p.UserID = userID
	p.IsLocation = isLocation
	p.Message = tg
}

// is message location
func isMessageLocation(update *tgbotapi.Update) bool {
	if update.Message.Location != nil {
		return true
	}
	return false
}

// check if the previoius message and the current message is written by the same user
func isSameUsPrevMsg(prevMsg *PreviousMessage, update *tgbotapi.Update) bool {
	isMsgText := update.Message.Text != ""
	if prevMsg.UserID == update.Message.From.ID && isMsgText {
		return true
	}
	return false
}

// checks if the coordinates already exist in the database
func isLocationAlreadyExist(location *tgbotapi.Location, spreadsheetId string, srv *sheets.Service) bool {
	// rng
	rng := "Sheet1!B:B"
	coordinates := getCoordinateValues(spreadsheetId, rng, srv)
	for _, coordinate := range coordinates {
		if coordinate == fmt.Sprintf("%f,%f", location.Latitude, location.Longitude) {
			return true
		}
	}
	return false

}

func getCoordinateValues(spreadsheetId, rng string, srv *sheets.Service) []string {
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, rng).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
		return nil
	}
	// change to a list of strings
	var values []string
	for _, row := range resp.Values {
		for _, col := range row {
			values = append(values, fmt.Sprintf("%v", col))
		}
	}
	return values
}
