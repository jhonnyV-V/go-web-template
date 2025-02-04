package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
	_ "github.com/joho/godotenv/autoload"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"

	"snippet-sharing/internal/database"
)

type Server struct {
	port int

	db database.Service
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	githubKey := os.Getenv("GITHUB_KEY")
	githubSecret := os.Getenv("GITHUB_SECRET")
	githubCallbackUrl := os.Getenv("GITHUB_CALLBACK_URL")
	sessionKeyPair := os.Getenv("SESSIONS_SECRET")

	store := sessions.NewCookieStore([]byte(sessionKeyPair))
	store.Options.Path = "/"
	store.Options.Domain = ""
	store.Options.MaxAge = 86400 * 30
	// store.Options.HttpOnly = true
	store.Options.Secure = false
	gothic.Store = store

	goth.UseProviders(
		github.New(githubKey, githubSecret, githubCallbackUrl),
	)

	NewServer := &Server{
		port: port,

		db: database.New(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
