package utility

import (
	"crypto/rand"
	"math/big"

	"github.com/jinzhu/gorm"
	"gitlab.com/shitposting/shitposting-bot/database/entities"
)

func NewDiscordMeme(db *gorm.DB) string {
	var meme []entities.Post

	var count int
	db.Not("posted_at", "NULL").Find(&meme).Count(&count)
	max := int64(count)

	db.Where("id = ?", getRand(max)).First(&meme)

	if len(meme) != 0 {
		return meme[0].Media
	}

	return ""
}

// getRand generates a random number
func getRand(max int64) (random int) {
	nBig, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		return
	}
	random = int(nBig.Int64())

	return
}
