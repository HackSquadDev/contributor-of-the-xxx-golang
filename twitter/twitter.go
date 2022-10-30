package twitter

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

func PostToTwitter(msg string, imgPath string) {

	// get TWITTER_API_URL from .env
	twitterApiUrl := os.Getenv("TWITTER_API_URL")
	// open the imgPath file to send to twitter
	img, err := os.Open(imgPath)
	if err != nil {
		panic(err)
	}
	log.Default().Println("Image opened successfully")
	defer img.Close()

	postBody, err := json.Marshal(map[string]interface{}{
		"text":   msg,
		"upload": img,
	})
	if err != nil {
		panic(err)
	}
	log.Default().Println("Post body created successfully")
	// make POST request to the url, with the body
	// and the auth header
	req, err := http.NewRequest("POST", twitterApiUrl, bytes.NewBuffer(postBody))
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	req.Header.Set("Content-Type", "application/json")
	// send request
	client := &http.Client{}
	resp, err := client.Do(req)
	log.Default().Println("Request sent")
	if err != nil {
		panic(err)
	}
	log.Default().Println("Response received")
	defer resp.Body.Close()
	// read response
	body, _ := io.ReadAll(resp.Body)
	log.Default().Println("Response: ", string(body))

}
