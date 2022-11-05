package twitter

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

func PostToTwitter(msg string, imgPath string) {

	// get TWITTER_API_URL from .env
	twitterApiUrl := os.Getenv("TWITTER_API_URL")
	var b bytes.Buffer

	w := multipart.NewWriter(&b)
	// use multipart writer to read image from path

	img, err := os.Open(imgPath)
	if err != nil {
		panic(err)
	}
	fw, err := w.CreateFormFile("upload", imgPath)
	if err != nil {
		panic(err)
	}
	if fw, err = w.CreateFormField("text"); err != nil {
		panic(err)
	}
	if _, err = io.Copy(fw, strings.NewReader(msg)); err != nil {
		panic(err)
	}
	w.Close()
	log.Default().Println("Image opened successfully")
	defer img.Close()

	if err != nil {
		panic(err)
	}
	log.Default().Println("Post body created successfully")
	// make POST request to the url, with the body
	// and the auth header
	req, err := http.NewRequest(http.MethodPost, twitterApiUrl, &b)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	req.Header.Set("Content-Type", w.FormDataContentType())
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
