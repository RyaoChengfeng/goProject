package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"time"
)

type KV struct {
	key    string
	values Values
}

type Values struct {
	value string
	time  string
}

type httpError struct {
	Message string
}

var kvMap = make(map[string]Values)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("api/key/:key", createJSON)
	e.GET("api/key/:key", sendJSON)
	e.PUT("api/key/:key", addKey)
	e.DELETE("api/key/:key", deleteKey)

	e.Logger.Fatal(e.Start(":1323"))
}

//  POST api/key/:key
func createJSON(c echo.Context) (err error) {
	key := c.Param(":key")
	_, ok := kvMap[key]
	if ok {
		return ErrorHandler(c, http.StatusBadRequest, "Key has already existed")
	}

	value := c.FormValue("value")
	if value == "" {
		return ErrorHandler(c, http.StatusBadRequest, "A valid value is needed")
	}

	tim := time.Now().Format("2006-01-02 15:04:05")
	values := Values{value, tim}

	kvMap[key] = values
	return c.JSONPretty(http.StatusOK, KV{key, values}, "	")
}

// GET api/key/:key
func sendJSON(c echo.Context) (err error) {

	key := c.Param("key")
	values, ok := kvMap[key]
	if !ok {
		return ErrorHandler(c, http.StatusBadRequest, "Key requested is not exist")
	}

	return c.JSONPretty(http.StatusOK, KV{key, values}, "	")
}

// PUT api/key/:key
func addKey(c echo.Context) (err error) {
	key := c.Param("key")
	_, ok := kvMap[key]
	if !ok {
		return ErrorHandler(c, http.StatusBadRequest, "Key requested is not exist")
	}
	value := c.FormValue("value")
	tim := time.Now().Format("2006-01-02 15:04:05")
	values := Values{value, tim}

	kvMap[key] = values

	return c.JSONPretty(http.StatusOK, KV{key, values}, "	")
}

// DELETE api/key/:key
func deleteKey(c echo.Context) (err error) {
	key := c.Param("key")
	_, ok := kvMap[key]
	if !ok {
		return ErrorHandler(c, http.StatusBadRequest, "Key requested is not exist")
	}

	delete(kvMap, key)

	return c.JSON(http.StatusOK, nil)
}

func ErrorHandler(c echo.Context, code int, msg string) error {
	return c.JSON(code, httpError{msg})
}
