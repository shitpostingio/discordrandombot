package utility

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/shitposting/tg-random-bot/database/entities"
)

//GetRandomFileID returns a random file_id from the database
func GetRandomFileID(db *gorm.DB) string {

	var meme entities.Post

	// SELECT message_id FROM `posts`  WHERE NOT (posted_at IS NULL OR message_id = 0) ORDER BY rand(),`posts`.`id` ASC LIMIT 1
	db.Select("media").Not("posted_at IS NULL OR message_id = 0").Order("rand()").First(&meme)

	return meme.Media
}
