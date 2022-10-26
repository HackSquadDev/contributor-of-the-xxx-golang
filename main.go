package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	ColorExtractor "github.com/HackSquadDev/contributor-of-the-xxx-golang/colorExtractor"

	"github.com/HackSquadDev/contributor-of-the-xxx-golang/handler"
	"github.com/HackSquadDev/contributor-of-the-xxx-golang/types"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/machinebox/graphql"
)

type TemplateRegistry struct {
	templates *template.Template
}

// Implement e.Renderer interface
func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
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

	// Initialize Echo
	e := echo.New()
	e.Renderer = &TemplateRegistry{
		templates: template.Must(template.ParseGlob("public/views/*.go.html")),
	}
	e.GET("/", home)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

func home(c echo.Context) error {
	client := graphql.NewClient("https://api.github.com/graphql")
	client2 := graphql.NewClient("https://api.github.com/graphql")
	orgData := getOrganizationDetails(client2)
	response := requestMergedPrs(client, "")
	hasNextPage := response.Search.PageInfo.HasNextPage
	// html := ""
	endCursor := ""
	PrCount := make(map[string]int)
	DataMap := make(map[string]types.RepositoryResponse)
	// bots have apps in their url.
	// Ex: https://github.com/apps/copybara-service
	r, err := regexp.Compile("^.*/apps/.*$")
	if err != nil {
		c.Logger().Panic("Unable to compile regexp: %v", err)
	}

	for hasNextPage {
		response = requestMergedPrs(client, endCursor)
		hasNextPage = response.Search.PageInfo.HasNextPage
		endCursor = response.Search.PageInfo.EndCursor
		for i := 0; i < len(response.Search.Nodes); i++ {
			// if it is a bot then don't count
			if r.MatchString(response.Search.Nodes[i].Author.Url) {
				continue
			}
			DataMap[response.Search.Nodes[i].Author.Login] = response.Search.Nodes[i]
			PrCount[response.Search.Nodes[i].Author.Login]++
		}
	}
	winnerName := ""
	highScore := 0

	for name, prs := range PrCount {
		/* html += name + ": " + fmt.Sprintf("%d", prs)
		html += "<br/>" */
		if prs > highScore {
			highScore = prs
			winnerName = name
		}
	}
	var winnerData types.RepositoryResponse
	c.Logger().Printf("Winner %s: PRs: %d", winnerName, highScore)
	for name, dude := range DataMap {
		if name == winnerName {
			winnerData = dude
			break
		}
	}
	// get organization details

	// get the organization avatar
	imageLink := orgData.Organization.AvatarUrl

	fmt.Printf("%v \n", imageLink)
	// get dominant colors from this image using colorExtractor module
	dominantColors, err := ColorExtractor.GetColors(imageLink)
	if err != nil {
		c.Logger().Panic("Unable to get colors: %v", err)
	}
	// log the dominant colors
	c.Logger().Printf("\nDominant colors: %v\n", dominantColors)

	// c.Logger().Printf("%v", winnerData)
	return handler.HomeHandler(c, winnerData, highScore, orgData, dominantColors)
	// return c.HTML(http.StatusOK, html)
}

// GraphQL Related Functions
func requestMergedPrs(client *graphql.Client, endCursor string) types.SearchResponse {
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
				author {
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
	}
	`, os.Getenv("GITHUB_ORG_NAME"), aWeekAgo, today, checkEndCursor(endCursor))
	// fmt.Printf(query)
	request := makeRequest(query)
	var resp types.SearchResponse
	err := client.Run(context.Background(), request, &resp)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

// graphql function to get organization details
func getOrganizationDetails(client *graphql.Client) types.OrganizationResponse {
	query := fmt.Sprintf(`
	{
		organization(login: "%s") {
			name
			url
			avatarUrl
			login
		}
	}
	`, os.Getenv("GITHUB_ORG_NAME"))
	request := makeRequest(query)
	var resp types.OrganizationResponse
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
