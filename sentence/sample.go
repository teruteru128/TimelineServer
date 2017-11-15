package sentence

import (
	"unicode/utf8"

	randomdata "github.com/Pallinder/go-randomdata"
	"github.com/TinyKitten/TimelineServer/models"
)

func getSampleText() string {
	text := randomdata.Paragraph()
	if utf8.RuneCountInString(text) > 140 {
		return text[:140]
	}
	return text
}

func GetRandomPost() (*models.User, *models.Post) {
	profile := randomdata.GenerateProfile(randomdata.RandomGender)
	sender := models.NewUser(
		profile.Login.Username,
		profile.Login.Md5,
		profile.Email,
	)
	sender.DisplayName = randomdata.Title(randomdata.RandomGender)

	post := models.NewPost(
		profile.Login.Username,
		getSampleText(),
	)

	return sender, post
}
