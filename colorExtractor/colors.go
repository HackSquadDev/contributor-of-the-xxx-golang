package ColorExtractor

import (
	"encoding/json"
	_ "image/color"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/HackSquadDev/contributor-of-the-xxx-golang/types"
)

func GetColors(url string) (types.ThemeColors, error) {
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
	var result types.ColorResponse
	json.Unmarshal([]byte(respBody), &result)
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
