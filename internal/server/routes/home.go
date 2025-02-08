package routes

import (
	"log"
	"net/http"
	"snippet-sharing/cmd/web/pages"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
)

func HomeWebHandler(c echo.Context) error {
	r := c.Request()
	w := c.Response()

	cc, ok := c.(*ContextWithUser)

	var user goth.User

	if !ok {
		user = goth.User{}
	} else {
		user = cc.User
	}

	component := pages.Home(user)
	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in HomeWebHandler: %e", err)
	}

	return nil
}
