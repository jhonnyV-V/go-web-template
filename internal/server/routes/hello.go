package routes

import (
	"log"
	"net/http"
	"snippet-sharing/cmd/web"

	"github.com/labstack/echo/v4"
)

func HelloWebHandler(c echo.Context) error {
	r := c.Request()
	w := c.Response()

	cc, ok := c.(*ContextWithUser)

	if !ok {
		c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	email := cc.User.Email
	component := web.HelloPost(email)
	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in HelloWebHandler: %e", err)
	}

	return nil
}
