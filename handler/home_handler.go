package handler

import (
	"github.com/HackSquadDev/contributor-of-the-xxx-golang/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HomeHandler(c echo.Context, person types.RepositoryResponse,prs int) error {
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
	})
}
