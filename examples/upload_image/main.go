package main

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"image/jpeg"

	"github.com/densestvoid/groupme"
)

// This is not a real token. Please find yours by logging
// into the GroupMe development website: https://dev.groupme.com/
const authorizationToken = "0123456789ABCDEF"

//go:embed groupme.jpeg
var imgBytes []byte

// A short program that gets the gets the first 5 groups
// the user is part of, and then the first 10 messages of
// the first group in that list
func main() {
	// Create a new client with your auth token
	client := groupme.NewClient(authorizationToken)

	img, err := jpeg.Decode(bytes.NewBuffer(imgBytes))
	if err != nil {
		fmt.Println(err)
		return
	}

	picture, err := client.UploadPicture(context.Background(), img, groupme.PictureEncodingJPEG)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(picture)
}
