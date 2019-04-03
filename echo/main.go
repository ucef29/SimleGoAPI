// using the echo framework

package main

import (
	"SimpleGoAPI/cache"
	"SimpleGoAPI/user"
	"net/http"
	"strings"

	"github.com/asdine/storm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/mgo.v2/bson"
)

type jsonResponse map[string]interface{}

func serveCache(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if cache.Serve(c.Response(), c.Request()) {
			return nil
		}
		return next(c)
	}
}

func cacheResponse(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Writer = cache.NewWriter(c.Response().Writer, c.Request())
		return next(c)
	}
}

func usersOptions(c echo.Context) error {
	methods := []string{http.MethodGet, http.MethodPost, http.MethodOptions}
	c.Response().Header().Set("Allow", strings.Join(methods, ","))
	return c.NoContent(http.StatusOK)
}

func userOptions(c echo.Context) error {
	methods := []string{http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions}
	c.Response().Header().Set("Allow", strings.Join(methods, ","))
	return c.NoContent(http.StatusOK)
}

func usersPostOne(c echo.Context) error {
	u := new(user.User)
	err := c.Bind(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	u.ID = bson.NewObjectId()
	err = u.Save()
	if err != nil {
		if err == user.ErrRecordInvalid {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	cache.Drop("/users")
	c.Response().Header().Set("Location", "/users/"+u.ID.Hex())
	return c.NoContent(http.StatusCreated)

}

func userGetOne(c echo.Context) error {
	if cache.Serve(c.Response(), c.Request()) {
		return nil
	}

	if !bson.IsObjectIdHex(c.Param("id")) {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	id := bson.ObjectIdHex(c.Param("id"))
	u, err := user.One(id)
	if err != nil {
		if err == storm.ErrNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, jsonResponse{"user": u})
}

func usersPutOne(c echo.Context) error {
	u := new(user.User)
	err := c.Bind(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	if !bson.IsObjectIdHex(c.Param("id")) {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	id := bson.ObjectIdHex(c.Param("id"))
	u.ID = id
	err = u.Save()
	if err != nil {
		if err == user.ErrRecordInvalid {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, jsonResponse{"user": u})
}

func usersPatchOne(c echo.Context) error {
	u := new(user.User)
	err := c.Bind(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	if !bson.IsObjectIdHex(c.Param("id")) {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	id := bson.ObjectIdHex(c.Param("id"))
	u.ID = id
	err = u.Save()
	if err != nil {
		if err == user.ErrRecordInvalid {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		return echo.NewHTTPError(http.StatusInternalServerError)

	}
	return c.JSON(http.StatusOK, jsonResponse{"user": u})
}

func usersDeleteOne(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	id := bson.ObjectIdHex(c.Param("id"))
	err := user.Delete(id)
	if err != nil {
		if err == storm.ErrNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)

}

func root(c echo.Context) error {
	return c.String(http.StatusOK, "Running API v1")
}

func auth(username, password string, c echo.Context) (bool, error) {
	if username == "joe" && password == "secret" {
		return true, nil
	}
	return false, nil
}

func main() {
	e := echo.New()
	e.GET("/", root)

	u := e.Group("/users")
	u.OPTIONS("", usersOptions)
	u.POST("", usersPostOne, middleware.BasicAuth(auth)) // implementing a basic auth middleWare

	uid := u.Group("/:id")
	uid.OPTIONS("", usersOptions)
	uid.GET("", userGetOne, serveCache, cacheResponse)                  // calling the serveCache and cacheResponse middleWares
	uid.PUT("", usersPutOne, middleware.BasicAuth(auth), cacheResponse) // could add multiple middleWares exec left -> rigjt
	uid.PATCH("", usersPatchOne, middleware.BasicAuth(auth), cacheResponse)
	uid.DELETE("", usersDeleteOne, middleware.BasicAuth(auth))

	e.Start(":12345")
}
