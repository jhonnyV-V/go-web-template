package server

import (
	"net/http"

	"snippet-sharing/cmd/web"
	"snippet-sharing/internal/server/routes"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	fileServer := http.FileServer(http.FS(web.Files))
	e.GET("/assets/*", echo.WrapHandler(fileServer))

	e.GET(
		"/login",
		echo.WrapHandler(
			templ.Handler(
				web.LoginForm("/auth/github"),
			),
		),
	)

	e.GET("/",
		echo.WrapHandler(
			templ.Handler(
				web.Home(),
			),
		),
	)

	e.GET("/hello", routes.HelloWebHandler, routes.ProtectedRoutesMiddlewarefunc)

	e.GET("/health", s.healthHandler)

	authGroup := e.Group("/auth")
	authGroup.GET("/:provider", routes.AuthBegin)
	authGroup.GET("/logout/:provider", routes.AuthLogout)
	authGroup.GET("/:provider/callback", routes.AuthCallback)

	return e
}

func (s *Server) HelloWorldHandler(c echo.Context) error {
	return c.String(http.StatusOK, "api is up and running")
}

func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}
