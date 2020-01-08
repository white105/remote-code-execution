package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"remote-code-execution/routers"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"timestamp":"${time_rfc3339}","request_id":"${id}","host":"${host}",` +
			`"method":"${method}","uri":"${uri}","status":${status},"error":"${error}","latency":${latency},` + "\n",
	}))

	routers.InitRouters(e)

	e.Logger.Fatal(e.Start(":1323"))
}
