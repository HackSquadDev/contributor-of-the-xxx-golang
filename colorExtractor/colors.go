package ColorExtractor

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/color"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/HackSquadDev/contributor-of-the-xxx-golang/types"
	"github.com/go-playground/colors"
	"github.com/mccutchen/palettor"
	"github.com/nfnt/resize"
)

func imageDownloader(url string) (string, error) {
	fmt.Printf("Downloading image from %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	file, err := os.Create("image.png")
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}
	return "image.png", nil
}

func FindDomiantColor(fileInput string, numOfColors int) ([]string, error) {
	f, err := os.Open(fileInput)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if err != nil {
		fmt.Println("File not found:", fileInput)
		return nil, err
	}
	originalImg, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	img := resize.Thumbnail(800, 800, originalImg, resize.Lanczos3)
	temp, err := palettor.Extract(numOfColors, 20, img)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Dominant colors: %v", temp)
	// convert each entry in temp to hex and return the new array
	var hexColors []string
	for _, color := range temp.Colors() {
		// convert color to hex
		log.Printf("\n\ncolor: %v; weight: %v\n\n", color, temp.Weight(color))
		temp := colors.FromStdColor(color)
		hexColors = append(hexColors, temp.ToHEX().String())
		// colors = append(colors, color)
	}
	return hexColors, nil

}

func GetColorsv1(url string, numOfColors int) ([]string, error) {
	// download the image
	fileName, err := imageDownloader(url)
	if err != nil {
		fmt.Println("AAAA Unable to download image: ", err)
		return nil, err
	}
	fmt.Printf("Downloaded image to %s", fileName)
	// find the dominant color
	color, err := FindDomiantColor(fileName, numOfColors)
	if err != nil {
		return nil, err
	}
	fmt.Println("Dominant color: ", color)
	return color, nil
}
func GetColorsv2(url string) (types.ThemeColors, error) {
	client := &http.Client{}
	imagga_api_key := os.Getenv("IMAGGA_API_KEY")
	imagga_api_secret := os.Getenv("IMAGGA_API_SECRET")
	req, err := http.NewRequest("GET", "https://api.imagga.com/v2/colors?image_url="+url, nil)
	if err != nil {
		return types.ThemeColors{}, err
	}
	req.SetBasicAuth(imagga_api_key, imagga_api_secret)
	resp, err := client.Do(req)
	if err != nil {
		return types.ThemeColors{}, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return types.ThemeColors{}, err
	}
	// log.Println(string(respBody))
	var result types.ColorResponse
	json.Unmarshal([]byte(respBody), &result)
	// log result
	// log.Printf("result: %v", result)
	var colors types.ThemeColors
	for _, hex := range result.Result.Colors.BackgroundColors {
		colors.BackgroundColors = append(colors.BackgroundColors, hex.HTMLCode)
	}
	for _, hex := range result.Result.Colors.ForegroundColors {
		colors.ForegroundColors = append(colors.ForegroundColors, hex.HTMLCode)
	}
	for _, hex := range result.Result.Colors.ImageColors {
		colors.ImageColors = append(colors.ImageColors, hex.HTMLCode)
	}
	log.Printf("colors: %v", colors)
	return colors, nil

}
