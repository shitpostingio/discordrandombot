package utility

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/shitposting/discord-random/entities"
)

//GetRandomFileID returns a random file_id from the database
func GetRandomMeme(db *gorm.DB) string {

	var meme entities.Post

	// SELECT message_id FROM `posts`  WHERE NOT (posted_at IS NULL OR message_id = 0) ORDER BY rand(),`posts`.`id` ASC LIMIT 1
	db.Select("file_id").Not("posted_at IS NULL OR message_id = 0").Order("rand()").First(&meme)

	return meme.FileID
}
