package ColorExtractor

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"

	"github.com/cenkalti/dominantcolor"
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
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	temp := dominantcolor.Find(img)
	// make temp an array of 1 string
	colors := make([]string, 1)
	colors = append(colors, dominantcolor.Hex(temp))
	// convert each entry in temp to hex and return the new array
	// var colors []string
	// for _, color := range temp {
	// 	colors = append(colors, dominantcolor.Hex(color))
	// }
	return colors, nil

}

func GetColors(url string, numOfColors int) ([]string, error) {
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
