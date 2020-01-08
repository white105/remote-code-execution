package controllers

import (
	"github.com/labstack/echo"
	"net/http"
)

func RCEController(c echo.Context) error {
	language := c.Param("language")
	return c.String(http.StatusOK, language)
}
