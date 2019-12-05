package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func downloadFile(filepath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer closeSafely(out)

	// Get the data
	resp, err := http.Get(url) // nolint: gosec
	if err != nil {
		return err
	}
	defer closeSafely(resp.Body)

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func closeSafely(toClose io.Closer) {
	err := toClose.Close()
	if err != nil {
		log.Println(err)
	}
}
