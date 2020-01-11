package controllers

import (
	"github.com/labstack/echo"
	"log"
	"net/http"
	"remote-code-execution/services"
)

func RCEController(c echo.Context) error {
	language := c.Param("language")
	fileName := c.Get("fileName").(string)
	log.Println(fileName)
	services.CompileCode(language, fileName)
	return c.String(http.StatusOK, language)
}
