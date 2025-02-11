package routes

import (
	"log"
	"net/http"
	"snippet-sharing/cmd/web/pages"

	"github.com/labstack/echo/v4"
)

func HelloWebHandler(c echo.Context) error {
	r := c.Request()
	w := c.Response()

	cc, ok := c.(*ContextWithUser)

	if !ok {
		c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	component := pages.HelloPost(cc.User)
	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in HelloWebHandler: %e", err)
	}

	return nil
}
