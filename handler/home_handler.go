package handler

import (
	"net/http"

	"github.com/HackSquadDev/contributor-of-the-xxx-golang/types"
	"github.com/labstack/echo/v4"
)

func HomeHandler(c echo.Context, person types.RepositoryResponse, prs int, orgData types.OrganizationResponse) error {
	// Please note the the second parameter "home.html" is the template name and should
	// be equal to the value stated in the {{ define }} statement in "view/home.html"
	return c.Render(http.StatusOK, "home.html", map[string]interface{}{
		"title":       "COTW",
		"msg":         "Hello, " + person.Author.Login,
		"name":        person.Author.Login,
		"profileLink": person.Author.Url,
		"avatarURL":   person.Author.AvatarURL,
		// number of prs made
		"prs": prs,
		// send the organization details
		"orgName":       orgData.Organization.Name,
		"orgAvatarURL":  orgData.Organization.AvatarUrl,
		"orgURL":        orgData.Organization.Url,
		"orgGithubName": orgData.Organization.Login,
		// TODO: send the dominant colors
		// "dominantColors": dominantColors,
	})
}
