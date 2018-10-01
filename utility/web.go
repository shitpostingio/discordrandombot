package utility

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"gitlab.com/shitposting/telegram-bot-api"
)

func GetFile(bot *tgbotapi.BotAPI, fileID string) (filePath string, err error) {

	imageDownloadURL, err := bot.GetFileDirectURL(fileID)

	if err != nil {
		return "", err
	}

	filePath = buildPath(fileID)
	err = downloadFile(filePath, imageDownloadURL)

	return
}

func buildPath(fileid string) string {

	fmt.Println(fileid)

	switch {
	case strings.HasPrefix(fileid, "BAAD"):
		return fileid + ".mp4"
	default:
		return fileid + ".jpg"
	}
}

//downloadFile downloads a file using a GET http request
func downloadFile(filepath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
