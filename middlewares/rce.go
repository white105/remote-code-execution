package middlewares

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"log"
	"os"
)

type SourceCode struct {
	Code string `json:"code"`
}

//Middleware create file & write code to file
func CreateSourceFile(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		language := c.Param("language")
		var content SourceCode
		if err := c.Bind(&content); err != nil {
			log.Fatal(err)
		}
		fileName := fmt.Sprintf("%s.%s", uuid.New().String(), language)
		file, err := os.Create(fmt.Sprintf("./project/%s", fileName))
		defer file.Close()
		bytesread, err := file.WriteString(content.Code)
		if err != nil {
			log.Fatal("Error create file: ", err)
		}
		log.Printf("Wrote %d byte", bytesread)
		file.Sync()
		c.Set("fileName", fileName)
		return next(c)
	}
}
