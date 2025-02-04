package routes

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)


var CompleteUserAuth = func(res http.ResponseWriter, req *http.Request) (goth.User, error) {
	providerName, err := gothic.GetProviderName(req)
	if err != nil {
		return goth.User{}, err
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return goth.User{}, err
	}

	value, err := gothic.GetFromSession(providerName, req)
	if err != nil {
		return goth.User{}, err
	}
	// defer Logout(res, req)
	sess, err := provider.UnmarshalSession(value)
	if err != nil {
		return goth.User{}, err
	}

	// err = gothic.validateState(req, sess)
	// if err != nil {
	// 	return goth.User{}, err
	// }

	user, err := provider.FetchUser(sess)
	if err == nil {
		// user can be found with existing session data
		return user, err
	}

	params := req.URL.Query()
	if params.Encode() == "" && req.Method == "POST" {
		req.ParseForm()
		params = req.Form
	}

	// get new token and retry fetch
	_, err = sess.Authorize(provider, params)
	if err != nil {
		return goth.User{}, err
	}

	err = gothic.StoreInSession(providerName, sess.Marshal(), req, res)

	if err != nil {
		return goth.User{}, err
	}

	gu, err := provider.FetchUser(sess)
	return gu, err
}

func AuthBegin(c echo.Context) error {
	provider := c.Param("provider")
	if provider == "" {
		return c.String(http.StatusBadRequest, "Provider not specified")
	}

	q := c.Request().URL.Query()
	q.Add("provider", c.Param("provider"))
	c.Request().URL.RawQuery = q.Encode()

	req := c.Request()
	res := c.Response().Writer
	//TODO: do something with the user
	_, err := CompleteUserAuth(res, req)
	if err == nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}
	gothic.BeginAuthHandler(res, req)
	return nil
}

func AuthCallback(c echo.Context) error {
	req := c.Request()
	res := c.Response().Writer

	//TODO: do something with the user
	_, err := CompleteUserAuth(res, req)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

func AuthLogout(c echo.Context) error {
	req := c.Request()
	res := c.Response().Writer
	err := gothic.Logout(res, req)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

type ContextWithUser struct {
	echo.Context
	User    goth.User
}

func (c ContextWithUser) GetUser() goth.User {
	return c.User
}

func ProtectedRoutesMiddlewarefunc(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		q := c.Request().URL.Query()
		// TODO: dynamically select the provider
		q.Add("provider", "github")
		c.Request().URL.RawQuery = q.Encode()

		req := c.Request()
		res := c.Response().Writer
		user, err := CompleteUserAuth(res, req)
		if err != nil {
			fmt.Println("ERROR")
			fmt.Println(err.Error())
			return c.Redirect(http.StatusTemporaryRedirect, "/login")
		}
		cc := &ContextWithUser{
			c,
			user,
		}
		return next(cc)
	}
}
