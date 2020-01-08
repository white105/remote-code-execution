package middlewares

import (
	"fmt"
	"github.com/labstack/echo"
)

//Middleware create file & write code to file
func CreateSourceFile(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("Hi", c.Param("language"))
		return next(c)
	}
}
