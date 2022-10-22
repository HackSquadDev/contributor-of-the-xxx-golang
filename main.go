package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/machinebox/graphql"
)

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
	// Initialize Echo
	e := echo.New()

	e.GET("/", home)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

func home(c echo.Context) error {
	client := graphql.NewClient("https://api.github.com/graphql")
	response := requestData(client)
	return c.HTML(http.StatusOK, response.Licenses[0].Name)
}

// GraphQL Related Functions
func requestData(client *graphql.Client) GitHubResponse {
	query := `
	{
		licenses {
			name
			featured
		}
	}
	`
	request := graphql.NewRequest(query)
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		log.Fatal("GITHUB_TOKEN is empty.")
	}
	request.Header.Set("Authorization", fmt.Sprintf("bearer %s", githubToken))
	var resp GitHubResponse
	err := client.Run(context.Background(), request, &resp)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

type GitHubResponse struct {
	Licenses []struct {
		Name     string
		Featured bool
	}
}
