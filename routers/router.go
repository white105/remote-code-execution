package routers

import (
	"github.com/labstack/echo"
	"remote-code-execution/controllers"
	"remote-code-execution/middlewares"
)

//Declare all the routers here
func InitRouters(e *echo.Echo) {
	api := e.Group("/api")

	rce := api.Group("/rce")

	rce.POST("/:language", controllers.RCEController, middlewares.CreateSourceFile)
}
