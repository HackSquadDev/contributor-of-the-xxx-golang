package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/playwright-community/playwright-go"
)

func Setup() string {
	// can remove this line if in the package the load has been called already
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Setup playwright
	err = playwright.Install()
	if err != nil {
		panic(err)
	}
	// get the url from .env
	url := os.Getenv("IMAGE_URL")
	if url == "" {
		panic("No IMAGE_URL in .env")
	}
	return url
}
func GetImage() {
	// Setup
	url := Setup()
	// Launch playwright
	pw, err := playwright.Run()
	if err != nil {
		panic(err)
	}
	// Launch browser
	browser, err := pw.Chromium.Launch()
	if err != nil {
		panic(err)
	}
	// Launch context
	context, err := browser.NewContext()
	if err != nil {
		panic(err)
	}
	// Launch page
	page, err := context.NewPage()
	if err != nil {
		panic(err)
	}
	// Navigate to url
	_, err = page.Goto(url)
	if err != nil {
		panic(err)
	}
	// wait max 200 sec for div with class card
	_, err = page.WaitForSelector(".card", playwright.PageWaitForSelectorOptions{
		Timeout: playwright.Float(200000),
	})
	if err != nil {
		panic(err)
	}
	// wait till div with class card is present, timeout 200 sec
	_, err = page.WaitForSelector(".card", playwright.PageWaitForSelectorOptions{
		Timeout: playwright.Float(200000),
	})
	if err != nil {
		panic(err)
	}
	val, err := page.Locator(".card", playwright.PageLocatorOptions{})
	if err != nil {
		panic(err)
	}
	// omit background color in screenshot

	_, err = val.Screenshot(playwright.LocatorScreenshotOptions{
		// save to file with today's date
		Path:           playwright.String("images/" + time.Now().Format("2006-01-02") + ".png"),
		OmitBackground: playwright.Bool(true),
	})
	if err != nil {
		panic(err)
	}
	// close the browser
	err = browser.Close()
	if err != nil {
		panic(err)
	}

}
func main() {
	GetImage()
}
