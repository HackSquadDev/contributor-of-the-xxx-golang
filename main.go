package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
)

type Organization struct {
	Repositories []Repository
	Logo         string
}
type Owner struct {
	Login string `json:"avatar_url"`
}
type Repository struct {
	Name      string
	Full_name string
}

func getPage(orgName string, pageNo int, ItemsPerPage int, GithubToken string) (Organization, error) {
	client := &http.Client{
		Timeout: time.Second * 20,
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/orgs/%s/repos?page=%d&per_page=%d", orgName, pageNo, ItemsPerPage), nil)
	if err != nil {
		return Organization{}, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", GithubToken))

	resp, err := client.Do(req)
	if err != nil {
		return Organization{}, err
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return Organization{}, err
	}
	var org Organization
	err = jsoniter.Unmarshal(bodyBytes, &org.Repositories)
	if err != nil {
		return Organization{}, err

	}
	org.Logo = jsoniter.Get(bodyBytes, 0, "owner", "avatar_url").ToString()
	return org, nil

}
func main() {
	// Load ENVs
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}
	e := echo.New()
	e.GET("/", home)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

func home(c echo.Context) error {
	orgName := os.Getenv("ORG")
	GithubToken := os.Getenv("GITHUB_TOKEN")
	var org Organization

	temp, err := getPage(orgName, 1, 100, GithubToken)
	for i := 2; len(temp.Repositories) != 0; i++ {
		c.Logger().Printf("Page %d", i)
		temp, err = getPage(orgName, i, 100, GithubToken)
		if err != nil {
			c.Logger().Fatal(err)
		}
		org.Repositories = append(org.Repositories, temp.Repositories...)
	}
	return c.HTML(http.StatusOK, fmt.Sprintf("%d\n", len(org.Repositories)))
}
