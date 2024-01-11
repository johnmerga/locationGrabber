package main

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
