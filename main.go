package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

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
	hasNextPage := response.Search.PageInfo.HasNextPage
	html := ""
	endCursor := ""
	PrCount := make(map[string]int)
	// bots have apps in their url.
	// Ex: https://github.com/apps/copybara-service
	r, err := regexp.Compile("^.*/apps/.*$")

	for hasNextPage {
		response = requestOrganization(client, endCursor)
		hasNextPage = response.Search.PageInfo.HasNextPage
		endCursor = response.Search.PageInfo.EndCursor
		// html += response.Search.Nodes[0].Author.Login + "<br />"
		for i := 0; i < len(response.Search.Nodes); i++ {
			// if it is a bot
			if r.MatchString(response.Search.Nodes[i].Author.Url) {
				continue
			}

			PrCount[response.Search.Nodes[i].Author.Login]++
		}
	}
	winnerName := ""
	mx := 0
	if err != nil {
		c.Logger().Panic("Unable to compile regexp: %v", err)
	}

	for name, prs := range PrCount {
		html += name + ": " + fmt.Sprintf("%d", prs)
		html += "<br/>"
		if prs > mx {
			mx = prs
			winnerName = name
		}
	}
	var winnerDude RepositoryResponse
	c.Logger().Printf("Winner %s: PRs: %d", winnerName, mx)
	for _, dude := range response.Search.Nodes {
		if dude.Author.Login == winnerName {
			winnerDude = dude
			break
		}

	}
	c.Logger().Printf("%v", winnerDude)
	// TODO: Store the data in some Data structure
	return c.HTML(http.StatusOK, html)
}

// GraphQL Related Functions
func requestOrganization(client *graphql.Client, endCursor string) SearchResponse {
	// NOTE: DON'T EVEN THINK ABOUT TOUCHING THESE LINES.
	// refer https://stackoverflow.com/questions/33119748/convert-time-time-to-string#comment70144458_33119937
	// YYYY-MM-DD

	aWeekAgo := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	today := time.Now().Format("2006-01-02")
	query := fmt.Sprintf(`
{
  search(
    first: 100
    type: ISSUE,
    query: "org:%s is:pr is:merged merged:%s..%s -author:robot"
    %s
  ) {
    nodes {
      ... on PullRequest {
        title
        url
        author{
          avatarUrl
          url
          login
        }
      }
    }
    pageInfo{
      hasNextPage
      endCursor
    }
  }
}	`, os.Getenv("GITHUB_ORG_NAME"), aWeekAgo, today, checkEndCursor(endCursor))
	// fmt.Printf(query)
	request := makeRequest(query)
	var resp SearchResponse
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

type SearchResponse struct {
	Search struct {
		Nodes    []RepositoryResponse
		PageInfo struct {
			HasNextPage bool
			EndCursor   string
		}
	}
}
type RepositoryResponse struct {
	Title  string
	Url    string
	Author struct {
		AvatarURL string
		Login     string
		Url       string
	}
}
