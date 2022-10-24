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
	response := requestOrganization(client, "")
	hasNextPage := response.Organization.Repositories.PageInfo.HasNextPage
	html := ""
	endCursor := ""

	for hasNextPage {
		response = requestOrganization(client, endCursor)
		hasNextPage = response.Organization.Repositories.PageInfo.HasNextPage
		endCursor = response.Organization.Repositories.PageInfo.EndCursor
		// sending data as HTML for now
		for i := 0; i < len(response.Organization.Repositories.Nodes); i++ {
			html += response.Organization.Repositories.Nodes[i].Name
			html += "<br/>"
		}
	}
	// TODO: Store the data in some Data structure
	return c.HTML(http.StatusOK, html)
}

// GraphQL Related Functions
func requestOrganization(client *graphql.Client, endCursor string) ResponseOrganization {
	query := fmt.Sprintf(`
	{
		organization(login:"%s") {
		  id
		  name
		  login
		  url
		  avatarUrl
		  repositories(orderBy:{field:PUSHED_AT, direction:DESC}, first:100, privacy:PUBLIC, isFork:false, %s) {
			totalCount
			pageInfo {
			  startCursor
			  endCursor
			  hasNextPage
			  hasPreviousPage
			}
			nodes {
			  id
			  name
			  description
			  url
			  stargazerCount
			}
		  }
		}
	  }
	`, os.Getenv("GITHUB_ORG_NAME"), checkEndCursor(endCursor))
	request := makeRequest(query)
	var resp ResponseOrganization
	err := client.Run(context.Background(), request, &resp)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func checkEndCursor(cursor string) string {
	if cursor != "" {
		return fmt.Sprintf("after:\"%s\"", cursor)
	}
	return ""
}

func makeRequest(query string) *graphql.Request {
	request := graphql.NewRequest(query)
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		log.Fatal("GITHUB_TOKEN is empty.")
	}
	request.Header.Set("Authorization", fmt.Sprintf("bearer %s", githubToken))
	return request
}

type ResponseOrganization struct {
	Organization struct {
		Id           string
		Name         string
		Login        string
		Url          string
		AvatarUrl    string
		Repositories struct {
			TotalCount int
			PageInfo   struct {
				StartCursor     string
				EndCursor       string
				HasNextPage     bool
				hasPreviousPage bool
			}
			Nodes []ResponseRepository
		}
	}
}

type ResponseRepository struct {
	Id             string
	Name           string
	Description    string
	Url            string
	StargazerCount int
}
